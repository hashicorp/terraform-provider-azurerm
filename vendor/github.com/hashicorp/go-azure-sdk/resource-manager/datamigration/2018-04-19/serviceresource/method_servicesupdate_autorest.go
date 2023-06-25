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

type ServicesUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ServicesUpdate ...
func (c ServiceResourceClient) ServicesUpdate(ctx context.Context, id ServiceId, input DataMigrationService) (result ServicesUpdateOperationResponse, err error) {
	req, err := c.preparerForServicesUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "serviceresource.ServiceResourceClient", "ServicesUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForServicesUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "serviceresource.ServiceResourceClient", "ServicesUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ServicesUpdateThenPoll performs ServicesUpdate then polls until it's completed
func (c ServiceResourceClient) ServicesUpdateThenPoll(ctx context.Context, id ServiceId, input DataMigrationService) error {
	result, err := c.ServicesUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ServicesUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ServicesUpdate: %+v", err)
	}

	return nil
}

// preparerForServicesUpdate prepares the ServicesUpdate request.
func (c ServiceResourceClient) preparerForServicesUpdate(ctx context.Context, id ServiceId, input DataMigrationService) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForServicesUpdate sends the ServicesUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c ServiceResourceClient) senderForServicesUpdate(ctx context.Context, req *http.Request) (future ServicesUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
