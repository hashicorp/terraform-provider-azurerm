package dimensions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ByExternalCloudProviderTypeResponse struct {
	HttpResponse *http.Response
	Model        *Dimension
}

type ByExternalCloudProviderTypeOptions struct {
	Expand *string
	Filter *string
	Top    *int64
}

func DefaultByExternalCloudProviderTypeOptions() ByExternalCloudProviderTypeOptions {
	return ByExternalCloudProviderTypeOptions{}
}

func (o ByExternalCloudProviderTypeOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Expand != nil {
		out["$expand"] = *o.Expand
	}

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// ByExternalCloudProviderType ...
func (c DimensionsClient) ByExternalCloudProviderType(ctx context.Context, id ExternalCloudProviderTypeId, options ByExternalCloudProviderTypeOptions) (result ByExternalCloudProviderTypeResponse, err error) {
	req, err := c.preparerForByExternalCloudProviderType(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dimensions.DimensionsClient", "ByExternalCloudProviderType", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "dimensions.DimensionsClient", "ByExternalCloudProviderType", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForByExternalCloudProviderType(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dimensions.DimensionsClient", "ByExternalCloudProviderType", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForByExternalCloudProviderType prepares the ByExternalCloudProviderType request.
func (c DimensionsClient) preparerForByExternalCloudProviderType(ctx context.Context, id ExternalCloudProviderTypeId, options ByExternalCloudProviderTypeOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/dimensions", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForByExternalCloudProviderType handles the response to the ByExternalCloudProviderType request. The method always
// closes the http.Response Body.
func (c DimensionsClient) responderForByExternalCloudProviderType(resp *http.Response) (result ByExternalCloudProviderTypeResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
