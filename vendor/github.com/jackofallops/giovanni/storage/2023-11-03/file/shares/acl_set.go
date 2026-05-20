package shares

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type SetAclResponse struct {
	HttpResponse *http.Response
}

type SetAclInput struct {
	SignedIdentifiers []SignedIdentifier `xml:"SignedIdentifier"`
	XMLName           xml.Name           `xml:"SignedIdentifiers"`
}

// SetACL sets the specified Access Control List on the specified Storage Share
func (c Client) SetACL(ctx context.Context, shareName string, input SetAclInput) (result SetAclResponse, err error) {
	if shareName == "" {
		err = fmt.Errorf("`shareName` cannot be an empty string")
		return
	}

	if strings.ToLower(shareName) != shareName {
		err = fmt.Errorf("`shareName` must be a lower-cased string")
		return
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPut,
		OptionsObject: setAclOptions{},
		Path:          fmt.Sprintf("/%s", shareName),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		err = fmt.Errorf("building request: %+v", err)
		return
	}

	b, err := xml.Marshal(&input)
	if err != nil {
		err = fmt.Errorf("marshalling input: %+v", err)
		return
	}
	withHeader := xml.Header + string(b)
	bytesWithHeader := []byte(withHeader)
	req.ContentLength = int64(len(bytesWithHeader))
	req.Header.Set("Content-Length", strconv.Itoa(len(bytesWithHeader)))
	req.Body = io.NopCloser(bytes.NewReader(bytesWithHeader))

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil && resp.Response != nil {
		result.HttpResponse = resp.Response
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	return
}

type setAclOptions struct {
	SignedIdentifiers []SignedIdentifier `xml:"SignedIdentifier"`
}

func (s setAclOptions) ToHeaders() *client.Headers {
	return nil
}

func (s setAclOptions) ToOData() *odata.Query {
	return nil
}

func (s setAclOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("restype", "share")
	out.Append("comp", "acl")
	return out
}
