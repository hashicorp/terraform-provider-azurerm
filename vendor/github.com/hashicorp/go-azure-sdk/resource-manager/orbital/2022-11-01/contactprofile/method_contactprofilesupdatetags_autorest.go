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

type ContactProfilesUpdateTagsOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ContactProfilesUpdateTags ...
func (c ContactProfileClient) ContactProfilesUpdateTags(ctx context.Context, id ContactProfileId, input TagsObject) (result ContactProfilesUpdateTagsOperationResponse, err error) {
	req, err := c.preparerForContactProfilesUpdateTags(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "contactprofile.ContactProfileClient", "ContactProfilesUpdateTags", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForContactProfilesUpdateTags(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "contactprofile.ContactProfileClient", "ContactProfilesUpdateTags", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ContactProfilesUpdateTagsThenPoll performs ContactProfilesUpdateTags then polls until it's completed
func (c ContactProfileClient) ContactProfilesUpdateTagsThenPoll(ctx context.Context, id ContactProfileId, input TagsObject) error {
	result, err := c.ContactProfilesUpdateTags(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ContactProfilesUpdateTags: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ContactProfilesUpdateTags: %+v", err)
	}

	return nil
}

// preparerForContactProfilesUpdateTags prepares the ContactProfilesUpdateTags request.
func (c ContactProfileClient) preparerForContactProfilesUpdateTags(ctx context.Context, id ContactProfileId, input TagsObject) (*http.Request, error) {
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

// senderForContactProfilesUpdateTags sends the ContactProfilesUpdateTags request. The method will close the
// http.Response Body if it receives an error.
func (c ContactProfileClient) senderForContactProfilesUpdateTags(ctx context.Context, req *http.Request) (future ContactProfilesUpdateTagsOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
