package tenantconfiguration

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListOperationResponse struct {
	HttpResponse *http.Response
	Model        *ConfigurationList
}

// List ...
func (c TenantConfigurationClient) List(ctx context.Context) (result ListOperationResponse, err error) {
	req, err := c.preparerForList(ctx)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tenantconfiguration.TenantConfigurationClient", "List", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "tenantconfiguration.TenantConfigurationClient", "List", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForList(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tenantconfiguration.TenantConfigurationClient", "List", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForList prepares the List request.
func (c TenantConfigurationClient) preparerForList(ctx context.Context) (*http.Request, error) {
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

// responderForList handles the response to the List request. The method always
// closes the http.Response Body.
func (c TenantConfigurationClient) responderForList(resp *http.Response) (result ListOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
