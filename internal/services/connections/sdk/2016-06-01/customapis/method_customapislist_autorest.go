package customapis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

type CustomApisListOperationResponse struct {
	HttpResponse *http.Response
	Model        *CustomApiDefinitionCollection
}

type CustomApisListOperationOptions struct {
	Skiptoken *string
	Top       *int64
}

func DefaultCustomApisListOperationOptions() CustomApisListOperationOptions {
	return CustomApisListOperationOptions{}
}

func (o CustomApisListOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o CustomApisListOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Skiptoken != nil {
		out["skiptoken"] = *o.Skiptoken
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// CustomApisList ...
func (c CustomAPIsClient) CustomApisList(ctx context.Context, id commonids.SubscriptionId, options CustomApisListOperationOptions) (result CustomApisListOperationResponse, err error) {
	req, err := c.preparerForCustomApisList(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "customapis.CustomAPIsClient", "CustomApisList", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "customapis.CustomAPIsClient", "CustomApisList", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCustomApisList(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "customapis.CustomAPIsClient", "CustomApisList", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCustomApisList prepares the CustomApisList request.
func (c CustomAPIsClient) preparerForCustomApisList(ctx context.Context, id commonids.SubscriptionId, options CustomApisListOperationOptions) (*http.Request, error) {
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

// responderForCustomApisList handles the response to the CustomApisList request. The method always
// closes the http.Response Body.
func (c CustomAPIsClient) responderForCustomApisList(resp *http.Response) (result CustomApisListOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
