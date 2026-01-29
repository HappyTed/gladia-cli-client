package prerecorderv2

// Конфигурация языка текста на записи
// Доступные языки: https://docs.gladia.io/api-reference/v2/pre-recorded/init#body-language-config
type LanguageConf struct {
	Languages     []string `json:"languages"`      // Если установлен один язык, он будет использоваться для транскрибации. В противном случае, язык будет автоматически обнаружен моделью.
	CodeSwitching bool     `json:"code_switching"` // Если true - язык будет автоматически обнаружен на каждом высказывании. В противном случае язык будет автоматически обнаружен при первом произнесении, а затем использоваться для остальной части транскрибации. Если установлен один язык, этот вариант будет проигнорирован.
}

// Конфигурация распознавания спикеров, если diarization включено
type DiarizationConf struct {
	Enhanced      bool   `json:"enhanced"`                     //
	NumOfSpeakers *uint8 `json:"number_of_speakers,omitempty"` // количество спикеров, >=1
	MinSpeakers   *uint8 `json:"min_speakers,omitempty"`       // минимальное количество спикеров в моменте
	MaxSpeakers   *uint8 `json:"max_speakers,omitempty"`       // максимальное количество спикеров в моменте
}

// Конфигурация перевода, если  translation включено
type TranslationConf struct {
	Model             *string  `json:"model,omitempty"`              // Модель, которую вы хотите использовать модель перевода для перевода. Доступные опции: base, enhanced
	TargetLanguages   []string `json:"target_languages"`             // Язык цели в  iso639-1
	ContextAdaptation *bool    `json:"context_adaptation,omitempty"` // Позволяет или отключает контекстно-зависимые функции перевода, которые позволяют модели адаптировать переводы на основе предоставленного контекста.
	Context           *string  `json:"context,omitempty"`            // Информация о контексте для повышения точности перевода
	Informal          *bool    `json:"informal,omitempty"`           // Заставляет перевод использовать неформальные языковые формы, когда они доступны на целевом языке.
}

// Конфигурация для генерации субтитров, если  subtitlesвключено
type SubtitlesConf struct {
	Formats []string `json:"formats"` // Форматы субтитров, которые вы хотите, чтобы ваша транскрибация была отформатирована. Доступные опции: srt, vtt
}

// Инициировать транскрибирование. POST /v2/pre-recorded
type PreRecorderBody struct {
	AudioUrl          string           `json:"audio_url"`                    // Uploaded audio file Gladia URL. Example: "https://api.gladia.io/file/6c09400e-23d2-4bd2-be55-96a5ececfa3b"
	Diarization       bool             `json:"diarization"`                  // Включить автоматические определение спикеров (формат диалога)
	DiarizationConf   *DiarizationConf `json:"diarization_config,omitempty"` // Конфиг для более точного определения спикеров
	LangConf          *LanguageConf    `json:"language_config,omitempty"`    // Информация по исходному языку записи
	Translation       bool             `json:"translation"`                  // Нужно ли переводить
	TranslationConf   *TranslationConf `json:"translation_config,omitempty"` // Настройки перевода
	Subtitle          bool             `json:"subtitles"`                    // Нужны ли субтитры
	SubtitlesConf     *SubtitlesConf   `json:"subtitles_config,omitempty"`   // Настройки субтитров
	SentimentAnalysis bool             `json:"sentiment_analysis"`           // Включить анализ настроений для этого аудио
}

type PreRecorderInitResponse struct {
	ResultUrl string `json:"result_url"` // Предварительно встроенный URL с вашей транскрибацией
	ID        string `json:"id"`         // id чтобы получить результат
}
