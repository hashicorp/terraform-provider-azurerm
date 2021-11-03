package attestationproviders

type JsonWebKeySet struct {
	Keys *[]JsonWebKey `json:"keys,omitempty"`
}
