package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"aislaydbot/internal/bot"
	"aislaydbot/internal/config"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config xatosi: %v", err)
	}

	b, err := bot.New(cfg)
	if err != nil {
		log.Fatalf("bot yaratish xatosi: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Println("🚀 AI Slayd Bot ishga tushdi (Ctrl+C bilan to'xtatish)")
	if err := b.Run(ctx); err != nil && err != context.Canceled {
		log.Fatalf("bot run xatosi: %v", err)
	}
	log.Println("👋 Bot to'xtatildi")
}
