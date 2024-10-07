package tables

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type GetResponse struct {
	HttpResponse *http.Response

	MetaData string          `json:"odata.metadata,omitempty"`
	Tables   []GetResultItem `json:"value"`
}

type QueryInput struct {
	MetaDataLevel MetaDataLevel
}

// Query returns a list of tables under the specified account.
func (c Client) Query(ctx context.Context, input QueryInput) (result GetResponse, err error) {

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

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil && resp.Response != nil {
		result.HttpResponse = resp.Response

		if err == nil {
			err = resp.Unmarshal(&result)
			if err != nil {
				err = fmt.Errorf("unmarshalling response: %+v", err)
				return
			}
		}
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
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
