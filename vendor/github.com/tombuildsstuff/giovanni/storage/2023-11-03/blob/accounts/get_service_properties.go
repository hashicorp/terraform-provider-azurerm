package accounts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
)

type GetServicePropertiesResult struct {
	StorageServiceProperties
	HttpResponse *http.Response
}

func (c Client) GetServiceProperties(ctx context.Context, accountName string) (result GetServicePropertiesResult, err error) {
	if accountName == "" {
		return result, fmt.Errorf("`accountName` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: servicePropertiesOptions{},
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
