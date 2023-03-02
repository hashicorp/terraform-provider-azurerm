package serviceresource

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

type ServicesCreateOrUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ServicesCreateOrUpdate ...
func (c ServiceResourceClient) ServicesCreateOrUpdate(ctx context.Context, id ServiceId, input DataMigrationService) (result ServicesCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForServicesCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "serviceresource.ServiceResourceClient", "ServicesCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForServicesCreateOrUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "serviceresource.ServiceResourceClient", "ServicesCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ServicesCreateOrUpdateThenPoll performs ServicesCreateOrUpdate then polls until it's completed
func (c ServiceResourceClient) ServicesCreateOrUpdateThenPoll(ctx context.Context, id ServiceId, input DataMigrationService) error {
	result, err := c.ServicesCreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ServicesCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ServicesCreateOrUpdate: %+v", err)
	}

	return nil
}

// preparerForServicesCreateOrUpdate prepares the ServicesCreateOrUpdate request.
func (c ServiceResourceClient) preparerForServicesCreateOrUpdate(ctx context.Context, id ServiceId, input DataMigrationService) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForServicesCreateOrUpdate sends the ServicesCreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c ServiceResourceClient) senderForServicesCreateOrUpdate(ctx context.Context, req *http.Request) (future ServicesCreateOrUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
