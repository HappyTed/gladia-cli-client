package prerecorderv2

type ListResponse struct {
	First   string     `json:"first"`
	Current string     `json:"current"`
	Next    string     `json:"next"`
	Items   []ListItem `json:"items"`
}

type ListItem struct {
	ID                  string         `json:"id"`
	RequestID           string         `json:"request_id"`
	Version             int            `json:"version"` // api version
	Status              string         `json:"status"`
	CreatedAT           string         `json:"created_at"`
	PostSessionMetadata any            `json:"post_session_metadata"` // Object
	Kind                string         `json:"kind"`
	CompletedAT         string         `json:"completed_at"`
	CustomMetadata      map[string]any `json:"custom_metadata"`
	ErrorCode           int            `json:"error_code"`
	File                FileInfo       `json:"file"`
	ReqParams           ReqParams      `json:"request_params"`
	Result              Result         `json:"result"`
}
