package partnertopics

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ActivateResponse struct {
	HttpResponse *http.Response
	Model        *PartnerTopic
}

// Activate ...
func (c PartnerTopicsClient) Activate(ctx context.Context, id PartnerTopicId) (result ActivateResponse, err error) {
	req, err := c.preparerForActivate(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "partnertopics.PartnerTopicsClient", "Activate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "partnertopics.PartnerTopicsClient", "Activate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForActivate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "partnertopics.PartnerTopicsClient", "Activate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForActivate prepares the Activate request.
func (c PartnerTopicsClient) preparerForActivate(ctx context.Context, id PartnerTopicId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/activate", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForActivate handles the response to the Activate request. The method always
// closes the http.Response Body.
func (c PartnerTopicsClient) responderForActivate(resp *http.Response) (result ActivateResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
