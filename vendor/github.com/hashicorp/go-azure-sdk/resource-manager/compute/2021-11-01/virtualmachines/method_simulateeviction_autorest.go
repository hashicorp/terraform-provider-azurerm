package virtualmachines

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SimulateEvictionOperationResponse struct {
	HttpResponse *http.Response
}

// SimulateEviction ...
func (c VirtualMachinesClient) SimulateEviction(ctx context.Context, id VirtualMachineId) (result SimulateEvictionOperationResponse, err error) {
	req, err := c.preparerForSimulateEviction(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "SimulateEviction", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "SimulateEviction", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSimulateEviction(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "SimulateEviction", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSimulateEviction prepares the SimulateEviction request.
func (c VirtualMachinesClient) preparerForSimulateEviction(ctx context.Context, id VirtualMachineId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/simulateEviction", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForSimulateEviction handles the response to the SimulateEviction request. The method always
// closes the http.Response Body.
func (c VirtualMachinesClient) responderForSimulateEviction(resp *http.Response) (result SimulateEvictionOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
