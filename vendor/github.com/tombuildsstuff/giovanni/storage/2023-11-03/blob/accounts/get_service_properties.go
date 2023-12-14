package accounts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
)

type GetServicePropertiesResult struct {
	HttpResponse *client.Response
	Model        *StorageServiceProperties
}

func (c Client) GetServiceProperties(ctx context.Context, accountName string) (resp GetServicePropertiesResult, err error) {
	if accountName == "" {
		return resp, fmt.Errorf("`accountName` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ContentType: "text/xml",
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
	resp.HttpResponse, err = req.Execute(ctx)
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	if resp.HttpResponse != nil {
		if err = resp.HttpResponse.Unmarshal(&resp.Model); err != nil {
			err = fmt.Errorf("unmarshaling response: %+v", err)
			return
		}
	}

	return
}
