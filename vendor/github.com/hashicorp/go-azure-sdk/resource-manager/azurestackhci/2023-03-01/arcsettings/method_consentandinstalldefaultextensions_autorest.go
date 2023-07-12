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

type ConsentAndInstallDefaultExtensionsOperationResponse struct {
	HttpResponse *http.Response
	Model        *ArcSetting
}

// ConsentAndInstallDefaultExtensions ...
func (c ArcSettingsClient) ConsentAndInstallDefaultExtensions(ctx context.Context, id ArcSettingId) (result ConsentAndInstallDefaultExtensionsOperationResponse, err error) {
	req, err := c.preparerForConsentAndInstallDefaultExtensions(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "arcsettings.ArcSettingsClient", "ConsentAndInstallDefaultExtensions", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "arcsettings.ArcSettingsClient", "ConsentAndInstallDefaultExtensions", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForConsentAndInstallDefaultExtensions(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "arcsettings.ArcSettingsClient", "ConsentAndInstallDefaultExtensions", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForConsentAndInstallDefaultExtensions prepares the ConsentAndInstallDefaultExtensions request.
func (c ArcSettingsClient) preparerForConsentAndInstallDefaultExtensions(ctx context.Context, id ArcSettingId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/consentAndInstallDefaultExtensions", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForConsentAndInstallDefaultExtensions handles the response to the ConsentAndInstallDefaultExtensions request. The method always
// closes the http.Response Body.
func (c ArcSettingsClient) responderForConsentAndInstallDefaultExtensions(resp *http.Response) (result ConsentAndInstallDefaultExtensionsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
