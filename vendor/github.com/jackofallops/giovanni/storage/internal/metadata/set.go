package metadata

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
)

// SetIntoHeaders sets the provided MetaData into the headers
func SetIntoHeaders(headers map[string]interface{}, metaData map[string]string) map[string]interface{} {
	for k, v := range metaData {
		key := fmt.Sprintf("x-ms-meta-%s", k)
		headers[key] = v
	}

	return headers
}

// SetMetaDataHeaders sets the provided MetaData into the headers
func SetMetaDataHeaders(metaData map[string]string) client.Headers {
	headers := client.Headers{}
	for k, v := range metaData {
		key := fmt.Sprintf("x-ms-meta-%s", k)
		headers.Append(key, v)
	}
	return headers
}
