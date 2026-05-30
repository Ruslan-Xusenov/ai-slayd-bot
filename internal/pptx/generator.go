package pptx

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"aislaydbot/internal/ai"
)

// Slayd o'lchamlari (EMU - English Metric Units), 16:9
const (
	slideW = 12192000
	slideH = 6858000
)

// Generator — PPTX fayllarni yaratadi
type Generator struct {
	tmpDir string
}

func NewGenerator(tmpDir string) *Generator {
	return &Generator{tmpDir: tmpDir}
}

// Generate — taqdimotni PPTX fayliga aylantiradi va saqlaydi
func (g *Generator) Generate(pres *ai.Presentation, themeID ThemeID, outputName string) (string, error) {
	theme := GetTheme(themeID)

	data, err := buildPPTX(pres, theme)
	if err != nil {
		return "", fmt.Errorf("pptx build: %w", err)
	}

	safe := sanitize(outputName)
	if safe == "" {
		safe = "slayd"
	}
	stamp := time.Now().Format("20060102_150405")
	path := filepath.Join(g.tmpDir, fmt.Sprintf("%s_%s.pptx", safe, stamp))

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", fmt.Errorf("pptx saqlash: %w", err)
	}
	return path, nil
}

// buildPPTX — barcha XML fayllarni ZIP arxivga yig'adi
func buildPPTX(pres *ai.Presentation, t Theme) ([]byte, error) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	slides := collectSlides(pres)
	total := len(slides)

	entries := map[string]string{
		"[Content_Types].xml":                        contentTypesXML(total),
		"_rels/.rels":                                rootRelsXML(),
		"docProps/app.xml":                           appXML(),
		"docProps/core.xml":                          coreXML(pres.Title),
		"ppt/presentation.xml":                       presentationXML(total),
		"ppt/_rels/presentation.xml.rels":            presentationRelsXML(total),
		"ppt/presProps.xml":                          presPropsXML(),
		"ppt/theme/theme1.xml":                       theme1XML(),
		"ppt/slideMasters/slideMaster1.xml":          slideMasterXML(),
		"ppt/slideMasters/_rels/slideMaster1.xml.rels": slideMasterRelsXML(),
		"ppt/slideLayouts/slideLayout1.xml":          slideLayoutXML(),
		"ppt/slideLayouts/_rels/slideLayout1.xml.rels": slideLayoutRelsXML(),
	}

	for name, content := range entries {
		if err := writeEntry(zw, name, content); err != nil {
			return nil, fmt.Errorf("entry %s: %w", name, err)
		}
	}

	for i, s := range slides {
		n := i + 1
		xml := s.render(t, pres.Title, n, total)
		if err := writeEntry(zw, fmt.Sprintf("ppt/slides/slide%d.xml", n), xml); err != nil {
			return nil, err
		}
		if err := writeEntry(zw, fmt.Sprintf("ppt/slides/_rels/slide%d.xml.rels", n), slideRelsXML()); err != nil {
			return nil, err
		}
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// -------- Slayd turlari --------

type slideData struct {
	kind     string // "title" | "content" | "end"
	title    string
	subtitle string
	bullets  []string
}

func (s slideData) render(t Theme, presTitle string, idx, total int) string {
	switch s.kind {
	case "title":
		return titleSlideXML(t, s.title, s.subtitle)
	case "end":
		return endSlideXML(t, presTitle)
	default:
		return contentSlideXML(t, s.title, s.bullets, presTitle, idx, total)
	}
}

func collectSlides(pres *ai.Presentation) []slideData {
	out := []slideData{{kind: "title", title: pres.Title, subtitle: pres.Subtitle}}
	for _, sl := range pres.Slides {
		if strings.EqualFold(strings.TrimSpace(pres.Title), strings.TrimSpace(sl.Title)) {
			continue
		}
		out = append(out, slideData{kind: "content", title: sl.Title, bullets: sl.Bullets})
	}
	out = append(out, slideData{kind: "end"})
	return out
}

// -------- XML yordamchi funksiyalar --------

func writeEntry(zw *zip.Writer, name, content string) error {
	w, err := zw.Create(name)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(content))
	return err
}

func xe(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

func sanitize(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == ' ', r == '-', r == '_':
			b.WriteRune('_')
		}
	}
	out := b.String()
	if len(out) > 40 {
		out = out[:40]
	}
	return strings.Trim(out, "_")
}

func trunc(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "…"
}
