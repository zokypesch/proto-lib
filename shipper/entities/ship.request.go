package entities

type ShipLogRequestPropertyResponse struct {
	StatusCode int                    `json:"status_code"`
	Content    string                 `json:"content"`
	Others     map[string]interface{} `json:"others"`
}

type ShipLogRequestPropertyRequest struct {
	Headers string                 `json:"headers"`
	Payload string                 `json:"payload"`
	Others  map[string]interface{} `json:"others"`
}

type ShipLogRequestProperty struct {
	Request  ShipLogRequestPropertyRequest `json:"request"`
	Response ShipLogRequestPropertyResponse `json:"response"`
}

type ShipLogRequest struct {
	Action    string `json:"action"`   // full rul
	Todo      string `json:"todo"`     // Http Method
	Property  string `json:"property"` // Contains request & response
	Timestamp int64  `json:"timestamp"`
}
