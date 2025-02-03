package containers

import (
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/jackofallops/giovanni/storage/internal/metadata"
)

var _ client.Options = containerOptions{}

type containerOptions struct {
	metaData map[string]string
}

func (o containerOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}

	metaDataHeaders := make(map[string]interface{})
	metadata.SetIntoHeaders(metaDataHeaders, o.metaData)
	for k, v := range metaDataHeaders {
		headers.Append(k, v.(string))
	}

	return headers
}

func (containerOptions) ToOData() *odata.Query {
	return nil
}

func (containerOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("restype", "container")
	return out
}
