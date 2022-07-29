package tenantconfiguration

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TenantConfigurationsDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// TenantConfigurationsDelete ...
func (c TenantConfigurationClient) TenantConfigurationsDelete(ctx context.Context) (result TenantConfigurationsDeleteOperationResponse, err error) {
	req, err := c.preparerForTenantConfigurationsDelete(ctx)
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
func (c TenantConfigurationClient) preparerForTenantConfigurationsDelete(ctx context.Context) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath("/providers/Microsoft.Portal/tenantConfigurations/default"),
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
