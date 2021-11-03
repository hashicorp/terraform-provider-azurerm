package videoanalyzer

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type VideosDeleteResponse struct {
	HttpResponse *http.Response
}

// VideosDelete ...
func (c VideoAnalyzerClient) VideosDelete(ctx context.Context, id VideoId) (result VideosDeleteResponse, err error) {
	req, err := c.preparerForVideosDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideosDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideosDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForVideosDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideosDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForVideosDelete prepares the VideosDelete request.
func (c VideoAnalyzerClient) preparerForVideosDelete(ctx context.Context, id VideoId) (*http.Request, error) {
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

// responderForVideosDelete handles the response to the VideosDelete request. The method always
// closes the http.Response Body.
func (c VideoAnalyzerClient) responderForVideosDelete(resp *http.Response) (result VideosDeleteResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
