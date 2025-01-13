package entities

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type QueryEntitiesInput struct {
	// An optional OData filter
	Filter *string

	// An optional comma-separated
	PropertyNamesToSelect *[]string

	// An optional OData top
	Top *int

	PartitionKey string
	RowKey       string

	// The Level of MetaData which should be returned
	MetaDataLevel MetaDataLevel

	// The Next Partition Key used to load data from a previous point
	NextPartitionKey *string

	// The Next Row Key used to load data from a previous point
	NextRowKey *string
}

type QueryEntitiesResponse struct {
	HttpResponse *http.Response

	NextPartitionKey string
	NextRowKey       string

	MetaData string                   `json:"odata.metadata,omitempty"`
	Entities []map[string]interface{} `json:"value"`
}

// Query queries entities in a table and includes the $filter and $select options.
func (c Client) Query(ctx context.Context, tableName string, input QueryEntitiesInput) (result QueryEntitiesResponse, err error) {
	if tableName == "" {
		return result, fmt.Errorf("`tableName` cannot be an empty string")
	}

	additionalParameters := make([]string, 0)
	if input.PartitionKey != "" {
		additionalParameters = append(additionalParameters, "PartitionKey='%s'", input.PartitionKey)
	}

	if input.RowKey != "" {
		additionalParameters = append(additionalParameters, "RowKey='%s'", input.RowKey)
	}

	path := fmt.Sprintf("/%s", tableName)
	if len(additionalParameters) > 0 {
		path += fmt.Sprintf("(%s)", strings.Join(additionalParameters, ","))
	}

	opts := client.RequestOptions{
		ContentType: "application/json",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		OptionsObject: queryOptions{
			input: input,
		},
		Path: path,
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
	input QueryEntitiesInput
}

func (q queryOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}
	headers.Append("Accept", fmt.Sprintf("application/json;odata=%s", q.input.MetaDataLevel))
	headers.Append("DataServiceVersion", "3.0;NetFx")
	headers.Append("MaxDataServiceVersion", "3.0;NetFx")
	return headers
}

func (q queryOptions) ToOData() *odata.Query {
	return nil
}

func (q queryOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}

	if q.input.Filter != nil {
		out.Append("$filter", *q.input.Filter)
	}

	if q.input.PropertyNamesToSelect != nil {
		out.Append("$select", strings.Join(*q.input.PropertyNamesToSelect, ","))
	}

	if q.input.Top != nil {
		out.Append("$top", strconv.Itoa(*q.input.Top))
	}

	if q.input.NextPartitionKey != nil {
		out.Append("NextPartitionKey", *q.input.NextPartitionKey)
	}

	if q.input.NextRowKey != nil {
		out.Append("NextRowKey", *q.input.NextRowKey)
	}
	return out
}
