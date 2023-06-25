package maintenanceconfigurations

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ForResourceGroupListOperationResponse struct {
	HttpResponse *http.Response
	Model        *ListMaintenanceConfigurationsResult
}

// ForResourceGroupList ...
func (c MaintenanceConfigurationsClient) ForResourceGroupList(ctx context.Context, id commonids.ResourceGroupId) (result ForResourceGroupListOperationResponse, err error) {
	req, err := c.preparerForForResourceGroupList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "maintenanceconfigurations.MaintenanceConfigurationsClient", "ForResourceGroupList", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "maintenanceconfigurations.MaintenanceConfigurationsClient", "ForResourceGroupList", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForForResourceGroupList(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "maintenanceconfigurations.MaintenanceConfigurationsClient", "ForResourceGroupList", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForForResourceGroupList prepares the ForResourceGroupList request.
func (c MaintenanceConfigurationsClient) preparerForForResourceGroupList(ctx context.Context, id commonids.ResourceGroupId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Maintenance/maintenanceConfigurations", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForForResourceGroupList handles the response to the ForResourceGroupList request. The method always
// closes the http.Response Body.
func (c MaintenanceConfigurationsClient) responderForForResourceGroupList(resp *http.Response) (result ForResourceGroupListOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
