package tables

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type GetResponse struct {
	HttpResponse *client.Response

	MetaData string          `json:"odata.metadata,omitempty"`
	Tables   []GetResultItem `json:"value"`
}

type QueryInput struct {
	MetaDataLevel MetaDataLevel
}

// Query returns a list of tables under the specified account.
func (c Client) Query(ctx context.Context, input QueryInput) (resp GetResponse, err error) {

	opts := client.RequestOptions{
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		OptionsObject: queryOptions{
			metaDataLevel: input.MetaDataLevel,
		},
		Path: "/Tables",
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
		if resp.HttpResponse.Body != nil {
			err = resp.HttpResponse.Unmarshal(&resp)
			if err != nil {
				return resp, fmt.Errorf("unmarshalling response: %v", err)
			}
		}
	}

	return
}

type queryOptions struct {
	metaDataLevel MetaDataLevel
}

func (q queryOptions) ToHeaders() *client.Headers {
	// NOTE: whilst this supports ContinuationTokens and 'Top'
	// it appears that 'Skip' returns a '501 Not Implemented'
	// as such, we intentionally don't support those right now
	headers := &client.Headers{}
	headers.Append("Accept", fmt.Sprintf("application/json;odata=%s", q.metaDataLevel))
	return headers
}

func (q queryOptions) ToOData() *odata.Query {
	return nil
}

func (q queryOptions) ToQuery() *client.QueryParams {
	return nil
}
