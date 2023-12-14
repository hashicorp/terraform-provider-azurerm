package shares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type GetACLResult struct {
	HttpResponse *client.Response

	SignedIdentifiers []SignedIdentifier `xml:"SignedIdentifier"`
}

// GetACL get the Access Control List for the specified Storage Share
func (c Client) GetACL(ctx context.Context, shareName string) (resp GetACLResult, err error) {

	if shareName == "" {
		return resp, fmt.Errorf("`shareName` cannot be an empty string")
	}
	if strings.ToLower(shareName) != shareName {
		return resp, fmt.Errorf("`shareName` must be a lower-cased string")
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: getAclOptions{},
		Path:          fmt.Sprintf("/%s", shareName),
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
		err = resp.HttpResponse.Unmarshal(&resp)
		if err != nil {
			return resp, fmt.Errorf("unmarshalling response: %v", err)
		}
	}
	return
}

type getAclOptions struct {
}

func (g getAclOptions) ToHeaders() *client.Headers {
	return nil
}

func (g getAclOptions) ToOData() *odata.Query {
	return nil
}

func (g getAclOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("restype", "share")
	out.Append("comp", "acl")
	return out
}
