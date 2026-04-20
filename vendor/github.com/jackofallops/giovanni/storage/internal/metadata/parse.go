package metadata

import (
	"net/http"
	"strings"
)

// ParseFromHeaders parses the metadata from the headers
func ParseFromHeaders(headers http.Header) map[string]string {
	metaData := make(map[string]string, 0)
	prefix := "x-ms-meta-"
	for k, v := range headers {
		key := strings.ToLower(k)
		if !strings.HasPrefix(key, prefix) {
			continue
		}

		// TODO: update this to support case-insensitive headers when the base layer is changed to `hashicorp/go-azure-sdk`
		// (e.g. trim off the first 10 characters, but this can't be done until the base layer is updated, since
		// `Azure/go-autorest` canonicalizes the header keys)
		key = strings.TrimPrefix(key, prefix)
		metaData[key] = v[0]
	}
	return metaData
}
