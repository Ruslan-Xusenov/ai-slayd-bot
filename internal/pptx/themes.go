package pptx

import "aislaydbot/internal/i18n"

// ThemeID - tema identifikatori
type ThemeID string

const (
	ThemeModernDark  ThemeID = "modern_dark"
	ThemeModernLight ThemeID = "modern_light"
	ThemeMinimal     ThemeID = "minimal"
	ThemeCorporate   ThemeID = "corporate"
	ThemeVibrant     ThemeID = "vibrant"
	ThemeAcademic    ThemeID = "academic"
)

// Theme - PPTX uchun rang palitrasini saqlaydi (hex RRGGBB)
type Theme struct {
	ID         ThemeID
	Background string // fon rangi (asosiy)
	BgGradient string // gradient oxiri (bo'sh = gradient yo'q)
	Title      string // sarlavha rangi
	Subtitle   string // kichik sarlavha rangi
	Body       string // asosiy matn rangi
	Accent     string // accent rangi (chiziqlar, bullet)
	Footer     string // footer matn rangi
	HeaderBar  string // tepa chiziq rangi
}

func AllThemes() []Theme {
	return []Theme{
		{
			ID:         ThemeModernDark,
			Background: "121624",
			BgGradient: "202040",
			Title:      "FFFFFF",
			Subtitle:   "B4BEDC",
			Body:       "E1E6F5",
			Accent:     "63B3FF",
			Footer:     "8C96B4",
			HeaderBar:  "63B3FF",
		},
		{
			ID:         ThemeModernLight,
			Background: "F8FAFF",
			BgGradient: "E6EEFA",
			Title:      "141C32",
			Subtitle:   "505A78",
			Body:       "283046",
			Accent:     "2563EB",
			Footer:     "788296",
			HeaderBar:  "2563EB",
		},
		{
			ID:         ThemeMinimal,
			Background: "FFFFFF",
			BgGradient: "",
			Title:      "0A0A0A",
			Subtitle:   "5A5A5A",
			Body:       "1E1E1E",
			Accent:     "0A0A0A",
			Footer:     "8C8C8C",
			HeaderBar:  "0A0A0A",
		},
		{
			ID:         ThemeCorporate,
			Background: "F5F7FA",
			BgGradient: "E1E8F0",
			Title:      "0F2346",
			Subtitle:   "3C506E",
			Body:       "1E2D46",
			Accent:     "005AA0",
			Footer:     "647896",
			HeaderBar:  "005AA0",
		},
		{
			ID:         ThemeVibrant,
			Background: "FF5E82",
			BgGradient: "7850DC",
			Title:      "FFFFFF",
			Subtitle:   "FFE6F0",
			Body:       "FFFFFF",
			Accent:     "FFDC64",
			Footer:     "FFDCEB",
			HeaderBar:  "FFDC64",
		},
		{
			ID:         ThemeAcademic,
			Background: "FCFAF4",
			BgGradient: "F0EADC",
			Title:      "3C1E14",
			Subtitle:   "6E503C",
			Body:       "2D1E14",
			Accent:     "A04628",
			Footer:     "786450",
			HeaderBar:  "A04628",
		},
	}
}

func GetTheme(id ThemeID) Theme {
	for _, t := range AllThemes() {
		if t.ID == id {
			return t
		}
	}
	return AllThemes()[0]
}

func ThemeLabel(id ThemeID, lang i18n.Lang) string {
	s := i18n.T(lang)
	switch id {
	case ThemeModernDark:
		return s.ThemeModernDark
	case ThemeModernLight:
		return s.ThemeModernLight
	case ThemeMinimal:
		return s.ThemeMinimal
	case ThemeCorporate:
		return s.ThemeCorporate
	case ThemeVibrant:
		return s.ThemeVibrant
	case ThemeAcademic:
		return s.ThemeAcademic
	}
	return string(id)
}
