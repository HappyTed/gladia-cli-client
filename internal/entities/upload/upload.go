package upload

type UploadResponce struct {
	AudioUrl string `json:"audio_url"`
	MetaData struct {
		ID            string  `json:"id"`
		FileName      string  `json:"filename"`
		Extension     string  `json:"extension"`
		Size          int     `json:"size"`
		Duration      float32 `json:"audio_duration"`
		NumOfChannels uint    `json:"number_of_channels"`
	} `json:"audio_metadata"`
}
