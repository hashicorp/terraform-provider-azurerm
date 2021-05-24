package entities

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
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

type QueryEntitiesResult struct {
	autorest.Response

	NextPartitionKey string
	NextRowKey       string

	MetaData string                   `json:"odata.metadata,omitempty"`
	Entities []map[string]interface{} `json:"value"`
}

// Query queries entities in a table and includes the $filter and $select options.
func (client Client) Query(ctx context.Context, accountName, tableName string, input QueryEntitiesInput) (result QueryEntitiesResult, err error) {
	if accountName == "" {
		return result, validation.NewError("entities.Client", "Query", "`accountName` cannot be an empty string.")
	}
	if tableName == "" {
		return result, validation.NewError("entities.Client", "Query", "`tableName` cannot be an empty string.")
	}

	req, err := client.QueryPreparer(ctx, accountName, tableName, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "entities.Client", "Query", nil, "Failure preparing request")
		return
	}

	resp, err := client.QuerySender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "entities.Client", "Query", resp, "Failure sending request")
		return
	}

	result, err = client.QueryResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "entities.Client", "Query", resp, "Failure responding to request")
		return
	}

	return
}

// QueryPreparer prepares the Query request.
func (client Client) QueryPreparer(ctx context.Context, accountName, tableName string, input QueryEntitiesInput) (*http.Request, error) {

	pathParameters := map[string]interface{}{
		"tableName":            autorest.Encode("path", tableName),
		"additionalParameters": "",
	}

	//PartitionKey='<partition-key>',RowKey='<row-key>'
	additionalParams := make([]string, 0)
	if input.PartitionKey != "" {
		additionalParams = append(additionalParams, fmt.Sprintf("PartitionKey='%s'", input.PartitionKey))
	}
	if input.RowKey != "" {
		additionalParams = append(additionalParams, fmt.Sprintf("RowKey='%s'", input.RowKey))
	}
	if len(additionalParams) > 0 {
		pathParameters["additionalParameters"] = autorest.Encode("path", strings.Join(additionalParams, ","))
	}

	queryParameters := map[string]interface{}{}

	if input.Filter != nil {
		queryParameters["$filter"] = autorest.Encode("query", *input.Filter)
	}

	if input.PropertyNamesToSelect != nil {
		queryParameters["$select"] = autorest.Encode("query", strings.Join(*input.PropertyNamesToSelect, ","))
	}

	if input.Top != nil {
		queryParameters["$top"] = autorest.Encode("query", *input.Top)
	}

	if input.NextPartitionKey != nil {
		queryParameters["NextPartitionKey"] = *input.NextPartitionKey
	}

	if input.NextRowKey != nil {
		queryParameters["NextRowKey"] = *input.NextRowKey
	}

	headers := map[string]interface{}{
		"x-ms-version":          APIVersion,
		"Accept":                fmt.Sprintf("application/json;odata=%s", input.MetaDataLevel),
		"DataServiceVersion":    "3.0;NetFx",
		"MaxDataServiceVersion": "3.0;NetFx",
	}

	// GET /myaccount/Customers()?$filter=(Rating%20ge%203)%20and%20(Rating%20le%206)&$select=PartitionKey,RowKey,Address,CustomerSince
	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(endpoints.GetTableEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{tableName}({additionalParameters})", pathParameters),
		autorest.WithQueryParameters(queryParameters),
		autorest.WithHeaders(headers))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// QuerySender sends the Query request. The method will close the
// http.Response Body if it receives an error.
func (client Client) QuerySender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// QueryResponder handles the response to the Query request. The method always
// closes the http.Response Body.
func (client Client) QueryResponder(resp *http.Response) (result QueryEntitiesResult, err error) {
	if resp != nil && resp.Header != nil {
		result.NextPartitionKey = resp.Header.Get("x-ms-continuation-NextPartitionKey")
		result.NextRowKey = resp.Header.Get("x-ms-continuation-NextRowKey")
	}

	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}

	return
}
