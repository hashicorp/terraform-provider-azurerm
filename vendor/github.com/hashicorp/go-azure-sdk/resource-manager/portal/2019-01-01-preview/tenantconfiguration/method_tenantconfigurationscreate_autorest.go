package tenantconfiguration

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TenantConfigurationsCreateOperationResponse struct {
	HttpResponse *http.Response
	Model        *Configuration
}

// TenantConfigurationsCreate ...
func (c TenantConfigurationClient) TenantConfigurationsCreate(ctx context.Context, input Configuration) (result TenantConfigurationsCreateOperationResponse, err error) {
	req, err := c.preparerForTenantConfigurationsCreate(ctx, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tenantconfiguration.TenantConfigurationClient", "TenantConfigurationsCreate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "tenantconfiguration.TenantConfigurationClient", "TenantConfigurationsCreate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForTenantConfigurationsCreate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tenantconfiguration.TenantConfigurationClient", "TenantConfigurationsCreate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForTenantConfigurationsCreate prepares the TenantConfigurationsCreate request.
func (c TenantConfigurationClient) preparerForTenantConfigurationsCreate(ctx context.Context, input Configuration) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath("/providers/Microsoft.Portal/tenantConfigurations/default"),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForTenantConfigurationsCreate handles the response to the TenantConfigurationsCreate request. The method always
// closes the http.Response Body.
func (c TenantConfigurationClient) responderForTenantConfigurationsCreate(resp *http.Response) (result TenantConfigurationsCreateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
