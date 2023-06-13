package contactprofile

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

type ContactProfilesDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ContactProfilesDelete ...
func (c ContactProfileClient) ContactProfilesDelete(ctx context.Context, id ContactProfileId) (result ContactProfilesDeleteOperationResponse, err error) {
	req, err := c.preparerForContactProfilesDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "contactprofile.ContactProfileClient", "ContactProfilesDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForContactProfilesDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "contactprofile.ContactProfileClient", "ContactProfilesDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ContactProfilesDeleteThenPoll performs ContactProfilesDelete then polls until it's completed
func (c ContactProfileClient) ContactProfilesDeleteThenPoll(ctx context.Context, id ContactProfileId) error {
	result, err := c.ContactProfilesDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ContactProfilesDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ContactProfilesDelete: %+v", err)
	}

	return nil
}

// preparerForContactProfilesDelete prepares the ContactProfilesDelete request.
func (c ContactProfileClient) preparerForContactProfilesDelete(ctx context.Context, id ContactProfileId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForContactProfilesDelete sends the ContactProfilesDelete request. The method will close the
// http.Response Body if it receives an error.
func (c ContactProfileClient) senderForContactProfilesDelete(ctx context.Context, req *http.Request) (future ContactProfilesDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
