package shares

import (
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

var _ client.Options = sharesOptions{}

type sharesOptions struct{}

func (o sharesOptions) ToHeaders() *client.Headers {
	return &client.Headers{}
}

func (sharesOptions) ToOData() *odata.Query {
	return nil
}

func (sharesOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("restype", "share")
	return out
}
