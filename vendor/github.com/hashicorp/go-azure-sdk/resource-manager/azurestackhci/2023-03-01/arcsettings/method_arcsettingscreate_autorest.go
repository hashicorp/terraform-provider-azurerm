package arcsettings

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ArcSettingsCreateOperationResponse struct {
	HttpResponse *http.Response
	Model        *ArcSetting
}

// ArcSettingsCreate ...
func (c ArcSettingsClient) ArcSettingsCreate(ctx context.Context, id ArcSettingId, input ArcSetting) (result ArcSettingsCreateOperationResponse, err error) {
	req, err := c.preparerForArcSettingsCreate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "arcsettings.ArcSettingsClient", "ArcSettingsCreate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "arcsettings.ArcSettingsClient", "ArcSettingsCreate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForArcSettingsCreate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "arcsettings.ArcSettingsClient", "ArcSettingsCreate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForArcSettingsCreate prepares the ArcSettingsCreate request.
func (c ArcSettingsClient) preparerForArcSettingsCreate(ctx context.Context, id ArcSettingId, input ArcSetting) (*http.Request, error) {
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

// responderForArcSettingsCreate handles the response to the ArcSettingsCreate request. The method always
// closes the http.Response Body.
func (c ArcSettingsClient) responderForArcSettingsCreate(resp *http.Response) (result ArcSettingsCreateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
