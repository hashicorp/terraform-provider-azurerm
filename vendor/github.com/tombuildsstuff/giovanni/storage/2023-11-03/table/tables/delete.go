package tables

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type DeleteTableResponse struct {
	HttpResponse *client.Response
}

// Delete deletes the specified table and any data it contains.
func (c Client) Delete(ctx context.Context, tableName string) (resp DeleteTableResponse, err error) {
	if tableName == "" {
		return resp, fmt.Errorf("`tableName` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ContentType: "application/json",
		ExpectedStatusCodes: []int{
			http.StatusNoContent,
		},
		HttpMethod:    http.MethodDelete,
		OptionsObject: deleteOptions{},
		Path:          fmt.Sprintf("/Tables('%s')", tableName),
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
}

func (d deleteOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}
	headers.Append("Accept", "application/json")
	return headers
}

func (d deleteOptions) ToOData() *odata.Query {
	return nil
}

func (d deleteOptions) ToQuery() *client.QueryParams {
	return nil
}
