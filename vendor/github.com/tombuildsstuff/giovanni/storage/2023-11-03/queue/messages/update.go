package messages

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type UpdateInput struct {
	// A message must be in a format that can be included in an XML request with UTF-8 encoding.
	// The encoded message can be up to 64 KB in size.
	Message string

	// Specifies the valid pop receipt value required to modify this message.
	PopReceipt string

	// Specifies the new visibility timeout value, in seconds, relative to server time.
	// The new value must be larger than or equal to 0, and cannot be larger than 7 days.
	// The visibility timeout of a message cannot be set to a value later than the expiry time.
	// A message can be updated until it has been deleted or has expired.
	VisibilityTimeout int
}

type UpdateResponse struct {
	HttpResponse *client.Response
}

// Update updates an existing message based on it's Pop Receipt
func (c Client) Update(ctx context.Context, queueName string, messageID string, input UpdateInput) (resp UpdateResponse, err error) {

	if queueName == "" {
		return resp, fmt.Errorf("`queueName` cannot be an empty string")
	}
	if strings.ToLower(queueName) != queueName {
		return resp, fmt.Errorf("`queueName` must be a lower-cased string")
	}
	if input.PopReceipt == "" {
		return resp, fmt.Errorf("`input.PopReceipt` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusNoContent,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: updateOptions{
			input: input,
		},
		Path: fmt.Sprintf("/%s/messages/%s", queueName, messageID),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		err = fmt.Errorf("building request: %+v", err)
		return
	}

	marshalledMsg, err := xml.Marshal(QueueMessage{
		MessageText: input.Message,
	})
	if err != nil {
		return resp, fmt.Errorf("marshalling request: %v", err)
	}

	body := xml.Header + string(marshalledMsg)
	req.Body = io.NopCloser(bytes.NewReader([]byte(body)))
	req.ContentLength = int64(len(body))
	req.Header.Set("Content-Length", strconv.Itoa(len(body)))

	resp.HttpResponse, err = req.Execute(ctx)
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	return
}

type updateOptions struct {
	input UpdateInput
}

func (u updateOptions) ToHeaders() *client.Headers {
	return nil
}

func (u updateOptions) ToOData() *odata.Query {
	return nil
}

func (u updateOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("visibilitytimeout", strconv.Itoa(u.input.VisibilityTimeout))
	out.Append("popreceipt", u.input.PopReceipt)
	return out
}
