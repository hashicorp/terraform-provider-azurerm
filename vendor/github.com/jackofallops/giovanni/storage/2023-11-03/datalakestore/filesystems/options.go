package filesystems

import (
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

var _ client.Options = fileSystemOptions{}

type fileSystemOptions struct {
	properties map[string]string
}

func (o fileSystemOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}
	props := buildProperties(o.properties)
	if props != "" {
		headers.Append("x-ms-properties", props)
	}

	return headers
}

func (fileSystemOptions) ToOData() *odata.Query {
	return nil
}

func (fileSystemOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("resource", "filesystem")
	return out
}
