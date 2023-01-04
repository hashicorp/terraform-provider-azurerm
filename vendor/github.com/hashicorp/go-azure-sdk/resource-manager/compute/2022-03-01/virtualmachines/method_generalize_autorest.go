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

type GeneralizeOperationResponse struct {
	HttpResponse *http.Response
}

// Generalize ...
func (c VirtualMachinesClient) Generalize(ctx context.Context, id VirtualMachineId) (result GeneralizeOperationResponse, err error) {
	req, err := c.preparerForGeneralize(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "Generalize", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "Generalize", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGeneralize(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "Generalize", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGeneralize prepares the Generalize request.
func (c VirtualMachinesClient) preparerForGeneralize(ctx context.Context, id VirtualMachineId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/generalize", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGeneralize handles the response to the Generalize request. The method always
// closes the http.Response Body.
func (c VirtualMachinesClient) responderForGeneralize(resp *http.Response) (result GeneralizeOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
