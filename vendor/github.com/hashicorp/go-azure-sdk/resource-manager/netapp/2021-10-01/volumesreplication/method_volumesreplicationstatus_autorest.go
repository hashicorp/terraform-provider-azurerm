package volumesreplication

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VolumesReplicationStatusOperationResponse struct {
	HttpResponse *http.Response
	Model        *ReplicationStatus
}

// VolumesReplicationStatus ...
func (c VolumesReplicationClient) VolumesReplicationStatus(ctx context.Context, id VolumeId) (result VolumesReplicationStatusOperationResponse, err error) {
	req, err := c.preparerForVolumesReplicationStatus(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "volumesreplication.VolumesReplicationClient", "VolumesReplicationStatus", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "volumesreplication.VolumesReplicationClient", "VolumesReplicationStatus", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForVolumesReplicationStatus(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "volumesreplication.VolumesReplicationClient", "VolumesReplicationStatus", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForVolumesReplicationStatus prepares the VolumesReplicationStatus request.
func (c VolumesReplicationClient) preparerForVolumesReplicationStatus(ctx context.Context, id VolumeId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/replicationStatus", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForVolumesReplicationStatus handles the response to the VolumesReplicationStatus request. The method always
// closes the http.Response Body.
func (c VolumesReplicationClient) responderForVolumesReplicationStatus(resp *http.Response) (result VolumesReplicationStatusOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
