package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	DeepSeekKey   string
	TmpDir        string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		TelegramToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		DeepSeekKey:   os.Getenv("DEEPSEEK_API_KEY"),
		TmpDir:        os.Getenv("TMP_DIR"),
	}

	if cfg.TelegramToken == "" {
		return nil, errors.New("TELEGRAM_BOT_TOKEN bo'sh — .env faylga qo'shing")
	}
	if cfg.DeepSeekKey == "" {
		return nil, errors.New("DEEPSEEK_API_KEY bo'sh — .env faylga qo'shing")
	}
	if cfg.TmpDir == "" {
		cfg.TmpDir = "./tmp"
	}

	if err := os.MkdirAll(cfg.TmpDir, 0o755); err != nil {
		return nil, err
	}

	return cfg, nil
}
