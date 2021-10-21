package videoanalyzer

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type EdgeModulesDeleteResponse struct {
	HttpResponse *http.Response
}

// EdgeModulesDelete ...
func (c VideoAnalyzerClient) EdgeModulesDelete(ctx context.Context, id EdgeModuleId) (result EdgeModulesDeleteResponse, err error) {
	req, err := c.preparerForEdgeModulesDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "EdgeModulesDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "EdgeModulesDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForEdgeModulesDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "EdgeModulesDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForEdgeModulesDelete prepares the EdgeModulesDelete request.
func (c VideoAnalyzerClient) preparerForEdgeModulesDelete(ctx context.Context, id EdgeModuleId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForEdgeModulesDelete handles the response to the EdgeModulesDelete request. The method always
// closes the http.Response Body.
func (c VideoAnalyzerClient) responderForEdgeModulesDelete(resp *http.Response) (result EdgeModulesDeleteResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
