package arcsettings

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ArcSettingsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *ArcSetting
}

// ArcSettingsGet ...
func (c ArcSettingsClient) ArcSettingsGet(ctx context.Context, id ArcSettingId) (result ArcSettingsGetOperationResponse, err error) {
	req, err := c.preparerForArcSettingsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "arcsettings.ArcSettingsClient", "ArcSettingsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "arcsettings.ArcSettingsClient", "ArcSettingsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForArcSettingsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "arcsettings.ArcSettingsClient", "ArcSettingsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForArcSettingsGet prepares the ArcSettingsGet request.
func (c ArcSettingsClient) preparerForArcSettingsGet(ctx context.Context, id ArcSettingId) (*http.Request, error) {
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

// responderForArcSettingsGet handles the response to the ArcSettingsGet request. The method always
// closes the http.Response Body.
func (c ArcSettingsClient) responderForArcSettingsGet(resp *http.Response) (result ArcSettingsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
