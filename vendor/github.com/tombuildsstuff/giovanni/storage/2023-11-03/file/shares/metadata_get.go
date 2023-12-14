package shares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/tombuildsstuff/giovanni/storage/internal/metadata"
)

type GetMetaDataResponse struct {
	HttpResponse *client.Response

	MetaData map[string]string
}

// GetMetaData returns the MetaData associated with the specified Storage Share
func (c Client) GetMetaData(ctx context.Context, shareName string) (resp GetMetaDataResponse, err error) {
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
		OptionsObject: GetMetaDataOptions{},
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
		if resp.HttpResponse.Header != nil {
			resp.MetaData = metadata.ParseFromHeaders(resp.HttpResponse.Header)
		}
	}

	return
}

type GetMetaDataOptions struct{}

func (g GetMetaDataOptions) ToHeaders() *client.Headers {
	return nil
}

func (g GetMetaDataOptions) ToOData() *odata.Query {
	return nil
}

func (g GetMetaDataOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("restype", "share")
	out.Append("comp", "metadata")
	return out
}
