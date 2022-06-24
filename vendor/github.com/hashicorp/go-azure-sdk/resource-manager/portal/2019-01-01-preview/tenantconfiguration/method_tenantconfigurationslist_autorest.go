package tenantconfiguration

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TenantConfigurationsListOperationResponse struct {
	HttpResponse *http.Response
	Model        *ConfigurationList
}

// TenantConfigurationsList ...
func (c TenantConfigurationClient) TenantConfigurationsList(ctx context.Context) (result TenantConfigurationsListOperationResponse, err error) {
	req, err := c.preparerForTenantConfigurationsList(ctx)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tenantconfiguration.TenantConfigurationClient", "TenantConfigurationsList", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "tenantconfiguration.TenantConfigurationClient", "TenantConfigurationsList", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForTenantConfigurationsList(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tenantconfiguration.TenantConfigurationClient", "TenantConfigurationsList", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForTenantConfigurationsList prepares the TenantConfigurationsList request.
func (c TenantConfigurationClient) preparerForTenantConfigurationsList(ctx context.Context) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath("/providers/Microsoft.Portal/tenantConfigurations"),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForTenantConfigurationsList handles the response to the TenantConfigurationsList request. The method always
// closes the http.Response Body.
func (c TenantConfigurationClient) responderForTenantConfigurationsList(resp *http.Response) (result TenantConfigurationsListOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
