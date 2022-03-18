package customapis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type CustomApisListWsdlInterfacesOperationResponse struct {
	HttpResponse *http.Response
	Model        *WsdlServiceCollection
}

// CustomApisListWsdlInterfaces ...
func (c CustomAPIsClient) CustomApisListWsdlInterfaces(ctx context.Context, id LocationId, input WsdlDefinition) (result CustomApisListWsdlInterfacesOperationResponse, err error) {
	req, err := c.preparerForCustomApisListWsdlInterfaces(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "customapis.CustomAPIsClient", "CustomApisListWsdlInterfaces", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "customapis.CustomAPIsClient", "CustomApisListWsdlInterfaces", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCustomApisListWsdlInterfaces(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "customapis.CustomAPIsClient", "CustomApisListWsdlInterfaces", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCustomApisListWsdlInterfaces prepares the CustomApisListWsdlInterfaces request.
func (c CustomAPIsClient) preparerForCustomApisListWsdlInterfaces(ctx context.Context, id LocationId, input WsdlDefinition) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listWsdlInterfaces", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCustomApisListWsdlInterfaces handles the response to the CustomApisListWsdlInterfaces request. The method always
// closes the http.Response Body.
func (c CustomAPIsClient) responderForCustomApisListWsdlInterfaces(resp *http.Response) (result CustomApisListWsdlInterfacesOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
