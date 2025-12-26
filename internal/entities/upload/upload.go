package upload

// POST /v2/upload
type UploadResponce struct {
	AudioUrl string `json:"audio_url"` // Uploaded audio file Gladia URL. Example: "https://api.gladia.io/file/6c09400e-23d2-4bd2-be55-96a5ececfa3b"
	MetaData struct {
		ID            string  `json:"id"`
		FileName      string  `json:"filename"`
		Extension     string  `json:"extension"`
		Size          int     `json:"size"`
		Duration      float32 `json:"audio_duration"`
		NumOfChannels uint    `json:"number_of_channels"`
	} `json:"audio_metadata"`
}

func (r *UploadResponce) Dump() []byte {
	return nil
}
