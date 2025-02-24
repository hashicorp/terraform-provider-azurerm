package queues

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/jackofallops/giovanni/storage/internal/metadata"
)

type CreateInput struct {
	MetaData map[string]string
}

type CreateResponse struct {
	HttpResponse *http.Response
}

// Create creates the specified Queue within the specified Storage Account
func (c Client) Create(ctx context.Context, queueName string, input CreateInput) (result CreateResponse, err error) {

	if queueName == "" {
		return result, fmt.Errorf("`queueName` cannot be an empty string")
	}

	if strings.ToLower(queueName) != queueName {
		return result, fmt.Errorf("`queueName` must be a lower-cased string")
	}

	if err := metadata.Validate(input.MetaData); err != nil {
		return result, fmt.Errorf("`metadata` is not valid: %s", err)
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusCreated,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: createOptions{
			metadata: input.MetaData,
		},
		Path: fmt.Sprintf("/%s", queueName),
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
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	return
}

type createOptions struct {
	metadata map[string]string
}

func (c createOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}

	if len(c.metadata) > 0 {
		headers.Merge(metadata.SetMetaDataHeaders(c.metadata))
	}
	return headers
}

func (c createOptions) ToOData() *odata.Query {
	return nil
}

func (c createOptions) ToQuery() *client.QueryParams {
	return nil
}
