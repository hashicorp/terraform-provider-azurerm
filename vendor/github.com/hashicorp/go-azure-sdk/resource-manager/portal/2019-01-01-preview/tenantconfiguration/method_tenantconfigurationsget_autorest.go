package tenantconfiguration

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TenantConfigurationsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *Configuration
}

// TenantConfigurationsGet ...
func (c TenantConfigurationClient) TenantConfigurationsGet(ctx context.Context) (result TenantConfigurationsGetOperationResponse, err error) {
	req, err := c.preparerForTenantConfigurationsGet(ctx)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tenantconfiguration.TenantConfigurationClient", "TenantConfigurationsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "tenantconfiguration.TenantConfigurationClient", "TenantConfigurationsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForTenantConfigurationsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tenantconfiguration.TenantConfigurationClient", "TenantConfigurationsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForTenantConfigurationsGet prepares the TenantConfigurationsGet request.
func (c TenantConfigurationClient) preparerForTenantConfigurationsGet(ctx context.Context) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath("/providers/Microsoft.Portal/tenantConfigurations/default"),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForTenantConfigurationsGet handles the response to the TenantConfigurationsGet request. The method always
// closes the http.Response Body.
func (c TenantConfigurationClient) responderForTenantConfigurationsGet(resp *http.Response) (result TenantConfigurationsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
