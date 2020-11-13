package tables

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
)

type createTableRequest struct {
	TableName string `json:"TableName"`
}

// Create creates a new table in the storage account.
func (client Client) Create(ctx context.Context, accountName, tableName string) (result autorest.Response, err error) {
	if accountName == "" {
		return result, validation.NewError("tables.Client", "Create", "`accountName` cannot be an empty string.")
	}
	if tableName == "" {
		return result, validation.NewError("tables.Client", "Create", "`tableName` cannot be an empty string.")
	}

	req, err := client.CreatePreparer(ctx, accountName, tableName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tables.Client", "Create", nil, "Failure preparing request")
		return
	}

	resp, err := client.CreateSender(req)
	if err != nil {
		result = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "tables.Client", "Create", resp, "Failure sending request")
		return
	}

	result, err = client.CreateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tables.Client", "Create", resp, "Failure responding to request")
		return
	}

	return
}

// CreatePreparer prepares the Create request.
func (client Client) CreatePreparer(ctx context.Context, accountName, tableName string) (*http.Request, error) {
	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
		// NOTE: we could support returning metadata here, but it doesn't appear to be directly useful
		// vs making a request using the Get methods as-necessary?
		"Accept": "application/json;odata=nometadata",
		"Prefer": "return-no-content",
	}

	body := createTableRequest{
		TableName: tableName,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json"),
		autorest.AsPost(),
		autorest.WithBaseURL(endpoints.GetTableEndpoint(client.BaseURI, accountName)),
		autorest.WithPath("/Tables"),
		autorest.WithJSON(body),
		autorest.WithHeaders(headers))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateSender sends the Create request. The method will close the
// http.Response Body if it receives an error.
func (client Client) CreateSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// CreateResponder handles the response to the Create request. The method always
// closes the http.Response Body.
func (client Client) CreateResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusNoContent),
		autorest.ByClosing())
	result = autorest.Response{Response: resp}

	return
}
