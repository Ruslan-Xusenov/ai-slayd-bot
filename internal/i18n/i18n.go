package i18n

type Lang string

const (
	UZ Lang = "uz"
	RU Lang = "ru"
	EN Lang = "en"
)

type Strings struct {
	Welcome          string
	ChooseLang       string
	LangChanged      string
	AskTopic         string
	AskSlideCount    string
	AskTheme         string
	ChooseFromList   string
	Generating       string
	Done             string
	Error            string
	Cancelled        string
	CancelBtn        string
	NewSlideBtn      string
	ChangeLangBtn    string
	HelpText         string
	InvalidTopic     string
	InvalidCount     string
	TopicTooLong     string
	ThemeModernDark  string
	ThemeModernLight string
	ThemeMinimal     string
	ThemeCorporate   string
	ThemeVibrant     string
	ThemeAcademic    string
	AIPrompt         string
	SlideCaption     string
	BackBtn          string
}

var bundles = map[Lang]Strings{
	UZ: {
		Welcome:          "👋 Salom! Men *AI Slayd Bot*man.\n\nMen sizga istalgan mavzuda chiroyli PDF taqdimotlarni tayyorlab beraman. Yangi slayd yaratish uchun /new buyrug'ini bering.",
		ChooseLang:       "🌐 Tilni tanlang:",
		LangChanged:      "✅ Til o'zbek tiliga o'zgartirildi.",
		AskTopic:         "📝 *Mavzuni yozing:*\n\nMisol: _Sun'iy intellektning tibbiyotdagi o'rni_",
		AskSlideCount:    "📊 *Nechta slayd kerak?*\n\nQuyidagi variantlardan birini tanlang yoki sonni yozing (3-25):",
		AskTheme:         "🎨 *Slayd dizaynini tanlang:*",
		ChooseFromList:   "Iltimos, ro'yxatdan birini tanlang.",
		Generating:       "⏳ AI siz uchun slayd tayyorlamoqda... Bu 30-60 soniya olishi mumkin.",
		Done:             "✅ Slayd tayyor!",
		Error:            "❌ Xatolik yuz berdi: %s\n\nQayta urinib ko'ring: /new",
		Cancelled:        "🚫 Bekor qilindi. Yangi slayd uchun /new",
		CancelBtn:        "❌ Bekor qilish",
		NewSlideBtn:      "🆕 Yangi slayd",
		ChangeLangBtn:    "🌐 Tilni o'zgartirish",
		HelpText:         "📚 *Buyruqlar:*\n/new - Yangi slayd yaratish\n/lang - Tilni o'zgartirish\n/cancel - Joriy jarayonni bekor qilish\n/help - Yordam",
		InvalidTopic:     "⚠️ Mavzu juda qisqa. Iltimos, batafsilroq yozing (kamida 3 ta belgi).",
		InvalidCount:     "⚠️ Slayd soni 3 dan 25 gacha bo'lishi kerak.",
		TopicTooLong:     "⚠️ Mavzu juda uzun (maksimum 300 belgi).",
		ThemeModernDark:  "🌙 Modern Dark",
		ThemeModernLight: "☀️ Modern Light",
		ThemeMinimal:     "⚪ Minimal",
		ThemeCorporate:   "💼 Korporativ",
		ThemeVibrant:     "🌈 Yorqin",
		ThemeAcademic:    "🎓 Akademik",
		AIPrompt: `"%s" mavzusida professional taqdimot uchun aniq %d ta slayd yarat.

FAQAT to'g'ri JSON qaytaring (markdown yoki qo'shimcha matn yo'q).

JSON format:
{
  "title": "Taqdimot sarlavhasi (5-8 so'z)",
  "subtitle": "Mavzuning asosiy mohiyatini 8-15 so'zda ifodalang",
  "slides": [
    {
      "title": "Slayd sarlavhasi (4-7 so'z)",
      "bullets": [
        "Har bir bullet kamida 15-30 so'zdan iborat bo'lib, chuqur va atroflicha yoritilgan ma'lumot, fakt yoki misolni o'z ichiga olishi shart",
        "Raqamlar, foizlar, sanalar, muammolar yoki ularning yechimlarini tushuntirib bering — umuman qisqa yoki yuzaki gaplar ishlatmang",
        "O'quvchi ushbu bulletdan juda qimmatli va to'liq bilim olishi kerak"
      ],
      "note": "Notiq uchun qisqa izoh"
    }
  ]
}

TALABLAR:
- QAT'IY TALAB: "slides" massivi roppa-rosa %d ta slayddan iborat bo'lishi SHART! Agar sizdan %d ta slayd so'ralgan bo'lsa, aniq shuncha slayd yarating.
- Birinchi slayd kirish (titul), oxirgi slayd xulosa
- Har bir slaydda 6-7 ta BATAFSIL VA KENG YORITILGAN bullet
- Bullet ichida qisqa so'zlar EMAS: mavzuni to'liq tushuntiruvchi murakkab va uzun gaplar
- Slayd sarlavhalari qisqa va kuchli
- Til: O'zbek tili (lotin alifbosi)
- FAQAT JSON qaytaring`,
		SlideCaption: "📊 *%s*\n_%d ta slayd · %s_",
		BackBtn:      "⬅️ Orqaga",
	},
	RU: {
		Welcome:          "👋 Привет! Я *AI Slide Bot*.\n\nЯ создам красивую PDF презентацию на любую тему. Чтобы создать новую презентацию, нажмите /new",
		ChooseLang:       "🌐 Выберите язык:",
		LangChanged:      "✅ Язык изменён на русский.",
		AskTopic:         "📝 *Напишите тему:*\n\nПример: _Роль искусственного интеллекта в медицине_",
		AskSlideCount:    "📊 *Сколько слайдов нужно?*\n\nВыберите вариант или введите число (3-25):",
		AskTheme:         "🎨 *Выберите дизайн слайдов:*",
		ChooseFromList:   "Пожалуйста, выберите из списка.",
		Generating:       "⏳ ИИ создаёт презентацию... Это займёт 30-60 секунд.",
		Done:             "✅ Презентация готова!",
		Error:            "❌ Произошла ошибка: %s\n\nПопробуйте снова: /new",
		Cancelled:        "🚫 Отменено. Для новой презентации: /new",
		CancelBtn:        "❌ Отменить",
		NewSlideBtn:      "🆕 Новая презентация",
		ChangeLangBtn:    "🌐 Сменить язык",
		HelpText:         "📚 *Команды:*\n/new - Создать новую презентацию\n/lang - Сменить язык\n/cancel - Отменить текущий процесс\n/help - Помощь",
		InvalidTopic:     "⚠️ Тема слишком короткая. Минимум 3 символа.",
		InvalidCount:     "⚠️ Количество слайдов должно быть от 3 до 25.",
		TopicTooLong:     "⚠️ Тема слишком длинная (максимум 300 символов).",
		ThemeModernDark:  "🌙 Modern Dark",
		ThemeModernLight: "☀️ Modern Light",
		ThemeMinimal:     "⚪ Минимал",
		ThemeCorporate:   "💼 Корпоративный",
		ThemeVibrant:     "🌈 Яркий",
		ThemeAcademic:    "🎓 Академический",
		AIPrompt: `Создай профессиональную презентацию на тему "%s" из ровно %d слайдов.

Возвращай ТОЛЬКО валидный JSON (без markdown и лишнего текста).

Формат JSON:
{
  "title": "Заголовок презентации (5-8 слов)",
  "subtitle": "Основная идея темы в 8-15 словах",
  "slides": [
    {
      "title": "Заголовок слайда (4-7 слов)",
      "bullets": [
        "Каждый буллет — минимум 15-30 слов с подробным объяснением, фактом, цифрой или примером",
        "Указывай статистику, причины, следствия, конкретные результаты — никаких коротких и поверхностных фраз",
        "Читатель должен получить исчерпывающую информацию из каждого буллета"
      ],
      "note": "Краткая заметка для докладчика"
    }
  ]
}

ТРЕБОВАНИЯ:
- СТРОГОЕ ТРЕБОВАНИЕ: Массив "slides" должен содержать ровно %d слайдов! Если запрошено %d слайдов, сгенерируйте ровно столько.
- Первый слайд — вводный, последний — заключение
- 4-5 ПОДРОБНЫХ И РАЗВЕРНУТЫХ буллетов на каждый слайд
- В буллетах НЕ короткие слова: сложные и длинные предложения, полностью раскрывающие суть
- Заголовки слайдов — краткие и ёмкие
- Язык: Русский
- ТОЛЬКО JSON`,
		SlideCaption: "📊 *%s*\n_%d слайдов · %s_",
		BackBtn:      "⬅️ Назад",
	},
	EN: {
		Welcome:          "👋 Hello! I'm *AI Slide Bot*.\n\nI'll create beautiful PDF presentations on any topic. To start, send /new",
		ChooseLang:       "🌐 Choose a language:",
		LangChanged:      "✅ Language changed to English.",
		AskTopic:         "📝 *Enter the topic:*\n\nExample: _The role of AI in healthcare_",
		AskSlideCount:    "📊 *How many slides?*\n\nPick an option below or type a number (3-25):",
		AskTheme:         "🎨 *Choose a slide theme:*",
		ChooseFromList:   "Please choose from the list.",
		Generating:       "⏳ AI is generating your slides... This may take 30-60 seconds.",
		Done:             "✅ Your presentation is ready!",
		Error:            "❌ An error occurred: %s\n\nTry again with /new",
		Cancelled:        "🚫 Cancelled. Use /new to start over.",
		CancelBtn:        "❌ Cancel",
		NewSlideBtn:      "🆕 New presentation",
		ChangeLangBtn:    "🌐 Change language",
		HelpText:         "📚 *Commands:*\n/new - Create a new presentation\n/lang - Change language\n/cancel - Cancel current process\n/help - Help",
		InvalidTopic:     "⚠️ Topic is too short. Please enter at least 3 characters.",
		InvalidCount:     "⚠️ Slide count must be between 3 and 25.",
		TopicTooLong:     "⚠️ Topic is too long (max 300 characters).",
		ThemeModernDark:  "🌙 Modern Dark",
		ThemeModernLight: "☀️ Modern Light",
		ThemeMinimal:     "⚪ Minimal",
		ThemeCorporate:   "💼 Corporate",
		ThemeVibrant:     "🌈 Vibrant",
		ThemeAcademic:    "🎓 Academic",
		AIPrompt: `Create a professional presentation on "%s" with exactly %d slides.

Return ONLY valid JSON (no markdown, no extra text).

JSON format:
{
  "title": "Presentation title (5-8 words)",
  "subtitle": "Core essence of the topic in 8-15 words",
  "slides": [
    {
      "title": "Slide title (4-7 words)",
      "bullets": [
        "Each bullet must be 15-30 words minimum with a detailed explanation, specific fact, statistic, or real-world example",
        "Include deep insights, causes, effects, percentages, names — avoid short or superficial generalities",
        "The reader must gain comprehensive knowledge from every single bullet point"
      ],
      "note": "Brief speaker note"
    }
  ]
}

REQUIREMENTS:
- STRICT REQUIREMENT: The "slides" array MUST contain exactly %d slides! If %d slides are requested, you must generate exactly that amount.
- First slide is introduction, last slide is conclusion
- 4-5 HIGHLY DETAILED bullets per slide
- NO short phrases in bullets: use complex, full sentences that completely cover the topic
- Slide titles must be concise and impactful
- Language: English
- Return ONLY JSON`,
		SlideCaption: "📊 *%s*\n_%d slides · %s_",
		BackBtn:      "⬅️ Back",
	},
}

func T(lang Lang) Strings {
	if s, ok := bundles[lang]; ok {
		return s
	}
	return bundles[EN]
}

func ParseLang(code string) Lang {
	switch code {
	case "uz":
		return UZ
	case "ru":
		return RU
	case "en":
		return EN
	default:
		return EN
	}
}
