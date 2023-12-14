package queues

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

// GetMetaData returns the metadata for this Queue
func (c Client) GetMetaData(ctx context.Context, queueName string) (resp GetMetaDataResponse, err error) {

	if queueName == "" {
		return resp, fmt.Errorf("`queueName` cannot be an empty string")
	}

	if strings.ToLower(queueName) != queueName {
		return resp, fmt.Errorf("`queueName` must be a lower-cased string")
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: getMetaDataOptions{},
		Path:          fmt.Sprintf("/%s", queueName),
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

type getMetaDataOptions struct{}

func (g getMetaDataOptions) ToHeaders() *client.Headers {
	return nil
}

func (g getMetaDataOptions) ToOData() *odata.Query {
	return nil
}

func (g getMetaDataOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("comp", "metadata")
	return out
}
