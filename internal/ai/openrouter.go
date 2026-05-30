package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	openrouter "github.com/OpenRouterTeam/go-sdk"
	"github.com/OpenRouterTeam/go-sdk/models/components"
	"github.com/OpenRouterTeam/go-sdk/optionalnullable"
)

// Qat'iy bitta model
const modelName = "anthropic/claude-3.5-haiku"

// Slide — bitta slayd ma'lumotlari
type Slide struct {
	Title   string   `json:"title"`
	Bullets []string `json:"bullets"`
	Note    string   `json:"note"`
}

// Presentation — to'liq taqdimot
type Presentation struct {
	Title    string  `json:"title"`
	Subtitle string  `json:"subtitle"`
	Slides   []Slide `json:"slides"`
}

// Client — faqatgina Llama 3.2 3B Free model uchun AI client
type Client struct {
	sdk *openrouter.OpenRouter
}

// New — yangi AI client yaratadi
func New(apiKey, _ string) *Client {
	return &Client{
		sdk: openrouter.New(openrouter.WithSecurity(apiKey)),
	}
}

// Generate — faqatgina bitta model orqali taqdimot yaratadi (hech qanday zaxirasiz)
func (c *Client) Generate(ctx context.Context, prompt string, slideCount int) (*Presentation, error) {
	log.Printf("🤖 Model ishlatilmoqda: %s", modelName)

	aiCtx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	pres, err := c.callModel(aiCtx, modelName, prompt, slideCount)
	if err != nil {
		return nil, fmt.Errorf("API xatosi: %w", err)
	}
	return pres, nil
}

// callModel — berilgan model orqali API so'rov yuboradi
func (c *Client) callModel(ctx context.Context, model string, prompt string, slideCount int) (*Presentation, error) {
	systemMsg := components.CreateChatMessagesSystem(components.ChatSystemMessage{
		Role: components.ChatSystemMessageRoleSystem,
		Content: components.CreateChatSystemMessageContentStr(
			`You are an expert presentation creator. Generate detailed, information-rich slides. ` +
				`Each bullet point must be a long, complete, meaningful sentence (15-30 words minimum). ` +
				`Include specific facts, detailed explanations, statistics, and examples. Never write short superficial sentences. ` +
				`Reply ONLY with valid JSON. Do NOT include any comments (like // or /*) inside the JSON. ` +
				`Do not write any conversational text before or after the JSON. Start with { and end with }.`,
		),
	})

	userMsg := components.CreateChatMessagesUser(components.ChatUserMessage{
		Role: components.ChatUserMessageRoleUser,
		Content: components.CreateChatUserMessageContentStr(prompt),
	})

	// Claude 3.5 Haiku JSON formatni rasman qo'llab-quvvatlaydi, shuning uchun uni yuboramiz
	respFormat := &components.ResponseFormat{
		FormatJSONObjectConfig: &components.FormatJSONObjectConfig{
			Type: components.FormatJSONObjectConfigTypeJSONObject,
		},
		Type: components.ResponseFormatTypeJSONObject,
	}

	maxTok := int64(2000)

	resp, err := c.sdk.Chat.Send(ctx, components.ChatRequest{
		Model:          openrouter.Pointer(modelName),
		MaxTokens:      optionalnullable.From(&maxTok),
		Messages:       []components.ChatMessages{systemMsg, userMsg},
		ResponseFormat: respFormat,
	})

	if err != nil {
		return nil, fmt.Errorf("API so'rov (%s): %w", modelName, err)
	}

	if resp.ChatResult == nil || len(resp.ChatResult.Choices) == 0 {
		return nil, errors.New("bo'sh javob")
	}

	raw := ""
	msgContent, ok := resp.ChatResult.Choices[0].Message.Content.GetOrZero()
	if !ok {
		return nil, errors.New("bo'sh content")
	}
	if msgContent.Str != nil {
		raw = strings.TrimSpace(*msgContent.Str)
	} else {
		return nil, errors.New("javob matn emas")
	}

	raw = stripFences(raw)

	var p Presentation
	if err := json.Unmarshal([]byte(raw), &p); err != nil {
		return nil, fmt.Errorf("JSON parse xatosi: %w | javob: %s", err, truncate(raw, 200))
	}
	if len(p.Slides) == 0 {
		return nil, errors.New("slaydlar bo'sh")
	}
	if len(p.Slides) > slideCount {
		p.Slides = p.Slides[:slideCount]
	}
	return &p, nil
}

func stripFences(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```") {
		if i := strings.Index(s, "\n"); i != -1 {
			s = s[i+1:]
		}
		s = strings.TrimSuffix(s, "```")
		s = strings.TrimSpace(s)
	}
	if start := strings.Index(s, "{"); start > 0 {
		s = s[start:]
	}
	if end := strings.LastIndex(s, "}"); end >= 0 && end < len(s)-1 {
		s = s[:end+1]
	}
	return s
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
