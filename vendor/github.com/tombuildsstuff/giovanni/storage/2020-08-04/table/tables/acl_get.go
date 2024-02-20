package tables

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type GetACLResponse struct {
	HttpResponse *client.Response

	SignedIdentifiers []SignedIdentifier `xml:"SignedIdentifier"`
}

// GetACL returns the Access Control List for the specified Table
func (c Client) GetACL(ctx context.Context, tableName string) (resp GetACLResponse, err error) {

	if tableName == "" {
		return resp, fmt.Errorf("`tableName` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: getAclTableOptions{},
		Path:          fmt.Sprintf("/%s", tableName),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		err = fmt.Errorf("building request: %+v", err)
		return
	}

	resp.HttpResponse, err = req.Execute(ctx)
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	if resp.HttpResponse != nil {
		if resp.HttpResponse.Body != nil {
			err = resp.HttpResponse.Unmarshal(&resp)
			if err != nil {
				return resp, fmt.Errorf("unmarshalling response body: %v", err)
			}
		}
	}

	return
}

type getAclTableOptions struct{}

func (g getAclTableOptions) ToHeaders() *client.Headers {
	return nil
}

func (g getAclTableOptions) ToOData() *odata.Query {
	return nil
}

func (g getAclTableOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("comp", "acl")
	return out
}
