package insights

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type WorkbooksRevisionGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *Workbook
}

// WorkbooksRevisionGet ...
func (c InsightsClient) WorkbooksRevisionGet(ctx context.Context, id RevisionId) (result WorkbooksRevisionGetOperationResponse, err error) {
	req, err := c.preparerForWorkbooksRevisionGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "insights.InsightsClient", "WorkbooksRevisionGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "insights.InsightsClient", "WorkbooksRevisionGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForWorkbooksRevisionGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "insights.InsightsClient", "WorkbooksRevisionGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForWorkbooksRevisionGet prepares the WorkbooksRevisionGet request.
func (c InsightsClient) preparerForWorkbooksRevisionGet(ctx context.Context, id RevisionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForWorkbooksRevisionGet handles the response to the WorkbooksRevisionGet request. The method always
// closes the http.Response Body.
func (c InsightsClient) responderForWorkbooksRevisionGet(resp *http.Response) (result WorkbooksRevisionGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
