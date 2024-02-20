package directories

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

// GetMetaData returns all user-defined metadata for the specified directory
func (c Client) GetMetaData(ctx context.Context, shareName, path string) (resp GetMetaDataResponse, err error) {
	if shareName == "" {
		return resp, fmt.Errorf("`shareName` cannot be an empty string")
	}

	if strings.ToLower(shareName) != shareName {
		return resp, fmt.Errorf("`shareName` must be a lower-cased string")
	}

	if path == "" {
		return resp, fmt.Errorf("`path` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: GetMetaDataOptions{},
		Path:          fmt.Sprintf("/%s/%s", shareName, path),
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
	out.Append("restype", "directory")
	out.Append("comp", "metadata")
	return out
}
