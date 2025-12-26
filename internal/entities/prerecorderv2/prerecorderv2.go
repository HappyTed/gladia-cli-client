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

type PreRecorderResultResponse struct {
	ID            string      `json:"id"`
	RequestID     string      `json:"request_id"`
	Version       int         `json:"version"` // версия api
	Status        string      `json:"status"`  // "queued": поставлена в очередь. "processing": обрабатывается. "done": результат доступен. "error": во время обработки работы произошла ошибка.
	CreatedAt     string      `json:"created_at"`
	CompletedAt   *string     `json:"completed_at,omitempty"`
	CustomMeta    interface{} `json:"custom_metadata,omitempty"` // Пользовательские метаданные, приведенные в первоначальном запросе
	ErrorCode     interface{} `json:"error_code,omitempty"`      // HTTP код состояния ошибки, если статус "error"
	Kind          *string     `json:"kind,omitempty"`
	File          *FileInfo   `json:"file,omitempty"`           // Данные по загруженному для обработки файлу
	RequestParams *ReqParams  `json:"request_params,omitempty"` // Параметры, используемые при транскрибации. Может быть нулевым, если статус является «ошибкой»
	Result        *Result     `json:"result,omitempty"`         // Предварительно записанный результат транскрибации, когда статус "сделан"
}

type FileInfo struct {
	ID               string      `json:"id"`
	Filename         string      `json:"filename"`
	Source           interface{} `json:"source"`
	AudioDuration    float64     `json:"audio_duration"`
	NumberOfChannels int         `json:"number_of_channels"`
}

type ReqParams struct {
	AudioURL                 string       `json:"audio_url"`
	Model                    string       `json:"model"`
	Sentences                bool         `json:"sentences"`
	Subtitles                bool         `json:"subtitles"`
	Moderation               bool         `json:"moderation"`
	Diarization              bool         `json:"diarization"`
	Translation              bool         `json:"translation"`
	AudioToLLM               bool         `json:"audio_to_llm"`
	DisplayMode              bool         `json:"display_mode"`
	Summarization            bool         `json:"summarization"`
	AudioEnhancer            bool         `json:"audio_enhancer"`
	Chapterization           bool         `json:"chapterization"`
	CustomSpelling           bool         `json:"custom_spelling"`
	DetectLanguage           bool         `json:"detect_language"`
	LanguageConfig           LanguageConf `json:"language_config"`
	NameConsistency          bool         `json:"name_consistency"`
	SentimentAnalysis        bool         `json:"sentiment_analysis"`
	DiarizationEnhanced      bool         `json:"diarization_enhanced"`
	PunctuationEnhanced      bool         `json:"punctuation_enhanced"`
	EnableCodeSwitching      bool         `json:"enable_code_switching"`
	NamedEntityRecognition   bool         `json:"named_entity_recognition"`
	SpeakerReidentification  bool         `json:"speaker_reidentification"`
	AccurateWordsTimestamps  bool         `json:"accurate_words_timestamps"`
	SkipChannelDeduplication bool         `json:"skip_channel_deduplication"`
	StructuredDataExtraction bool         `json:"structured_data_extraction"`
}

type Result struct {
	Metadata                 Metadata          `json:"metadata"`
	Transcription            Transcription     `json:"transcription"`
	Translation              TaskResult        `json:"translation"`
	Summarization            SimpleTaskResult  `json:"summarization"`
	Moderation               SimpleTaskResult  `json:"moderation"`
	NamedEntityRecognition   NamedEntityResult `json:"named_entity_recognition"`
	NameConsistency          SimpleTaskResult  `json:"name_consistency"`
	SpeakerReidentification  SimpleTaskResult  `json:"speaker_reidentification"`
	StructuredDataExtraction SimpleTaskResult  `json:"structured_data_extraction"`
	SentimentAnalysis        SimpleTaskResult  `json:"sentiment_analysis"`
	AudioToLLM               AudioToLLMResult  `json:"audio_to_llm"`
	Sentences                SimpleTaskResult  `json:"sentences"`
	DisplayMode              SimpleTaskResult  `json:"display_mode"`
	Chapterization           SimpleTaskResult  `json:"chapterization"`
	Diarization              DiarizationResult `json:"diarization"`
}

// составные части структуры Result

type Metadata struct {
	AudioDuration            float64 `json:"audio_duration"`
	NumberOfDistinctChannels int     `json:"number_of_distinct_channels"`
	BillingTime              float64 `json:"billing_time"`
	TranscriptionTime        float64 `json:"transcription_time"`
}

type Transcription struct {
	Languages      []string    `json:"languages"`
	Utterances     []Utterance `json:"utterances"`
	FullTranscript string      `json:"full_transcript"`
}

type Utterance struct {
	Text       string  `json:"text"`
	Language   string  `json:"language"`
	Start      float64 `json:"start"`
	End        float64 `json:"end"`
	Confidence float64 `json:"confidence"`
	Channel    int     `json:"channel"`
	Words      []Word  `json:"words"`
	Speaker    *int    `json:"speaker,omitempty"`
}

type Word struct {
	Word       string  `json:"word"`
	Start      float64 `json:"start"`
	End        float64 `json:"end"`
	Confidence float64 `json:"confidence"`
}

type Subtitle struct {
	Format    string `json:"format"`
	Subtitles string `json:"subtitles"`
}

type ErrorInfo struct {
	StatusCode int    `json:"status_code"`
	Exception  string `json:"exception"`
	Message    string `json:"message"`
}

type SentenceResult struct {
	Success  bool      `json:"success"`
	IsEmpty  bool      `json:"is_empty"`
	ExecTime int       `json:"exec_time"`
	Error    ErrorInfo `json:"error"`
	Results  []string  `json:"results"`
}

type SimpleTaskResult struct {
	Success  bool      `json:"success"`
	IsEmpty  bool      `json:"is_empty"`
	ExecTime int       `json:"exec_time"`
	Error    ErrorInfo `json:"error"`
	Results  any       `json:"results"`
}

type TaskResult struct {
	Success  bool            `json:"success"`
	IsEmpty  bool            `json:"is_empty"`
	ExecTime int             `json:"exec_time"`
	Error    ErrorInfo       `json:"error"`
	Results  []Transcription `json:"results"`
}

type NamedEntityResult struct {
	Success  bool      `json:"success"`
	IsEmpty  bool      `json:"is_empty"`
	ExecTime int       `json:"exec_time"`
	Error    ErrorInfo `json:"error"`
	Entity   string    `json:"entity"`
}

type AudioToLLMResult struct {
	Success  bool             `json:"success"`
	IsEmpty  bool             `json:"is_empty"`
	ExecTime int              `json:"exec_time"`
	Error    ErrorInfo        `json:"error"`
	Results  []AudioToLLMItem `json:"results"`
}

type AudioToLLMItem struct {
	Success  bool      `json:"success"`
	IsEmpty  bool      `json:"is_empty"`
	ExecTime int       `json:"exec_time"`
	Error    ErrorInfo `json:"error"`
	Results  LLMResult `json:"results"`
}

type LLMResult struct {
	Prompt   string `json:"prompt"`
	Response string `json:"response"`
}

type DiarizationResult struct {
	Success  bool        `json:"success"`
	IsEmpty  bool        `json:"is_empty"`
	ExecTime int         `json:"exec_time"`
	Error    ErrorInfo   `json:"error"`
	Results  []Utterance `json:"results"`
}
