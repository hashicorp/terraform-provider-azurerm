package afdendpoints

type AfdPurgeParameters struct {
	ContentPaths []string  `json:"contentPaths"`
	Domains      *[]string `json:"domains,omitempty"`
}
