package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"go-gladia.io-client/pkg/logger"
)

type Config struct {
	Token          string          `env:"API_KEY" env-required:"true"`
	BaseUrl        string          `env:"URL" env-default:"https://api.gladia.io"`
	AppEnv         string          `env:"ENV" env-default:"DEV"` // Среда, где запускается приложение
	LogLevel       logger.LogLevel `env:"LOG_LEVEL" env-default:"0"`
	LogPath        string          `env:"LOG_PATH" env-default:"logs/app.log"`
	RequestTimeout time.Duration   `env:"LOG_FORMATTER" env-default:"0"`
	Flags
	Transcription
	Database
}

type (
	Flags struct {
		AudioFile     string
		AwaitResults  bool
		AwaitInterval time.Duration
		AwaitTimeout  time.Duration
		OutputFile    string
	}

	// read from env
	Transcription struct {
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

	// read from env
	Database struct {
		DSN    string
		Driver string
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
