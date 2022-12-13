package volumegroups

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VolumeGroupsListByNetAppAccountOperationResponse struct {
	HttpResponse *http.Response
	Model        *VolumeGroupList
}

// VolumeGroupsListByNetAppAccount ...
func (c VolumeGroupsClient) VolumeGroupsListByNetAppAccount(ctx context.Context, id NetAppAccountId) (result VolumeGroupsListByNetAppAccountOperationResponse, err error) {
	req, err := c.preparerForVolumeGroupsListByNetAppAccount(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "volumegroups.VolumeGroupsClient", "VolumeGroupsListByNetAppAccount", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "volumegroups.VolumeGroupsClient", "VolumeGroupsListByNetAppAccount", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForVolumeGroupsListByNetAppAccount(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "volumegroups.VolumeGroupsClient", "VolumeGroupsListByNetAppAccount", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForVolumeGroupsListByNetAppAccount prepares the VolumeGroupsListByNetAppAccount request.
func (c VolumeGroupsClient) preparerForVolumeGroupsListByNetAppAccount(ctx context.Context, id NetAppAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/volumeGroups", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForVolumeGroupsListByNetAppAccount handles the response to the VolumeGroupsListByNetAppAccount request. The method always
// closes the http.Response Body.
func (c VolumeGroupsClient) responderForVolumeGroupsListByNetAppAccount(resp *http.Response) (result VolumeGroupsListByNetAppAccountOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
