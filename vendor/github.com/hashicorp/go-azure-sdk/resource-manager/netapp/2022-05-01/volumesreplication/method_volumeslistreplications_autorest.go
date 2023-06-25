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

type VolumesListReplicationsOperationResponse struct {
	HttpResponse *http.Response
	Model        *ListReplications
}

// VolumesListReplications ...
func (c VolumesReplicationClient) VolumesListReplications(ctx context.Context, id VolumeId) (result VolumesListReplicationsOperationResponse, err error) {
	req, err := c.preparerForVolumesListReplications(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "volumesreplication.VolumesReplicationClient", "VolumesListReplications", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "volumesreplication.VolumesReplicationClient", "VolumesListReplications", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForVolumesListReplications(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "volumesreplication.VolumesReplicationClient", "VolumesListReplications", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForVolumesListReplications prepares the VolumesListReplications request.
func (c VolumesReplicationClient) preparerForVolumesListReplications(ctx context.Context, id VolumeId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listReplications", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForVolumesListReplications handles the response to the VolumesListReplications request. The method always
// closes the http.Response Body.
func (c VolumesReplicationClient) responderForVolumesListReplications(resp *http.Response) (result VolumesListReplicationsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
