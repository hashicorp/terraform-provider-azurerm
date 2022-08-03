package contactprofile

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContactProfilesGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *ContactProfile
}

// ContactProfilesGet ...
func (c ContactProfileClient) ContactProfilesGet(ctx context.Context, id ContactProfileId) (result ContactProfilesGetOperationResponse, err error) {
	req, err := c.preparerForContactProfilesGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "contactprofile.ContactProfileClient", "ContactProfilesGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "contactprofile.ContactProfileClient", "ContactProfilesGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForContactProfilesGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "contactprofile.ContactProfileClient", "ContactProfilesGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForContactProfilesGet prepares the ContactProfilesGet request.
func (c ContactProfileClient) preparerForContactProfilesGet(ctx context.Context, id ContactProfileId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForContactProfilesGet handles the response to the ContactProfilesGet request. The method always
// closes the http.Response Body.
func (c ContactProfileClient) responderForContactProfilesGet(resp *http.Response) (result ContactProfilesGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
