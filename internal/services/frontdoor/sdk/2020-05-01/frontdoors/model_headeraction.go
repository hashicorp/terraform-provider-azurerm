package frontdoors

type HeaderAction struct {
	HeaderActionType HeaderActionType `json:"headerActionType"`
	HeaderName       string           `json:"headerName"`
	Value            *string          `json:"value,omitempty"`
}
