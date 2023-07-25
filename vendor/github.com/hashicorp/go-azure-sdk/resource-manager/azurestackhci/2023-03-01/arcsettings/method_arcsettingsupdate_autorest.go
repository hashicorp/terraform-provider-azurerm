package arcsettings

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ArcSettingsUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *ArcSetting
}

// ArcSettingsUpdate ...
func (c ArcSettingsClient) ArcSettingsUpdate(ctx context.Context, id ArcSettingId, input ArcSettingsPatch) (result ArcSettingsUpdateOperationResponse, err error) {
	req, err := c.preparerForArcSettingsUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "arcsettings.ArcSettingsClient", "ArcSettingsUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "arcsettings.ArcSettingsClient", "ArcSettingsUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForArcSettingsUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "arcsettings.ArcSettingsClient", "ArcSettingsUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForArcSettingsUpdate prepares the ArcSettingsUpdate request.
func (c ArcSettingsClient) preparerForArcSettingsUpdate(ctx context.Context, id ArcSettingId, input ArcSettingsPatch) (*http.Request, error) {
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

// responderForArcSettingsUpdate handles the response to the ArcSettingsUpdate request. The method always
// closes the http.Response Body.
func (c ArcSettingsClient) responderForArcSettingsUpdate(resp *http.Response) (result ArcSettingsUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
