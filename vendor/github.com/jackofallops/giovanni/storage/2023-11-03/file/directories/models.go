package directories

import "encoding/xml"

type ErrorResponse struct {
	XMLName xml.Name `xml:"Error"`
	Code    *string  `xml:"Code"`
	Message *string  `xml:"Message"`
}
