package managedclusters

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

type ListClusterMonitoringUserCredentialsOperationResponse struct {
	HttpResponse *http.Response
	Model        *CredentialResults
}

type ListClusterMonitoringUserCredentialsOperationOptions struct {
	ServerFqdn *string
}

func DefaultListClusterMonitoringUserCredentialsOperationOptions() ListClusterMonitoringUserCredentialsOperationOptions {
	return ListClusterMonitoringUserCredentialsOperationOptions{}
}

func (o ListClusterMonitoringUserCredentialsOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ListClusterMonitoringUserCredentialsOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.ServerFqdn != nil {
		out["server-fqdn"] = *o.ServerFqdn
	}

	return out
}

// ListClusterMonitoringUserCredentials ...
func (c ManagedClustersClient) ListClusterMonitoringUserCredentials(ctx context.Context, id commonids.KubernetesClusterId, options ListClusterMonitoringUserCredentialsOperationOptions) (result ListClusterMonitoringUserCredentialsOperationResponse, err error) {
	req, err := c.preparerForListClusterMonitoringUserCredentials(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "ListClusterMonitoringUserCredentials", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "ListClusterMonitoringUserCredentials", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListClusterMonitoringUserCredentials(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "ListClusterMonitoringUserCredentials", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListClusterMonitoringUserCredentials prepares the ListClusterMonitoringUserCredentials request.
func (c ManagedClustersClient) preparerForListClusterMonitoringUserCredentials(ctx context.Context, id commonids.KubernetesClusterId, options ListClusterMonitoringUserCredentialsOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/listClusterMonitoringUserCredential", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListClusterMonitoringUserCredentials handles the response to the ListClusterMonitoringUserCredentials request. The method always
// closes the http.Response Body.
func (c ManagedClustersClient) responderForListClusterMonitoringUserCredentials(resp *http.Response) (result ListClusterMonitoringUserCredentialsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
