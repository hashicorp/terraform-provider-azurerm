package customapis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type CustomApisExtractApiDefinitionFromWsdlOperationResponse struct {
	HttpResponse *http.Response
	Model        *interface{}
}

// CustomApisExtractApiDefinitionFromWsdl ...
func (c CustomAPIsClient) CustomApisExtractApiDefinitionFromWsdl(ctx context.Context, id LocationId, input WsdlDefinition) (result CustomApisExtractApiDefinitionFromWsdlOperationResponse, err error) {
	req, err := c.preparerForCustomApisExtractApiDefinitionFromWsdl(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "customapis.CustomAPIsClient", "CustomApisExtractApiDefinitionFromWsdl", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "customapis.CustomAPIsClient", "CustomApisExtractApiDefinitionFromWsdl", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCustomApisExtractApiDefinitionFromWsdl(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "customapis.CustomAPIsClient", "CustomApisExtractApiDefinitionFromWsdl", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCustomApisExtractApiDefinitionFromWsdl prepares the CustomApisExtractApiDefinitionFromWsdl request.
func (c CustomAPIsClient) preparerForCustomApisExtractApiDefinitionFromWsdl(ctx context.Context, id LocationId, input WsdlDefinition) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/extractApiDefinitionFromWsdl", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCustomApisExtractApiDefinitionFromWsdl handles the response to the CustomApisExtractApiDefinitionFromWsdl request. The method always
// closes the http.Response Body.
func (c CustomAPIsClient) responderForCustomApisExtractApiDefinitionFromWsdl(resp *http.Response) (result CustomApisExtractApiDefinitionFromWsdlOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
