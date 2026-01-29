package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Token   string `env:"API_KEY" env-required:"true"`
	BaseUrl string `env:"BASE_URL" env-default:"https://api.gladia.io"`
	IsDebug bool
	Flags
	TranscriptionConfig
	HTTPClientConfig
	WSClientConfig
}

type (
	Flags struct {
		AudioFile     string
		AwaitResults  bool
		AwaitInterval time.Duration
		AwaitTimeout  time.Duration
		OutputFile    string
	}

	HTTPClientConfig struct {
		Timeout    time.Duration
		MaxRetries uint8
	}

	WSClientConfig struct{}

	// read from env
	TranscriptionConfig struct {
		Diarization       bool
		Enhanced          bool
		Speakers          *uint8
		MaxSpeakers       *uint8
		MinSpeakers       *uint8
		Translation       bool
		TargetLanguages   []string
		SentimentAnalysis bool
		InputLanguages    []string
	}
)

func LoadConfig() *Config {
	cfg := &Config{
		Token:          "",
		BaseUrl:        "https://api.gladia.io",
		AppEnv:         "DEV",
		LogLevel:       logger.DEBUG,
		LogPath:        "/tmp/app.log",
		RequestTimeout: 0,
		Transcription: Transcription{
			Diarization:       false,
			Enhanced:          true,
			Translation:       false,
			SentimentAnalysis: true,
		},
		Flags: Flags{
			AwaitInterval: time.Second * 5,
			AwaitTimeout:  0,
			OutputFile:    "result.txt",
		},
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		log.Fatal("fail to read env: ", err)
	}

	return cfg
}
