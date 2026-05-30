package bot

import (
	"fmt"

	"aislaydbot/internal/i18n"
	"aislaydbot/internal/pptx"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func langKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇺🇿 O'zbekcha", "lang:uz"),
			tgbotapi.NewInlineKeyboardButtonData("🇷🇺 Русский", "lang:ru"),
			tgbotapi.NewInlineKeyboardButtonData("🇬🇧 English", "lang:en"),
		),
	)
}

func countKeyboard(lang i18n.Lang) tgbotapi.InlineKeyboardMarkup {
	s := i18n.T(lang)
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("3", "count:3"),
			tgbotapi.NewInlineKeyboardButtonData("5", "count:5"),
			tgbotapi.NewInlineKeyboardButtonData("7", "count:7"),
			tgbotapi.NewInlineKeyboardButtonData("10", "count:10"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("12", "count:12"),
			tgbotapi.NewInlineKeyboardButtonData("15", "count:15"),
			tgbotapi.NewInlineKeyboardButtonData("20", "count:20"),
			tgbotapi.NewInlineKeyboardButtonData("25", "count:25"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(s.CancelBtn, "action:cancel"),
		),
	)
}

func themeKeyboard(lang i18n.Lang) tgbotapi.InlineKeyboardMarkup {
	rows := [][]tgbotapi.InlineKeyboardButton{}
	themes := pptx.AllThemes()
	for i := 0; i < len(themes); i += 2 {
		row := []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(
				pptx.ThemeLabel(themes[i].ID, lang),
				fmt.Sprintf("theme:%s", themes[i].ID),
			),
		}
		if i+1 < len(themes) {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(
				pptx.ThemeLabel(themes[i+1].ID, lang),
				fmt.Sprintf("theme:%s", themes[i+1].ID),
			))
		}
		rows = append(rows, row)
	}
	s := i18n.T(lang)
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(s.CancelBtn, "action:cancel"),
	))
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func mainMenuKeyboard(lang i18n.Lang) tgbotapi.InlineKeyboardMarkup {
	s := i18n.T(lang)
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(s.NewSlideBtn, "action:new"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(s.ChangeLangBtn, "action:lang"),
		),
	)
}
