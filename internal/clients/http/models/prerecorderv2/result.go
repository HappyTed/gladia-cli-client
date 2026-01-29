package prerecorderv2

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
