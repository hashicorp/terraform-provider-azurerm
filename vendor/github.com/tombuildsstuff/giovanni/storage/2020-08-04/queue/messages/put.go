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

type PutInput struct {
	// A message must be in a format that can be included in an XML request with UTF-8 encoding.
	// The encoded message can be up to 64 KB in size.
	Message string

	// The maximum time-to-live can be any positive number,
	// as well as -1 indicating that the message does not expire.
	// If this parameter is omitted, the default time-to-live is 7 days.
	MessageTtl *int

	// Specifies the new visibility timeout value, in seconds, relative to server time.
	// The new value must be larger than or equal to 0, and cannot be larger than 7 days.
	// The visibility timeout of a message cannot be set to a value later than the expiry time.
	// visibilitytimeout should be set to a value smaller than the time-to-live value.
	// If not specified, the default value is 0.
	VisibilityTimeout *int
}

// Put adds a new message to the back of the message queue
func (c Client) Put(ctx context.Context, queueName string, input PutInput) (resp QueueMessagesListResponse, err error) {
	if queueName == "" {
		return resp, fmt.Errorf("`queueName` cannot be an empty string")
	}

	if strings.ToLower(queueName) != queueName {
		return resp, fmt.Errorf("`queueName` must be a lower-cased string")
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusCreated,
		},
		HttpMethod: http.MethodPost,
		OptionsObject: putOptions{
			input: input,
		},
		Path: fmt.Sprintf("/%s/messages", queueName),
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

	if resp.HttpResponse != nil {
		if err = resp.HttpResponse.Unmarshal(&resp); err != nil {
			return resp, fmt.Errorf("unmarshalling response: %+v", err)
		}
	}

	return
}

type putOptions struct {
	input PutInput
}

func (p putOptions) ToHeaders() *client.Headers {
	return nil
}

func (p putOptions) ToOData() *odata.Query {
	return nil
}

func (p putOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}

	if p.input.MessageTtl != nil {
		out.Append("messagettl", strconv.Itoa(*p.input.MessageTtl))
	}

	if p.input.VisibilityTimeout != nil {
		out.Append("visibilitytimeout", strconv.Itoa(*p.input.VisibilityTimeout))
	}

	return out
}
