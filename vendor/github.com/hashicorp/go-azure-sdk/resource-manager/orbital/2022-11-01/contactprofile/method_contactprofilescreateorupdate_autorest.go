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

type ContactProfilesCreateOrUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ContactProfilesCreateOrUpdate ...
func (c ContactProfileClient) ContactProfilesCreateOrUpdate(ctx context.Context, id ContactProfileId, input ContactProfile) (result ContactProfilesCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForContactProfilesCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "contactprofile.ContactProfileClient", "ContactProfilesCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForContactProfilesCreateOrUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "contactprofile.ContactProfileClient", "ContactProfilesCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ContactProfilesCreateOrUpdateThenPoll performs ContactProfilesCreateOrUpdate then polls until it's completed
func (c ContactProfileClient) ContactProfilesCreateOrUpdateThenPoll(ctx context.Context, id ContactProfileId, input ContactProfile) error {
	result, err := c.ContactProfilesCreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ContactProfilesCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ContactProfilesCreateOrUpdate: %+v", err)
	}

	return nil
}

// preparerForContactProfilesCreateOrUpdate prepares the ContactProfilesCreateOrUpdate request.
func (c ContactProfileClient) preparerForContactProfilesCreateOrUpdate(ctx context.Context, id ContactProfileId, input ContactProfile) (*http.Request, error) {
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

// senderForContactProfilesCreateOrUpdate sends the ContactProfilesCreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c ContactProfileClient) senderForContactProfilesCreateOrUpdate(ctx context.Context, req *http.Request) (future ContactProfilesCreateOrUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
