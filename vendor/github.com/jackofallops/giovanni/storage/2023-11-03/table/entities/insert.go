package entities

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type InsertEntityInput struct {
	// The level of MetaData provided for this Entity
	MetaDataLevel MetaDataLevel

	// The Entity which should be inserted, by default all values are strings
	// To explicitly type a property, specify the appropriate OData data type by setting
	// the m:type attribute within the property definition
	Entity map[string]interface{}

	// When inserting an entity into a table, you must specify values for the PartitionKey and RowKey system properties.
	// Together, these properties form the primary key and must be unique within the table.
	// Both the PartitionKey and RowKey values must be string values; each key value may be up to 64 KB in size.
	// If you are using an integer value for the key value, you should convert the integer to a fixed-width string,
	// because they are canonically sorted. For example, you should convert the value 1 to 0000001 to ensure proper sorting.
	RowKey       string
	PartitionKey string
}

type InsertResponse struct {
	HttpResponse *http.Response
}

// Insert inserts a new entity into a table.
func (c Client) Insert(ctx context.Context, tableName string, input InsertEntityInput) (result InsertResponse, err error) {
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
		HttpMethod: http.MethodPost,
		OptionsObject: insertOptions{
			MetaDataLevel: input.MetaDataLevel,
		},
		Path: fmt.Sprintf("/%s", tableName),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		err = fmt.Errorf("building request: %+v", err)
		return
	}

	input.Entity["PartitionKey"] = input.PartitionKey
	input.Entity["RowKey"] = input.RowKey

	err = req.Marshal(&input.Entity)
	if err != nil {
		return result, fmt.Errorf("marshalling request: %+v", err)
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

type insertOptions struct {
	MetaDataLevel MetaDataLevel
}

func (i insertOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}
	headers.Append("Accept", fmt.Sprintf("application/json;odata=%s", i.MetaDataLevel))
	headers.Append("Prefer", "return-no-content")
	return headers
}

func (i insertOptions) ToOData() *odata.Query {
	return nil
}

func (i insertOptions) ToQuery() *client.QueryParams {
	return nil
}
