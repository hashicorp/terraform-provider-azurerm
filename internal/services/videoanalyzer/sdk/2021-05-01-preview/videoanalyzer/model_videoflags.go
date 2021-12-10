package videoanalyzer

type VideoFlags struct {
	CanStream   bool `json:"canStream"`
	HasData     bool `json:"hasData"`
	IsRecording bool `json:"isRecording"`
}
