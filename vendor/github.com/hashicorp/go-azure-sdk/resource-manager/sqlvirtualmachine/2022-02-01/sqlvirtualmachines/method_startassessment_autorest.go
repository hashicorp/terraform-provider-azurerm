package sqlvirtualmachines

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StartAssessmentOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// StartAssessment ...
func (c SqlVirtualMachinesClient) StartAssessment(ctx context.Context, id SqlVirtualMachineId) (result StartAssessmentOperationResponse, err error) {
	req, err := c.preparerForStartAssessment(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "sqlvirtualmachines.SqlVirtualMachinesClient", "StartAssessment", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForStartAssessment(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "sqlvirtualmachines.SqlVirtualMachinesClient", "StartAssessment", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// StartAssessmentThenPoll performs StartAssessment then polls until it's completed
func (c SqlVirtualMachinesClient) StartAssessmentThenPoll(ctx context.Context, id SqlVirtualMachineId) error {
	result, err := c.StartAssessment(ctx, id)
	if err != nil {
		return fmt.Errorf("performing StartAssessment: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after StartAssessment: %+v", err)
	}

	return nil
}

// preparerForStartAssessment prepares the StartAssessment request.
func (c SqlVirtualMachinesClient) preparerForStartAssessment(ctx context.Context, id SqlVirtualMachineId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/startAssessment", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForStartAssessment sends the StartAssessment request. The method will close the
// http.Response Body if it receives an error.
func (c SqlVirtualMachinesClient) senderForStartAssessment(ctx context.Context, req *http.Request) (future StartAssessmentOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
