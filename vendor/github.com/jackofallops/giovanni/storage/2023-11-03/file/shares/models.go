package shares

import "encoding/xml"

type SignedIdentifier struct {
	Id           string       `xml:"Id"`
	AccessPolicy AccessPolicy `xml:"AccessPolicy"`
}

type AccessPolicy struct {
	Start      string `xml:"Start"`
	Expiry     string `xml:"Expiry"`
	Permission string `xml:"Permission"`
}

type ShareProtocol string

const (
	// SMB indicates the share can be accessed by SMBv3.0, SMBv2.1 and REST.
	SMB ShareProtocol = "SMB"

	// NFS indicates the share can be accessed by NFSv4.1. A premium account is required for this option.
	NFS ShareProtocol = "NFS"
)

type ErrorResponse struct {
	XMLName xml.Name `xml:"Error"`
	Code    *string  `xml:"Code"`
	Message *string  `xml:"Message"`
}
