package directories

import (
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

var _ client.Options = directoriesOptions{}

type directoriesOptions struct{}

func (o directoriesOptions) ToHeaders() *client.Headers {
	return &client.Headers{}
}

func (directoriesOptions) ToOData() *odata.Query {
	return nil
}

func (directoriesOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("restype", "directory")
	return out
}
