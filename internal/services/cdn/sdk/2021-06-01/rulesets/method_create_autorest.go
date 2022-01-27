package rulesets

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type CreateResponse struct {
	HttpResponse *http.Response
	Model        *RuleSet
}

// Create ...
func (c RuleSetsClient) Create(ctx context.Context, id RuleSetId) (result CreateResponse, err error) {
	req, err := c.preparerForCreate(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "rulesets.RuleSetsClient", "Create", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "rulesets.RuleSetsClient", "Create", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCreate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "rulesets.RuleSetsClient", "Create", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCreate prepares the Create request.
func (c RuleSetsClient) preparerForCreate(ctx context.Context, id RuleSetId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCreate handles the response to the Create request. The method always
// closes the http.Response Body.
func (c RuleSetsClient) responderForCreate(resp *http.Response) (result CreateResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
