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

type InstanceViewOperationResponse struct {
	HttpResponse *http.Response
	Model        *VirtualMachineInstanceView
}

// InstanceView ...
func (c VirtualMachinesClient) InstanceView(ctx context.Context, id VirtualMachineId) (result InstanceViewOperationResponse, err error) {
	req, err := c.preparerForInstanceView(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "InstanceView", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "InstanceView", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForInstanceView(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "InstanceView", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForInstanceView prepares the InstanceView request.
func (c VirtualMachinesClient) preparerForInstanceView(ctx context.Context, id VirtualMachineId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/instanceView", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForInstanceView handles the response to the InstanceView request. The method always
// closes the http.Response Body.
func (c VirtualMachinesClient) responderForInstanceView(resp *http.Response) (result InstanceViewOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
