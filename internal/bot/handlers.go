package bot

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"aislaydbot/internal/i18n"
	"aislaydbot/internal/pptx"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleUpdate(ctx context.Context, upd tgbotapi.Update) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic handler: %v", r)
		}
	}()

	switch {
	case upd.CallbackQuery != nil:
		b.handleCallback(ctx, upd.CallbackQuery)
	case upd.Message != nil:
		b.handleMessage(ctx, upd.Message)
	}
}

func (b *Bot) handleMessage(ctx context.Context, m *tgbotapi.Message) {
	if m.IsCommand() {
		b.handleCommand(ctx, m)
		return
	}
	sess := b.store.Get(m.Chat.ID)
	switch sess.Step {
	case StepAwaitingTopic:
		b.acceptTopic(m, sess)
	case StepAwaitingCount:
		b.acceptCustomCount(m, sess)
	default:
		b.sendMarkdown(m.Chat.ID, i18n.T(sess.Lang).HelpText, mainMenuKeyboard(sess.Lang))
	}
}

func (b *Bot) handleCommand(ctx context.Context, m *tgbotapi.Message) {
	sess := b.store.Get(m.Chat.ID)
	s := i18n.T(sess.Lang)

	switch m.Command() {
	case "start":
		b.sendMarkdown(m.Chat.ID, s.Welcome, mainMenuKeyboard(sess.Lang))
		b.sendMarkdown(m.Chat.ID, s.ChooseLang, langKeyboard())
	case "lang":
		b.sendMarkdown(m.Chat.ID, s.ChooseLang, langKeyboard())
	case "new":
		b.startNew(m.Chat.ID, sess)
	case "cancel":
		b.store.Reset(m.Chat.ID)
		b.sendMarkdown(m.Chat.ID, s.Cancelled, mainMenuKeyboard(sess.Lang))
	case "help":
		b.sendMarkdown(m.Chat.ID, s.HelpText, mainMenuKeyboard(sess.Lang))
	default:
		b.sendMarkdown(m.Chat.ID, s.HelpText, mainMenuKeyboard(sess.Lang))
	}
}

func (b *Bot) handleCallback(ctx context.Context, cq *tgbotapi.CallbackQuery) {
	b.api.Request(tgbotapi.NewCallback(cq.ID, ""))

	chatID := cq.Message.Chat.ID
	sess := b.store.Get(chatID)
	s := i18n.T(sess.Lang)

	data := cq.Data
	switch {
	case strings.HasPrefix(data, "lang:"):
		code := strings.TrimPrefix(data, "lang:")
		newLang := i18n.ParseLang(code)
		b.store.Update(chatID, func(ss *Session) { ss.Lang = newLang })
		b.editText(chatID, cq.Message.MessageID, i18n.T(newLang).LangChanged)
		b.sendMarkdown(chatID, i18n.T(newLang).HelpText, mainMenuKeyboard(newLang))

	case strings.HasPrefix(data, "action:"):
		switch strings.TrimPrefix(data, "action:") {
		case "new":
			b.startNew(chatID, sess)
		case "lang":
			b.sendMarkdown(chatID, s.ChooseLang, langKeyboard())
		case "cancel":
			b.store.Reset(chatID)
			b.editText(chatID, cq.Message.MessageID, s.Cancelled)
		}

	case strings.HasPrefix(data, "count:"):
		if sess.Step != StepAwaitingCount {
			return
		}
		n, err := strconv.Atoi(strings.TrimPrefix(data, "count:"))
		if err != nil || n < 3 || n > 25 {
			b.sendMarkdown(chatID, s.InvalidCount, nil)
			return
		}
		b.store.Update(chatID, func(ss *Session) {
			ss.Count = n
			ss.Step = StepAwaitingTheme
		})
		b.editText(chatID, cq.Message.MessageID, fmt.Sprintf("✔️ %d", n))
		b.sendMarkdown(chatID, s.AskTheme, themeKeyboard(sess.Lang))

	case strings.HasPrefix(data, "theme:"):
		if sess.Step != StepAwaitingTheme {
			return
		}
		id := pptx.ThemeID(strings.TrimPrefix(data, "theme:"))
		b.store.Update(chatID, func(ss *Session) {
			ss.Theme = id
			ss.Step = StepGenerating
		})
		b.editText(chatID, cq.Message.MessageID, "✔️ "+pptx.ThemeLabel(id, sess.Lang))
		b.generateAndSend(ctx, chatID, sess.Lang, sess.Topic, sess.Count, id)
	}
}

func (b *Bot) startNew(chatID int64, sess *Session) {
	s := i18n.T(sess.Lang)
	b.store.Update(chatID, func(ss *Session) {
		ss.Step = StepAwaitingTopic
		ss.Topic = ""
		ss.Count = 0
		ss.Theme = ""
	})
	b.sendMarkdown(chatID, s.AskTopic, nil)
}

func (b *Bot) acceptTopic(m *tgbotapi.Message, sess *Session) {
	s := i18n.T(sess.Lang)
	topic := strings.TrimSpace(m.Text)
	count := utf8.RuneCountInString(topic)
	if count < 3 {
		b.sendMarkdown(m.Chat.ID, s.InvalidTopic, nil)
		return
	}
	if count > 300 {
		b.sendMarkdown(m.Chat.ID, s.TopicTooLong, nil)
		return
	}
	b.store.Update(m.Chat.ID, func(ss *Session) {
		ss.Topic = topic
		ss.Step = StepAwaitingCount
	})
	b.sendMarkdown(m.Chat.ID, s.AskSlideCount, countKeyboard(sess.Lang))
}

func (b *Bot) acceptCustomCount(m *tgbotapi.Message, sess *Session) {
	s := i18n.T(sess.Lang)
	n, err := strconv.Atoi(strings.TrimSpace(m.Text))
	if err != nil || n < 3 || n > 25 {
		b.sendMarkdown(m.Chat.ID, s.InvalidCount, countKeyboard(sess.Lang))
		return
	}
	b.store.Update(m.Chat.ID, func(ss *Session) {
		ss.Count = n
		ss.Step = StepAwaitingTheme
	})
	b.sendMarkdown(m.Chat.ID, s.AskTheme, themeKeyboard(sess.Lang))
}

func (b *Bot) generateAndSend(ctx context.Context, chatID int64, lang i18n.Lang, topic string, count int, themeID pptx.ThemeID) {
	s := i18n.T(lang)

	progressMsg, _ := b.api.Send(tgbotapi.NewMessage(chatID, s.Generating))

	prompt := fmt.Sprintf(s.AIPrompt, topic, count, count, count)

	pres, err := b.ai.Generate(ctx, prompt, count)
	if err != nil {
		log.Printf("AI generate error: %v", err)
		b.deleteMessage(chatID, progressMsg.MessageID)
		b.sendPlain(chatID, fmt.Sprintf(s.Error, "AI: "+err.Error()), mainMenuKeyboard(lang))
		b.store.Reset(chatID)
		return
	}

	pptxPath, err := b.pptxGen.Generate(pres, themeID, topic)
	if err != nil {
		log.Printf("PPTX generate error: %v", err)
		b.deleteMessage(chatID, progressMsg.MessageID)
		b.sendPlain(chatID, fmt.Sprintf(s.Error, "PPTX: "+err.Error()), mainMenuKeyboard(lang))
		b.store.Reset(chatID)
		return
	}

	b.deleteMessage(chatID, progressMsg.MessageID)

	doc := tgbotapi.NewDocument(chatID, tgbotapi.FilePath(pptxPath))
	doc.Caption = fmt.Sprintf(s.SlideCaption, pres.Title, count, pptx.ThemeLabel(themeID, lang))
	doc.ParseMode = "Markdown"
	if _, err := b.api.Send(doc); err != nil {
		log.Printf("send doc error: %v", err)
		b.sendPlain(chatID, fmt.Sprintf(s.Error, err.Error()), mainMenuKeyboard(lang))
	} else {
		b.sendMarkdown(chatID, s.Done, mainMenuKeyboard(lang))
	}

	_ = os.Remove(pptxPath)
	b.store.Reset(chatID)
}

func (b *Bot) sendMarkdown(chatID int64, text string, kb interface{}) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.DisableWebPagePreview = true
	if kb != nil {
		msg.ReplyMarkup = kb
	}
	if _, err := b.api.Send(msg); err != nil {
		log.Printf("send error: %v, retrying without markdown", err)
		b.sendPlain(chatID, text, kb)
	}
}

func (b *Bot) sendPlain(chatID int64, text string, kb interface{}) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.DisableWebPagePreview = true
	if kb != nil {
		msg.ReplyMarkup = kb
	}
	b.api.Send(msg)
}

func (b *Bot) editText(chatID int64, msgID int, text string) {
	edit := tgbotapi.NewEditMessageText(chatID, msgID, text)
	if _, err := b.api.Send(edit); err != nil {
		log.Printf("edit error: %v", err)
	}
}

func (b *Bot) deleteMessage(chatID int64, msgID int) {
	if msgID == 0 {
		return
	}
	b.api.Request(tgbotapi.NewDeleteMessage(chatID, msgID))
}
