package messages

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type GetInput struct {
	// VisibilityTimeout specifies the new visibility timeout value, in seconds, relative to server time.
	// The new value must be larger than or equal to 0, and cannot be larger than 7 days.
	VisibilityTimeout *int

	// NumberOfMessages specifies the (maximum) number of messages that should be retrieved from the queue.
	// This can be a maximum of 32.
	NumberOfMessages int
}

type GetResponse struct {
	HttpResponse *client.Response
}

// Get retrieves one or more messages from the front of the queue
func (c Client) Get(ctx context.Context, queueName string, input GetInput) (resp QueueMessagesListResponse, err error) {
	if queueName == "" {
		return resp, fmt.Errorf("`queueName` cannot be an empty string")
	}
	if strings.ToLower(queueName) != queueName {
		return resp, fmt.Errorf("`queueName` must be a lower-cased string")
	}
	if input.NumberOfMessages < 1 || input.NumberOfMessages > 32 {
		return resp, fmt.Errorf("`input.NumberOfMessages` must be between 1 and 32")
	}
	if input.VisibilityTimeout != nil {
		t := *input.VisibilityTimeout
		maxTime := (time.Hour * 24 * 7).Seconds()
		if t < 1 || t < int(maxTime) {
			return resp, fmt.Errorf("`input.VisibilityTimeout` must be larger than or equal to 1 second, and cannot be larger than 7 days")
		}
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		OptionsObject: getOptions{
			visibilityTimeout: input.VisibilityTimeout,
			numberOfMessages:  input.NumberOfMessages,
		},
		Path: fmt.Sprintf("/%s/messages", queueName),
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
		if err = resp.HttpResponse.Unmarshal(&resp); err != nil {
			return resp, fmt.Errorf("unmarshalling response: %+v", err)
		}
	}

	return
}

type getOptions struct {
	visibilityTimeout *int
	numberOfMessages  int
}

func (g getOptions) ToHeaders() *client.Headers {
	return nil
}

func (g getOptions) ToOData() *odata.Query {
	return nil
}

func (g getOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	if g.visibilityTimeout != nil {
		out.Append("visibilitytimeout", strconv.Itoa(*g.visibilityTimeout))
	}
	out.Append("numofmessages", strconv.Itoa(g.numberOfMessages))
	return out
}
