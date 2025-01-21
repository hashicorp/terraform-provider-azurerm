package queues

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type GetStorageServicePropertiesResponse struct {
	StorageServiceProperties
	HttpResponse *http.Response
}

// GetServiceProperties gets the properties for this queue
func (c Client) GetServiceProperties(ctx context.Context) (result GetStorageServicePropertiesResponse, err error) {

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: getStorageServicePropertiesOptions{},
		Path:          "/",
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

type getStorageServicePropertiesOptions struct{}

func (g getStorageServicePropertiesOptions) ToHeaders() *client.Headers {
	return nil
}

func (g getStorageServicePropertiesOptions) ToOData() *odata.Query {
	return nil
}

func (g getStorageServicePropertiesOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("comp", "properties")
	out.Append("restype", "service")
	return out
}
