package modellibrary

type WatchdogStartRequest struct {
	FileName        string `json:"file_name"`
	IntervalSeconds int    `json:"interval_seconds"`
}

type WatchdogStartResponse struct{}
