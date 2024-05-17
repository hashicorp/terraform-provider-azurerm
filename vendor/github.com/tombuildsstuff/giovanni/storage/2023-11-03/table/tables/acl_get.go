package tables

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type GetACLResponse struct {
	HttpResponse *http.Response

	SignedIdentifiers []SignedIdentifier `xml:"SignedIdentifier"`
}

// GetACL returns the Access Control List for the specified Table
func (c Client) GetACL(ctx context.Context, tableName string) (result GetACLResponse, err error) {

	if tableName == "" {
		err = fmt.Errorf("`tableName` cannot be an empty string")
		return
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

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil && resp.Response != nil {
		result.HttpResponse = resp.Response

		if err == nil {
			err = resp.Unmarshal(&result)
			if err != nil {
				err = fmt.Errorf("unmarshalling response: %+v", err)
				return
			}
		}
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
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
