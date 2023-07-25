package configurationassignments

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeleteParentOperationResponse struct {
	HttpResponse *http.Response
	Model        *ConfigurationAssignment
}

// DeleteParent ...
func (c ConfigurationAssignmentsClient) DeleteParent(ctx context.Context, id ScopedConfigurationAssignmentId) (result DeleteParentOperationResponse, err error) {
	req, err := c.preparerForDeleteParent(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurationassignments.ConfigurationAssignmentsClient", "DeleteParent", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurationassignments.ConfigurationAssignmentsClient", "DeleteParent", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDeleteParent(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurationassignments.ConfigurationAssignmentsClient", "DeleteParent", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDeleteParent prepares the DeleteParent request.
func (c ConfigurationAssignmentsClient) preparerForDeleteParent(ctx context.Context, id ScopedConfigurationAssignmentId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDeleteParent handles the response to the DeleteParent request. The method always
// closes the http.Response Body.
func (c ConfigurationAssignmentsClient) responderForDeleteParent(resp *http.Response) (result DeleteParentOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
