package topics

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ExtensionTopicsGetResponse struct {
	HttpResponse *http.Response
	Model        *ExtensionTopic
}

// ExtensionTopicsGet ...
func (c TopicsClient) ExtensionTopicsGet(ctx context.Context, id ScopeId) (result ExtensionTopicsGetResponse, err error) {
	req, err := c.preparerForExtensionTopicsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "topics.TopicsClient", "ExtensionTopicsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "topics.TopicsClient", "ExtensionTopicsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForExtensionTopicsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "topics.TopicsClient", "ExtensionTopicsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForExtensionTopicsGet prepares the ExtensionTopicsGet request.
func (c TopicsClient) preparerForExtensionTopicsGet(ctx context.Context, id ScopeId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.EventGrid/extensionTopics/default", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForExtensionTopicsGet handles the response to the ExtensionTopicsGet request. The method always
// closes the http.Response Body.
func (c TopicsClient) responderForExtensionTopicsGet(resp *http.Response) (result ExtensionTopicsGetResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
