package accounts

import (
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

var _ client.Options = servicePropertiesOptions{}

type servicePropertiesOptions struct {
}

func (servicePropertiesOptions) ToHeaders() *client.Headers {
	return nil
}

func (servicePropertiesOptions) ToOData() *odata.Query {
	return nil
}

func (servicePropertiesOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("comp", "properties")
	out.Append("restype", "service")
	return out
}
