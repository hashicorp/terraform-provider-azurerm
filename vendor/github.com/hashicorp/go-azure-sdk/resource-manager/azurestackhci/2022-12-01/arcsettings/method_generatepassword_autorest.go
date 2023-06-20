package arcsettings

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GeneratePasswordOperationResponse struct {
	HttpResponse *http.Response
	Model        *PasswordCredential
}

// GeneratePassword ...
func (c ArcSettingsClient) GeneratePassword(ctx context.Context, id ArcSettingId) (result GeneratePasswordOperationResponse, err error) {
	req, err := c.preparerForGeneratePassword(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "arcsettings.ArcSettingsClient", "GeneratePassword", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "arcsettings.ArcSettingsClient", "GeneratePassword", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGeneratePassword(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "arcsettings.ArcSettingsClient", "GeneratePassword", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGeneratePassword prepares the GeneratePassword request.
func (c ArcSettingsClient) preparerForGeneratePassword(ctx context.Context, id ArcSettingId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/generatePassword", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGeneratePassword handles the response to the GeneratePassword request. The method always
// closes the http.Response Body.
func (c ArcSettingsClient) responderForGeneratePassword(resp *http.Response) (result GeneratePasswordOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
