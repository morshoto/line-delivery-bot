package httpapi

import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "io"
    "net/http"
    "strings"

    "github.com/go-chi/chi/v5"

    "example.com/delivery-bot/server/internal/dedupe"
    "example.com/delivery-bot/server/internal/line"
    "example.com/delivery-bot/server/internal/logging"
    "example.com/delivery-bot/server/internal/qr"
    "example.com/delivery-bot/server/internal/security"
)

type ScanRequest struct {
    GroupID     string `json:"group_id"`
    QRText      string `json:"qr_text"`
    DisplayName string `json:"display_name"`
}

type ScanResponse struct {
    Status   string `json:"status"`
    Rescan   bool   `json:"rescan"`
    Carrier  string `json:"carrier"`
    Tracking string `json:"tracking_no"`
}

type Handler struct {
    Dedupe          dedupe.Store
    Parser          qr.Parser
    Pusher          line.Pusher
    SharedVerifier  security.SharedTokenVerifier
    SigVerifier     security.SignatureVerifier
    Logger          logging.Logger
}

func (h *Handler) Register(r chi.Router) {
    r.Post("/api/scan", h.handleScan())
    r.Post("/callback", h.handleCallback())
}

func (h *Handler) handleScan() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if err := h.SharedVerifier.Verify(r); err != nil {
            http.Error(w, "forbidden", http.StatusForbidden)
            return
        }
        var req ScanRequest
        dec := json.NewDecoder(r.Body)
        dec.DisallowUnknownFields()
        if err := dec.Decode(&req); err != nil {
            http.Error(w, "bad request", http.StatusBadRequest)
            return
        }
        if req.GroupID == "" || req.QRText == "" || req.DisplayName == "" {
            http.Error(w, "missing fields", http.StatusBadRequest)
            return
        }
        carrier, tracking := h.Parser.Parse(req.QRText)
        key := carrier + ":" + tracking
        if tracking == "" {
            // Avoid global collision for unknown / no-tracking cases
            sum := sha256.Sum256([]byte(req.QRText))
            key = carrier + ":" + hex.EncodeToString(sum[:8])
        }
        isDup := h.Dedupe.Seen(key)

        // assemble message
        tag := ""
        if isDup {
            tag = "（再スキャン）"
        }
        lines := []string{"配送スキャン" + tag}
        if carrier != "" && carrier != "unknown" {
            lines = append(lines, "配送会社: "+carrier)
        }
        if tracking != "" {
            lines = append(lines, "伝票番号: "+tracking)
        }
        lines = append(lines, "送信者: "+req.DisplayName)
        lines = append(lines, "QR: "+req.QRText)
        message := strings.Join(lines, "\n")

        if err := h.Pusher.Push(req.GroupID, message); err != nil {
            h.Logger.Error(map[string]any{"event": "push_fail", "error": err.Error()})
            http.Error(w, "failed to push", http.StatusBadGateway)
            return
        }

        h.Logger.Info(map[string]any{
            "event": "scan", "group_id": req.GroupID, "carrier": carrier, "tracking_no": tracking, "dedupe": isDup,
        })

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(ScanResponse{
            Status:   "ok",
            Rescan:   isDup,
            Carrier:  carrier,
            Tracking: tracking,
        })
    }
}

func (h *Handler) handleCallback() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Read body for signature verification
        // 1MB cap
        body := http.MaxBytesReader(w, r.Body, 1<<20) 
        b, err := io.ReadAll(body)
        if err != nil {
            http.Error(w, "bad request", http.StatusBadRequest)
            return
        }
        sig := r.Header.Get("X-Line-Signature")
        if err := h.SigVerifier.Verify(sig, b); err != nil {
            if err.Error() == "not configured" {
                http.Error(w, "not configured", http.StatusServiceUnavailable)
                return
            }
            http.Error(w, "forbidden", http.StatusForbidden)
            return
        }
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("ok"))
    }
}
