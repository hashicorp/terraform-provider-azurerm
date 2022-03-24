package connections

type ApiConnectionTestLink struct {
	Method     *string `json:"method,omitempty"`
	RequestUri *string `json:"requestUri,omitempty"`
}
