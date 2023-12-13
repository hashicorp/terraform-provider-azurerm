package accounts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
)

type SetServicePropertiesResult struct {
	HttpResponse *client.Response
}

func (c Client) SetServiceProperties(ctx context.Context, accountName string, input StorageServiceProperties) (resp SetServicePropertiesResult, err error) {
	if accountName == "" {
		return resp, fmt.Errorf("`accountName` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ContentType: "text/xml",
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
	resp.HttpResponse, err = req.Execute(ctx)
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	return
}
