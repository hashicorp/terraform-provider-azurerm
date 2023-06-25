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

type ListApplicableSchedulesOperationResponse struct {
	HttpResponse *http.Response
	Model        *ApplicableSchedule
}

// ListApplicableSchedules ...
func (c VirtualMachinesClient) ListApplicableSchedules(ctx context.Context, id VirtualMachineId) (result ListApplicableSchedulesOperationResponse, err error) {
	req, err := c.preparerForListApplicableSchedules(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "ListApplicableSchedules", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "ListApplicableSchedules", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListApplicableSchedules(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "ListApplicableSchedules", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListApplicableSchedules prepares the ListApplicableSchedules request.
func (c VirtualMachinesClient) preparerForListApplicableSchedules(ctx context.Context, id VirtualMachineId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listApplicableSchedules", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListApplicableSchedules handles the response to the ListApplicableSchedules request. The method always
// closes the http.Response Body.
func (c VirtualMachinesClient) responderForListApplicableSchedules(resp *http.Response) (result ListApplicableSchedulesOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
