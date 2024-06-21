package entities

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type GetEntityInput struct {
	PartitionKey string
	RowKey       string

	// The Level of MetaData which should be returned
	MetaDataLevel MetaDataLevel
}

type GetEntityResponse struct {
	HttpResponse *http.Response

	Entity map[string]interface{}
}

// Get queries entities in a table and includes the $filter and $select options.
func (c Client) Get(ctx context.Context, tableName string, input GetEntityInput) (result GetEntityResponse, err error) {
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
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		OptionsObject: getEntitiesOptions{
			MetaDataLevel: input.MetaDataLevel,
		},
		Path: fmt.Sprintf("/%s(PartitionKey='%s', RowKey='%s')", tableName, input.PartitionKey, input.RowKey),
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
			err = resp.Unmarshal(&result.Entity)
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

type getEntitiesOptions struct {
	MetaDataLevel MetaDataLevel
}

func (g getEntitiesOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}
	headers.Append("Accept", fmt.Sprintf("application/json;odata=%s", g.MetaDataLevel))
	headers.Append("DataServiceVersion", "3.0;NetFx")
	headers.Append("MaxDataServiceVersion", "3.0;NetFx")
	return headers
}

func (g getEntitiesOptions) ToOData() *odata.Query {
	return nil
}

func (g getEntitiesOptions) ToQuery() *client.QueryParams {
	return nil
}
