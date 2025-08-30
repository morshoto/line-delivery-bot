package main

import (
    "log"
    "net/http"
    "os"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

func main() {
    r := chi.NewRouter()
    r.Use(middleware.RealIP, middleware.Logger, middleware.Recoverer)

    // 健康チェック
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("ok"))
    })

    // 将来のエンドポイントの置き場所（空でOK）
    r.Post("/api/scan", func(w http.ResponseWriter, r *http.Request) {
        http.Error(w, "not implemented", http.StatusNotImplemented)
    })

    // PORT 環境変数（Renderは自動で渡す）
    port := os.Getenv("PORT")
    if port == "" { port = "10000" }
    addr := ":" + port
    log.Printf("listening on %s", addr)
    if err := http.ListenAndServe(addr, r); err != nil {
        log.Fatal(err)
    }
}

