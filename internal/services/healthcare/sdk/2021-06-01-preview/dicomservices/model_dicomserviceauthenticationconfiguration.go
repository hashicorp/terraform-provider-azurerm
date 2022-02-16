package dicomservices

type DicomServiceAuthenticationConfiguration struct {
	Audiences *[]string `json:"audiences,omitempty"`
	Authority *string   `json:"authority,omitempty"`
}
