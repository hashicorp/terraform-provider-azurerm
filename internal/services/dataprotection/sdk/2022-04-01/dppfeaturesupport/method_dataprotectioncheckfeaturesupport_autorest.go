package dppfeaturesupport

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

type DataProtectionCheckFeatureSupportOperationResponse struct {
	HttpResponse *http.Response
	Model        *FeatureValidationResponseBase
}

// DataProtectionCheckFeatureSupport ...
func (c DppFeatureSupportClient) DataProtectionCheckFeatureSupport(ctx context.Context, id LocationId, input FeatureValidationRequestBase) (result DataProtectionCheckFeatureSupportOperationResponse, err error) {
	req, err := c.preparerForDataProtectionCheckFeatureSupport(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dppfeaturesupport.DppFeatureSupportClient", "DataProtectionCheckFeatureSupport", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "dppfeaturesupport.DppFeatureSupportClient", "DataProtectionCheckFeatureSupport", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDataProtectionCheckFeatureSupport(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dppfeaturesupport.DppFeatureSupportClient", "DataProtectionCheckFeatureSupport", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDataProtectionCheckFeatureSupport prepares the DataProtectionCheckFeatureSupport request.
func (c DppFeatureSupportClient) preparerForDataProtectionCheckFeatureSupport(ctx context.Context, id LocationId, input FeatureValidationRequestBase) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/checkFeatureSupport", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDataProtectionCheckFeatureSupport handles the response to the DataProtectionCheckFeatureSupport request. The method always
// closes the http.Response Body.
func (c DppFeatureSupportClient) responderForDataProtectionCheckFeatureSupport(resp *http.Response) (result DataProtectionCheckFeatureSupportOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
