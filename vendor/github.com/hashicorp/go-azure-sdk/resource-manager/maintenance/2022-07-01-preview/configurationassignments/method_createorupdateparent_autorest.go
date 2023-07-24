package configurationassignments

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CreateOrUpdateParentOperationResponse struct {
	HttpResponse *http.Response
	Model        *ConfigurationAssignment
}

// CreateOrUpdateParent ...
func (c ConfigurationAssignmentsClient) CreateOrUpdateParent(ctx context.Context, id ScopedConfigurationAssignmentId, input ConfigurationAssignment) (result CreateOrUpdateParentOperationResponse, err error) {
	req, err := c.preparerForCreateOrUpdateParent(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurationassignments.ConfigurationAssignmentsClient", "CreateOrUpdateParent", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurationassignments.ConfigurationAssignmentsClient", "CreateOrUpdateParent", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCreateOrUpdateParent(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurationassignments.ConfigurationAssignmentsClient", "CreateOrUpdateParent", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCreateOrUpdateParent prepares the CreateOrUpdateParent request.
func (c ConfigurationAssignmentsClient) preparerForCreateOrUpdateParent(ctx context.Context, id ScopedConfigurationAssignmentId, input ConfigurationAssignment) (*http.Request, error) {
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

// responderForCreateOrUpdateParent handles the response to the CreateOrUpdateParent request. The method always
// closes the http.Response Body.
func (c ConfigurationAssignmentsClient) responderForCreateOrUpdateParent(resp *http.Response) (result CreateOrUpdateParentOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
