package queues

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/jackofallops/giovanni/storage/internal/metadata"
)

type GetMetaDataResponse struct {
	HttpResponse *http.Response

	MetaData map[string]string
}

// GetMetaData returns the metadata for this Queue
func (c Client) GetMetaData(ctx context.Context, queueName string) (result GetMetaDataResponse, err error) {

	if queueName == "" {
		return result, fmt.Errorf("`queueName` cannot be an empty string")
	}

	if strings.ToLower(queueName) != queueName {
		return result, fmt.Errorf("`queueName` must be a lower-cased string")
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

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil && resp.Response != nil {
		result.HttpResponse = resp.Response

		if err == nil {
			if resp.Header != nil {
				result.MetaData = metadata.ParseFromHeaders(resp.Header)
			}
		}
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	if result.HttpResponse != nil {
		if result.HttpResponse.Header != nil {
			result.MetaData = metadata.ParseFromHeaders(result.HttpResponse.Header)
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
