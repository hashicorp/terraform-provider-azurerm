package tables

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
)

type GetResult struct {
	autorest.Response

	MetaData string          `json:"odata.metadata,omitempty"`
	Tables   []GetResultItem `json:"value"`
}

// Query returns a list of tables under the specified account.
func (client Client) Query(ctx context.Context, accountName string, metaDataLevel MetaDataLevel) (result GetResult, err error) {
	if accountName == "" {
		return result, validation.NewError("tables.Client", "Query", "`accountName` cannot be an empty string.")
	}

	req, err := client.QueryPreparer(ctx, accountName, metaDataLevel)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tables.Client", "Query", nil, "Failure preparing request")
		return
	}

	resp, err := client.QuerySender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "tables.Client", "Query", resp, "Failure sending request")
		return
	}

	result, err = client.QueryResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tables.Client", "Query", resp, "Failure responding to request")
		return
	}

	return
}

// QueryPreparer prepares the Query request.
func (client Client) QueryPreparer(ctx context.Context, accountName string, metaDataLevel MetaDataLevel) (*http.Request, error) {
	// NOTE: whilst this supports ContinuationTokens and 'Top'
	// it appears that 'Skip' returns a '501 Not Implemented'
	// as such, we intentionally don't support those right now

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
		"Accept":       fmt.Sprintf("application/json;odata=%s", metaDataLevel),
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(endpoints.GetTableEndpoint(client.BaseURI, accountName)),
		autorest.WithPath("/Tables"),
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
func (client Client) QueryResponder(resp *http.Response) (result GetResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}

	return
}
