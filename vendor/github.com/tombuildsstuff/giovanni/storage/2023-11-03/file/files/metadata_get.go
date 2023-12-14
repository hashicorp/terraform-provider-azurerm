package files

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

// GetMetaData returns the MetaData for the specified File.
func (c Client) GetMetaData(ctx context.Context, shareName, path, fileName string) (resp GetMetaDataResponse, err error) {
	if shareName == "" {
		return resp, fmt.Errorf("`shareName` cannot be an empty string")
	}

	if strings.ToLower(shareName) != shareName {
		return resp, fmt.Errorf("`shareName` must be a lower-cased string")
	}

	if fileName == "" {
		return resp, fmt.Errorf("`fileName` cannot be an empty string")
	}

	if path != "" {
		path = fmt.Sprintf("%s/", path)
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: GetMetadataOptions{},
		Path:          fmt.Sprintf("/%s/%s%s", shareName, path, fileName),
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

type GetMetadataOptions struct{}

func (m GetMetadataOptions) ToHeaders() *client.Headers {
	return nil
}

func (m GetMetadataOptions) ToOData() *odata.Query {
	return nil
}

func (m GetMetadataOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("comp", "metadata")
	return out
}
