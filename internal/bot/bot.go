package bot

import (
	"context"
	"log"

	"aislaydbot/internal/ai"
	"aislaydbot/internal/config"
	"aislaydbot/internal/pptx"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api     *tgbotapi.BotAPI
	cfg     *config.Config
	store   *Store
	ai      *ai.DeepSeekClient
	pptxGen *pptx.Generator
}

func New(cfg *config.Config) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		return nil, err
	}
	api.Debug = false

	log.Printf("✅ Bot ulandi: @%s", api.Self.UserName)

	return &Bot{
		api:     api,
		cfg:     cfg,
		store:   NewStore(),
		ai:      ai.NewDeepSeek(cfg.DeepSeekKey),
		pptxGen: pptx.NewGenerator(cfg.TmpDir),
	}, nil
}

func (b *Bot) Run(ctx context.Context) error {
	if _, err := b.api.Request(tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{Command: "start", Description: "Start the bot / Botni ishga tushirish"},
		tgbotapi.BotCommand{Command: "new", Description: "Create new presentation / Yangi taqdimot"},
		tgbotapi.BotCommand{Command: "lang", Description: "Change language / Tilni o'zgartirish"},
		tgbotapi.BotCommand{Command: "cancel", Description: "Cancel current action / Bekor qilish"},
		tgbotapi.BotCommand{Command: "help", Description: "Help / Yordam"},
	)); err != nil {
		log.Printf("commands set xatosi: %v", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			b.api.StopReceivingUpdates()
			return ctx.Err()
		case upd, ok := <-updates:
			if !ok {
				return nil
			}
			go b.handleUpdate(ctx, upd)
		}
	}
}
