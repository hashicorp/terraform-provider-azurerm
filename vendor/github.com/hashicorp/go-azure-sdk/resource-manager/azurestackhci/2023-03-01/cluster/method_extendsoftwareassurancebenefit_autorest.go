package cluster

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

type ExtendSoftwareAssuranceBenefitOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ExtendSoftwareAssuranceBenefit ...
func (c ClusterClient) ExtendSoftwareAssuranceBenefit(ctx context.Context, id ClusterId, input SoftwareAssuranceChangeRequest) (result ExtendSoftwareAssuranceBenefitOperationResponse, err error) {
	req, err := c.preparerForExtendSoftwareAssuranceBenefit(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cluster.ClusterClient", "ExtendSoftwareAssuranceBenefit", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForExtendSoftwareAssuranceBenefit(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cluster.ClusterClient", "ExtendSoftwareAssuranceBenefit", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ExtendSoftwareAssuranceBenefitThenPoll performs ExtendSoftwareAssuranceBenefit then polls until it's completed
func (c ClusterClient) ExtendSoftwareAssuranceBenefitThenPoll(ctx context.Context, id ClusterId, input SoftwareAssuranceChangeRequest) error {
	result, err := c.ExtendSoftwareAssuranceBenefit(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ExtendSoftwareAssuranceBenefit: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ExtendSoftwareAssuranceBenefit: %+v", err)
	}

	return nil
}

// preparerForExtendSoftwareAssuranceBenefit prepares the ExtendSoftwareAssuranceBenefit request.
func (c ClusterClient) preparerForExtendSoftwareAssuranceBenefit(ctx context.Context, id ClusterId, input SoftwareAssuranceChangeRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/extendSoftwareAssuranceBenefit", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForExtendSoftwareAssuranceBenefit sends the ExtendSoftwareAssuranceBenefit request. The method will close the
// http.Response Body if it receives an error.
func (c ClusterClient) senderForExtendSoftwareAssuranceBenefit(ctx context.Context, req *http.Request) (future ExtendSoftwareAssuranceBenefitOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
