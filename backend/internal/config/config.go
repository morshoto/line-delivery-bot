package config

import "os"

type Config struct {
    Port                  string
    SharedToken           string
    LineChannelSecret     string
    LineChannelAccessToken string
}

func FromEnv() Config {
    port := os.Getenv("PORT")
    if port == "" {
        port = "10000"
    }
    return Config{
        Port:                   port,
        SharedToken:            os.Getenv("SHARED_TOKEN"),
        LineChannelSecret:      os.Getenv("LINE_CHANNEL_SECRET"),
        LineChannelAccessToken: os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
    }
}

