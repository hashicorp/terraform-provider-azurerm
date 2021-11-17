package partnertopics

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type DeactivateResponse struct {
	HttpResponse *http.Response
	Model        *PartnerTopic
}

// Deactivate ...
func (c PartnerTopicsClient) Deactivate(ctx context.Context, id PartnerTopicId) (result DeactivateResponse, err error) {
	req, err := c.preparerForDeactivate(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "partnertopics.PartnerTopicsClient", "Deactivate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "partnertopics.PartnerTopicsClient", "Deactivate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDeactivate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "partnertopics.PartnerTopicsClient", "Deactivate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDeactivate prepares the Deactivate request.
func (c PartnerTopicsClient) preparerForDeactivate(ctx context.Context, id PartnerTopicId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/deactivate", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDeactivate handles the response to the Deactivate request. The method always
// closes the http.Response Body.
func (c PartnerTopicsClient) responderForDeactivate(resp *http.Response) (result DeactivateResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
