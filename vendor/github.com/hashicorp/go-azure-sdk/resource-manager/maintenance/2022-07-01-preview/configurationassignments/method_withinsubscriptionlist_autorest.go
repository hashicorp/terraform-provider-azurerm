package configurationassignments

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WithinSubscriptionListOperationResponse struct {
	HttpResponse *http.Response
	Model        *ListConfigurationAssignmentsResult
}

// WithinSubscriptionList ...
func (c ConfigurationAssignmentsClient) WithinSubscriptionList(ctx context.Context, id commonids.SubscriptionId) (result WithinSubscriptionListOperationResponse, err error) {
	req, err := c.preparerForWithinSubscriptionList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurationassignments.ConfigurationAssignmentsClient", "WithinSubscriptionList", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurationassignments.ConfigurationAssignmentsClient", "WithinSubscriptionList", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForWithinSubscriptionList(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurationassignments.ConfigurationAssignmentsClient", "WithinSubscriptionList", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForWithinSubscriptionList prepares the WithinSubscriptionList request.
func (c ConfigurationAssignmentsClient) preparerForWithinSubscriptionList(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Maintenance/configurationAssignments", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForWithinSubscriptionList handles the response to the WithinSubscriptionList request. The method always
// closes the http.Response Body.
func (c ConfigurationAssignmentsClient) responderForWithinSubscriptionList(resp *http.Response) (result WithinSubscriptionListOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
