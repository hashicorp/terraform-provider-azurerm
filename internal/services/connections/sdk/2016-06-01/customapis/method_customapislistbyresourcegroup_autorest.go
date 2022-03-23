package customapis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

type CustomApisListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *CustomApiDefinitionCollection
}

type CustomApisListByResourceGroupOperationOptions struct {
	Skiptoken *string
	Top       *int64
}

func DefaultCustomApisListByResourceGroupOperationOptions() CustomApisListByResourceGroupOperationOptions {
	return CustomApisListByResourceGroupOperationOptions{}
}

func (o CustomApisListByResourceGroupOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o CustomApisListByResourceGroupOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Skiptoken != nil {
		out["skiptoken"] = *o.Skiptoken
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// CustomApisListByResourceGroup ...
func (c CustomAPIsClient) CustomApisListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId, options CustomApisListByResourceGroupOperationOptions) (result CustomApisListByResourceGroupOperationResponse, err error) {
	req, err := c.preparerForCustomApisListByResourceGroup(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "customapis.CustomAPIsClient", "CustomApisListByResourceGroup", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "customapis.CustomAPIsClient", "CustomApisListByResourceGroup", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCustomApisListByResourceGroup(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "customapis.CustomAPIsClient", "CustomApisListByResourceGroup", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCustomApisListByResourceGroup prepares the CustomApisListByResourceGroup request.
func (c CustomAPIsClient) preparerForCustomApisListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId, options CustomApisListByResourceGroupOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Web/customApis", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCustomApisListByResourceGroup handles the response to the CustomApisListByResourceGroup request. The method always
// closes the http.Response Body.
func (c CustomAPIsClient) responderForCustomApisListByResourceGroup(resp *http.Response) (result CustomApisListByResourceGroupOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
