package main

import (
	"fmt"
	"log"
	"os"

	"aislaydbot/internal/ai"
	"aislaydbot/internal/pptx"
)

func main() {
	pres := &ai.Presentation{
		Title:    "Sun'iy intellekt va kelajak",
		Subtitle: "Yangi davrning asosiy texnologiyasi",
		Slides: []ai.Slide{
			{
				Title:   "Sun'iy intellekt va kelajak",
				Bullets: []string{"Kirish slayd"},
			},
			{
				Title: "Sun'iy intellekt nima?",
				Bullets: []string{
					"Inson tafakkuriga taqlid qiluvchi tizimlar",
					"Mashina o'rganish (Machine Learning) asosiy yo'nalish",
					"Chuqur o'rganish (Deep Learning) neyron tarmoqlar",
					"Generativ AI: matn, rasm, video yaratish",
				},
			},
			{
				Title: "AI ning afzalliklari",
				Bullets: []string{
					"Katta hajmdagi ma'lumotlarni tez tahlil qilish",
					"24/7 to'xtovsiz ishlash imkoniyati",
					"Inson xatolarini minimumga tushirish",
					"Murakkab vazifalarni avtomatlashtirish",
				},
			},
			{
				Title: "Asosiy qo'llanish sohalari",
				Bullets: []string{
					"Tibbiyot: diagnostika va davolash",
					"Ta'lim: shaxsiy o'qituvchi yordamchilari",
					"Moliya: firibgarlikni aniqlash",
					"Transport: avtonom haydash",
				},
			},
			{
				Title: "Xulosa",
				Bullets: []string{
					"AI hayotimizning ajralmas qismiga aylanmoqda",
					"Yangi imkoniyatlar va vazifalar paydo bo'lmoqda",
					"O'rganishni boshlash uchun eng yaxshi vaqt — hozir",
				},
			},
		},
	}

	tmpDir := "./tmp"
	_ = os.MkdirAll(tmpDir, 0o755)

	gen := pptx.NewGenerator(tmpDir)

	themes := []pptx.ThemeID{
		pptx.ThemeModernDark,
		pptx.ThemeModernLight,
		pptx.ThemeMinimal,
		pptx.ThemeCorporate,
		pptx.ThemeVibrant,
		pptx.ThemeAcademic,
	}

	for _, th := range themes {
		path, err := gen.Generate(pres, th, fmt.Sprintf("test_%s", th))
		if err != nil {
			log.Fatalf("[%s] xato: %v", th, err)
		}
		log.Printf("✔ %s -> %s", th, path)
	}
}
