package entities

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type DeleteEntityInput struct {
	// When inserting an entity into a table, you must specify values for the PartitionKey and RowKey system properties.
	// Together, these properties form the primary key and must be unique within the table.
	// Both the PartitionKey and RowKey values must be string values; each key value may be up to 64 KB in size.
	// If you are using an integer value for the key value, you should convert the integer to a fixed-width string,
	// because they are canonically sorted. For example, you should convert the value 1 to 0000001 to ensure proper sorting.
	RowKey       string
	PartitionKey string
}

type DeleteEntityResponse struct {
	HttpResponse *http.Response
}

// Delete deletes an existing entity in a table.
func (c Client) Delete(ctx context.Context, tableName string, input DeleteEntityInput) (result DeleteEntityResponse, err error) {

	if tableName == "" {
		return result, fmt.Errorf("`tableName` cannot be an empty string")
	}

	if input.PartitionKey == "" {
		return result, fmt.Errorf("`input.PartitionKey` cannot be an empty string")
	}

	if input.RowKey == "" {
		return result, fmt.Errorf("`input.RowKey` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ContentType: "application/json",
		ExpectedStatusCodes: []int{
			http.StatusNoContent,
		},
		HttpMethod:    http.MethodDelete,
		OptionsObject: deleteEntitiesOptions{},
		Path:          fmt.Sprintf("/%s(PartitionKey='%s', RowKey='%s')", tableName, input.PartitionKey, input.RowKey),
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

type deleteEntitiesOptions struct{}

func (d deleteEntitiesOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}
	headers.Append("Accept", "application/json")
	headers.Append("If-Match", "*")
	return headers
}

func (d deleteEntitiesOptions) ToOData() *odata.Query {
	return nil
}

func (d deleteEntitiesOptions) ToQuery() *client.QueryParams {
	return nil
}
