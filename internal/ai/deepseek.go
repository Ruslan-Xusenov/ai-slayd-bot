package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const maxRetries = 3

const deepseekBaseURL = "https://api.deepseek.com/v1"
const deepseekModel = "deepseek-chat"

// deepseekMessage — bitta chat xabari
type deepseekMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// deepseekRequest — DeepSeek API so'rovi
type deepseekRequest struct {
	Model          string            `json:"model"`
	Messages       []deepseekMessage `json:"messages"`
	MaxTokens      int               `json:"max_tokens"`
	Temperature    float64           `json:"temperature"`
	ResponseFormat map[string]string `json:"response_format"`
}

// deepseekChoice — javobdagi tanlov
type deepseekChoice struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
}

// deepseekResponse — DeepSeek API javobi
type deepseekResponse struct {
	Choices []deepseekChoice `json:"choices"`
	Error   *struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"error,omitempty"`
}

// DeepSeekClient — DeepSeek AI uchun client
type DeepSeekClient struct {
	apiKey     string
	httpClient *http.Client
}

// NewDeepSeek — yangi DeepSeek client yaratadi
func NewDeepSeek(apiKey string) *DeepSeekClient {
	return &DeepSeekClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// Generate — DeepSeek orqali taqdimot yaratadi (retry bilan)
func (c *DeepSeekClient) Generate(ctx context.Context, prompt string, slideCount int) (*Presentation, error) {
	log.Printf("🤖 DeepSeek model ishlatilmoqda: %s", deepseekModel)

	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		if attempt > 1 {
			waitSec := time.Duration(attempt*attempt) * time.Second
			log.Printf("🔄 Qayta urinish %d/%d (%.0f soniyadan keyin)...", attempt, maxRetries, waitSec.Seconds())
			select {
			case <-time.After(waitSec):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
		result, err := c.doGenerate(ctx, prompt, slideCount)
		if err == nil {
			return result, nil
		}
		lastErr = err
		log.Printf("⚠️  Urinish %d xatosi: %v", attempt, err)
	}
	return nil, fmt.Errorf("%d urinishdan keyin ham xato: %w", maxRetries, lastErr)
}

// doGenerate — bir marta API so'rov yuboradi
func (c *DeepSeekClient) doGenerate(ctx context.Context, prompt string, slideCount int) (*Presentation, error) {

	systemPrompt := `You are an expert presentation creator. Generate detailed, information-rich slides. ` +
		`Each bullet point must be a long, complete, meaningful sentence (15-30 words minimum). ` +
		`Include specific facts, detailed explanations, statistics, and examples. Never write short superficial sentences. ` +
		`Reply ONLY with valid JSON. Do NOT include any comments (like // or /*) inside the JSON. ` +
		`Do not write any conversational text before or after the JSON. Start with { and end with }.` +
		fmt.Sprintf(`
Generate exactly %d slides. The JSON structure must be:
{
  "title": "Presentation Title",
  "subtitle": "Subtitle",
  "slides": [
    {
      "title": "Slide Title",
      "bullets": ["bullet 1", "bullet 2", "bullet 3"],
      "note": "Speaker note"
    }
  ]
}`, slideCount)

	reqBody := deepseekRequest{
		Model: deepseekModel,
		Messages: []deepseekMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: prompt},
		},
		MaxTokens:   8000,
		Temperature: 0.7,
		ResponseFormat: map[string]string{
			"type": "json_object",
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("so'rov yaratishda xato: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		deepseekBaseURL+"/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("HTTP so'rov yaratishda xato: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("DeepSeek API xatosi: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("javobni o'qishda xato: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DeepSeek HTTP %d: %s", resp.StatusCode, string(respBytes))
	}

	var dsResp deepseekResponse
	if err := json.Unmarshal(respBytes, &dsResp); err != nil {
		return nil, fmt.Errorf("JSON parse xatosi: %w | javob: %s", err, truncate(string(respBytes), 300))
	}

	if dsResp.Error != nil {
		return nil, fmt.Errorf("DeepSeek xatosi (kod %d): %s", dsResp.Error.Code, dsResp.Error.Message)
	}

	if len(dsResp.Choices) == 0 {
		return nil, errors.New("bo'sh javob")
	}

	raw := strings.TrimSpace(dsResp.Choices[0].Message.Content)
	raw = stripFences(raw)

	// Agar JSON kesilgan bo'lsa, tuzatishga harakat qilamiz
	raw = fixTruncatedJSON(raw)

	var p Presentation
	if err := json.Unmarshal([]byte(raw), &p); err != nil {
		return nil, fmt.Errorf("taqdimot JSON parse xatosi: %w | javob: %s", err, truncate(raw, 300))
	}
	if len(p.Slides) == 0 {
		return nil, errors.New("slaydlar bo'sh")
	}
	if len(p.Slides) > slideCount {
		p.Slides = p.Slides[:slideCount]
	}
	return &p, nil
}

// fixTruncatedJSON — kesilgan JSONni tuzatishga harakat qiladi
func fixTruncatedJSON(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	// Agar allaqachon to'g'ri bo'lsa, o'zgartirmaymiz
	var tmp interface{}
	if json.Unmarshal([]byte(s), &tmp) == nil {
		return s
	}

	// Ochilmagan qo'shtirnoqni yopamiz
	inString := false
	escaped := false
	for _, ch := range s {
		if escaped {
			escaped = false
			continue
		}
		if ch == '\\' {
			escaped = true
			continue
		}
		if ch == '"' {
			inString = !inString
		}
	}
	if inString {
		s += "\""
	}

	// Yopilmagan ob'ekt/massivlarni yopamiz
	opens := 0
	squares := 0
	inStr2 := false
	esc2 := false
	for _, ch := range s {
		if esc2 {
			esc2 = false
			continue
		}
		if ch == '\\' && inStr2 {
			esc2 = true
			continue
		}
		if ch == '"' {
			inStr2 = !inStr2
			continue
		}
		if inStr2 {
			continue
		}
		switch ch {
		case '{':
			opens++
		case '}':
			opens--
		case '[':
			squares++
		case ']':
			squares--
		}
	}
	// Yetishmayotgan qavslarni qo'shamiz
	for squares > 0 {
		s += "]"
		squares--
	}
	for opens > 0 {
		s += "}"
		opens--
	}

	return s
}