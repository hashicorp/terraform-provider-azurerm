package agentpools

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

type GetAvailableAgentPoolVersionsOperationResponse struct {
	HttpResponse *http.Response
	Model        *AgentPoolAvailableVersions
}

// GetAvailableAgentPoolVersions ...
func (c AgentPoolsClient) GetAvailableAgentPoolVersions(ctx context.Context, id commonids.KubernetesClusterId) (result GetAvailableAgentPoolVersionsOperationResponse, err error) {
	req, err := c.preparerForGetAvailableAgentPoolVersions(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "agentpools.AgentPoolsClient", "GetAvailableAgentPoolVersions", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "agentpools.AgentPoolsClient", "GetAvailableAgentPoolVersions", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetAvailableAgentPoolVersions(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "agentpools.AgentPoolsClient", "GetAvailableAgentPoolVersions", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetAvailableAgentPoolVersions prepares the GetAvailableAgentPoolVersions request.
func (c AgentPoolsClient) preparerForGetAvailableAgentPoolVersions(ctx context.Context, id commonids.KubernetesClusterId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/availableAgentPoolVersions", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGetAvailableAgentPoolVersions handles the response to the GetAvailableAgentPoolVersions request. The method always
// closes the http.Response Body.
func (c AgentPoolsClient) responderForGetAvailableAgentPoolVersions(resp *http.Response) (result GetAvailableAgentPoolVersionsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
