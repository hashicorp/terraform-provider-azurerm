package messages

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type DeleteResponse struct {
	HttpResponse *client.Response
}

type DeleteInput struct {
	PopReceipt string
}

// Delete deletes a specific message
func (c Client) Delete(ctx context.Context, queueName, messageID string, input DeleteInput) (resp DeleteResponse, err error) {

	if queueName == "" {
		return resp, fmt.Errorf("`queueName` cannot be an empty string")
	}

	if strings.ToLower(queueName) != queueName {
		return resp, fmt.Errorf("`queueName` must be a lower-cased string")
	}

	if messageID == "" {
		return resp, fmt.Errorf("`messageID` cannot be an empty string")
	}

	if input.PopReceipt == "" {
		return resp, fmt.Errorf("`input.PopReceipt` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusNoContent,
		},
		HttpMethod: http.MethodDelete,
		OptionsObject: deleteOptions{
			popReceipt: input.PopReceipt,
		},
		Path: fmt.Sprintf("/%s/messages/%s", queueName, messageID),
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

type deleteOptions struct {
	popReceipt string
}

func (d deleteOptions) ToHeaders() *client.Headers {
	return nil
}

func (d deleteOptions) ToOData() *odata.Query {
	return nil
}

func (d deleteOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("popreceipt", d.popReceipt)
	return out
}
