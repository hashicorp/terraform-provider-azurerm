package endpoints

type HttpErrorRangeParameters struct {
	Begin *int64 `json:"begin,omitempty"`
	End   *int64 `json:"end,omitempty"`
}
