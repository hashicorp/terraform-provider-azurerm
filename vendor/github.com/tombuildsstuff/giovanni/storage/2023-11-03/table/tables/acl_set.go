package tables

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type setAcl struct {
	SignedIdentifiers []SignedIdentifier `xml:"SignedIdentifier"`

	XMLName xml.Name `xml:"SignedIdentifiers"`
}

type SetACLResponse struct {
	HttpResponse *client.Response
}

// SetACL sets the specified Access Control List for the specified Table
func (c Client) SetACL(ctx context.Context, tableName string, acls []SignedIdentifier) (resp SetACLResponse, err error) {
	if tableName == "" {
		return resp, fmt.Errorf("`tableName` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusNoContent,
		},
		HttpMethod:    http.MethodPut,
		OptionsObject: setAclTableOptions{},
		Path:          fmt.Sprintf("/%s", tableName),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		err = fmt.Errorf("building request: %+v", err)
		return
	}

	err = req.Marshal(setAcl{SignedIdentifiers: acls})
	if err != nil {
		return resp, fmt.Errorf("marshalling request: %+v", err)
	}

	resp.HttpResponse, err = req.Execute(ctx)
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	return
}

type setAclTableOptions struct{}

func (s setAclTableOptions) ToHeaders() *client.Headers {
	return nil
}

func (s setAclTableOptions) ToOData() *odata.Query {
	return nil
}

func (s setAclTableOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("comp", "acl")
	return out
}
