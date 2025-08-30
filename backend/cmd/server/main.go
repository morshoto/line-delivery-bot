package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"

    "example.com/delivery-bot/server/internal/config"
    "example.com/delivery-bot/server/internal/dedupe"
    "example.com/delivery-bot/server/internal/httpapi"
    "example.com/delivery-bot/server/internal/line"
    "example.com/delivery-bot/server/internal/logging"
    "example.com/delivery-bot/server/internal/qr"
    "example.com/delivery-bot/server/internal/security"
)

func main() {
    cfg := config.FromEnv()

    r := chi.NewRouter()
    r.Use(middleware.RealIP, middleware.Logger, middleware.Recoverer)

    // Health check
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("ok"))
    })

    logger := logging.JSONLogger{}
    store := dedupe.NewTTLStore(30 * time.Minute)
    parser := qr.NaiveParser{}

    var pusher line.Pusher
    if cfg.LineChannelAccessToken != "" {
        pusher = line.NewHTTPPusher(cfg.LineChannelAccessToken)
    } else {
        logger.Warn(map[string]any{"event": "push_skip", "reason": "missing_token"})
        pusher = line.NoopPusher{}
    }

    sharedVerifier := security.NewHeaderSharedTokenVerifier(cfg.SharedToken)
    sigVerifier := security.NewLineSignatureVerifier(cfg.LineChannelSecret)

    handler := httpapi.Handler{
        Dedupe:         store,
        Parser:         parser,
        Pusher:         pusher,
        SharedVerifier: sharedVerifier,
        SigVerifier:    sigVerifier,
        Logger:         logger,
    }
    handler.Register(r)

    addr := ":" + cfg.Port

    srv := &http.Server{
        Addr:         addr,
        Handler:      r,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    // Start server in background
    go func() {
        log.Printf("listening on %s", addr)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("server error: %v", err)
        }
    }()

    // Handle graceful shutdown on SIGINT/SIGTERM
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
    <-stop
    log.Printf("shutdown signal received, draining...")

    ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Printf("graceful shutdown failed: %v", err)
        if cerr := srv.Close(); cerr != nil {
            log.Printf("force close failed: %v", cerr)
        }
    }
    log.Printf("server stopped")
}
