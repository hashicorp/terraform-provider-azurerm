package shares

type SignedIdentifier struct {
	Id           string       `xml:"Id"`
	AccessPolicy AccessPolicy `xml:"AccessPolicy"`
}

type AccessPolicy struct {
	Start      string `xml:"Start"`
	Expiry     string `xml:"Expiry"`
	Permission string `xml:"Permission"`
}
