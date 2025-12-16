package prerecorderv2

type LanguageConf struct {
	Languages     []string `json:"languages"`
	CodeSwitching bool     `json:"code_switching"` //
}

type DiarizationConf struct {
	Enhanced      bool   `json:"enhanced"`
	NumOfSpeakers *uint8 `json:"number_of_speakers,omitempty"`
	MinSpeakers   *uint8 `json:"min_speakers,omitempty"`
	MaxSpeakers   *uint8 `json:"max_speakers,omitempty"`
}

type TranslationConf struct {
	Model             *string  `json:"model,omitempty"`
	TargetLanguages   []string `json:"target_languages"`
	ContextAdaptation *bool    `json:"context_adaptation,omitempty"`
	Context           *string  `json:"context,omitempty"`
	Informal          *bool    `json:"informal,omitempty"`
}

type SubtitlesConf struct {
	Formats []string `json:"formats"`
}

type PreRecorderBody struct {
	AudioUrl          string           `json:"audio_url"`
	Diarization       bool             `json:"diarization"`
	LangConf          *LanguageConf    `json:"language_config,omitempty"`
	DiarizationConf   *DiarizationConf `json:"diarization_config,omitempty"`
	Translation       bool             `json:"translation"`
	TranslationConf   *TranslationConf `json:"translation_config,omitempty"`
	Subtitle          bool             `json:"subtitles"`
	SubtitlesConf     *SubtitlesConf   `json:"subtitles_config,omitempty"`
	SentimentAnalysis bool             `json:"sentiment_analysis"`
}

type PreRecorderInitResponse struct {
	ResultUrl string `json:"result_url"`
	ID        string `json:"id"`
}

type PreRecorderResultResponse struct {
	ID            string      `json:"id"`
	RequestID     string      `json:"request_id"`
	Version       int         `json:"version"`
	Status        string      `json:"status"`
	CreatedAt     string      `json:"created_at"`
	CompletedAt   *string     `json:"completed_at,omitempty"`
	CustomMeta    interface{} `json:"custom_metadata,omitempty"`
	ErrorCode     interface{} `json:"error_code,omitempty"`
	Kind          *string     `json:"kind,omitempty"`
	File          *FileInfo   `json:"file,omitempty"`
	RequestParams *ReqParams  `json:"request_params,omitempty"`
	Result        *Result     `json:"result,omitempty"`
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
	Metadata      Metadata      `json:"metadata"`
	Transcription Transcription `json:"transcription"`
}

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
