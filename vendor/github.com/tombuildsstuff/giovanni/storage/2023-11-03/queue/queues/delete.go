package queues

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
)

type DeleteResponse struct {
	HttpResponse *client.Response
}

// Delete deletes the specified Queue within the specified Storage Account
func (c Client) Delete(ctx context.Context, queueName string) (resp DeleteResponse, err error) {

	if queueName == "" {
		return resp, fmt.Errorf("`queueName` cannot be an empty string")
	}

	if strings.ToLower(queueName) != queueName {
		return resp, fmt.Errorf("`queueName` must be a lower-cased string")
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusNoContent,
		},
		HttpMethod:    http.MethodDelete,
		OptionsObject: nil,
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

	return
}
