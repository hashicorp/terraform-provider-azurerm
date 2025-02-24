package accounts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
)

type SetServicePropertiesResult struct {
	HttpResponse *http.Response
}

func (c Client) SetServiceProperties(ctx context.Context, accountName string, input StorageServiceProperties) (result SetServicePropertiesResult, err error) {
	if accountName == "" {
		return result, fmt.Errorf("`accountName` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
		},
		HttpMethod:    http.MethodPut,
		OptionsObject: servicePropertiesOptions{},
		Path:          "/",
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		err = fmt.Errorf("building request: %+v", err)
		return
	}

	if err = req.Marshal(&input); err != nil {
		err = fmt.Errorf("marshaling request: %+v", err)
		return
	}

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
