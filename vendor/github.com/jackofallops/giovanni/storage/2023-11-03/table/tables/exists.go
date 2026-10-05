package tables

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type TableExistsResponse struct {
	HttpResponse *http.Response
}

// Exists checks that the specified table exists
func (c Client) Exists(ctx context.Context, tableName string) (result TableExistsResponse, err error) {
	if tableName == "" {
		err = fmt.Errorf("`tableName` cannot be an empty string")
		return
	}

	opts := client.RequestOptions{
		ContentType: "application/json",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: tableExistsOptions{},
		Path:          fmt.Sprintf("/Tables('%s')", tableName),
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

type tableExistsOptions struct{}

func (t tableExistsOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}
	headers.Append("Accept", "application/json;odata=nometadata")
	return headers
}

func (t tableExistsOptions) ToOData() *odata.Query {
	return nil
}

func (t tableExistsOptions) ToQuery() *client.QueryParams {
	return nil
}
