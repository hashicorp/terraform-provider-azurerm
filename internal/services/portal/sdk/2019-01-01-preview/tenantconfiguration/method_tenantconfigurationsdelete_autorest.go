package tenantconfiguration

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type TenantConfigurationsDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// TenantConfigurationsDelete ...
func (c TenantConfigurationClient) TenantConfigurationsDelete(ctx context.Context, id ConfigurationId) (result TenantConfigurationsDeleteOperationResponse, err error) {
	req, err := c.preparerForTenantConfigurationsDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tenantconfiguration.TenantConfigurationClient", "TenantConfigurationsDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "tenantconfiguration.TenantConfigurationClient", "TenantConfigurationsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForTenantConfigurationsDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tenantconfiguration.TenantConfigurationClient", "TenantConfigurationsDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForTenantConfigurationsDelete prepares the TenantConfigurationsDelete request.
func (c TenantConfigurationClient) preparerForTenantConfigurationsDelete(ctx context.Context, id ConfigurationId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForTenantConfigurationsDelete handles the response to the TenantConfigurationsDelete request. The method always
// closes the http.Response Body.
func (c TenantConfigurationClient) responderForTenantConfigurationsDelete(resp *http.Response) (result TenantConfigurationsDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
