package tables

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type createTableRequest struct {
	TableName string `json:"TableName"`
}

type CreateTableResponse struct {
	HttpResponse *client.Response
}

// Create creates a new table in the storage account.
func (c Client) Create(ctx context.Context, tableName string) (resp CreateTableResponse, err error) {
	if tableName == "" {
		return resp, fmt.Errorf("`tableName` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ContentType: "application/json",
		ExpectedStatusCodes: []int{
			http.StatusNoContent,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: createTableOptions{},
		Path:          "/Tables",
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		err = fmt.Errorf("building request: %+v", err)
		return
	}

	err = req.Marshal(&createTableRequest{TableName: tableName})
	if err != nil {
		return resp, fmt.Errorf("marshalling request")
	}

	resp.HttpResponse, err = req.Execute(ctx)
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	return
}

type createTableOptions struct{}

func (c createTableOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}
	headers.Append("Accept", "application/json;odata=nometadata")
	headers.Append("Prefer", "return-no-content")
	return headers
}

func (c createTableOptions) ToOData() *odata.Query {
	return nil
}

func (c createTableOptions) ToQuery() *client.QueryParams {
	return nil
}
