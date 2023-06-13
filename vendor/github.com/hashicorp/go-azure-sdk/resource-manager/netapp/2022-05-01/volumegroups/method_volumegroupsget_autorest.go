package volumegroups

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VolumeGroupsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *VolumeGroupDetails
}

// VolumeGroupsGet ...
func (c VolumeGroupsClient) VolumeGroupsGet(ctx context.Context, id VolumeGroupId) (result VolumeGroupsGetOperationResponse, err error) {
	req, err := c.preparerForVolumeGroupsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "volumegroups.VolumeGroupsClient", "VolumeGroupsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "volumegroups.VolumeGroupsClient", "VolumeGroupsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForVolumeGroupsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "volumegroups.VolumeGroupsClient", "VolumeGroupsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForVolumeGroupsGet prepares the VolumeGroupsGet request.
func (c VolumeGroupsClient) preparerForVolumeGroupsGet(ctx context.Context, id VolumeGroupId) (*http.Request, error) {
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

// responderForVolumeGroupsGet handles the response to the VolumeGroupsGet request. The method always
// closes the http.Response Body.
func (c VolumeGroupsClient) responderForVolumeGroupsGet(resp *http.Response) (result VolumeGroupsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
