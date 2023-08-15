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

type ListClusterAdminCredentialsOperationResponse struct {
	HttpResponse *http.Response
	Model        *CredentialResults
}

type ListClusterAdminCredentialsOperationOptions struct {
	ServerFqdn *string
}

func DefaultListClusterAdminCredentialsOperationOptions() ListClusterAdminCredentialsOperationOptions {
	return ListClusterAdminCredentialsOperationOptions{}
}

func (o ListClusterAdminCredentialsOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ListClusterAdminCredentialsOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.ServerFqdn != nil {
		out["server-fqdn"] = *o.ServerFqdn
	}

	return out
}

// ListClusterAdminCredentials ...
func (c ManagedClustersClient) ListClusterAdminCredentials(ctx context.Context, id commonids.KubernetesClusterId, options ListClusterAdminCredentialsOperationOptions) (result ListClusterAdminCredentialsOperationResponse, err error) {
	req, err := c.preparerForListClusterAdminCredentials(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "ListClusterAdminCredentials", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "ListClusterAdminCredentials", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListClusterAdminCredentials(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "ListClusterAdminCredentials", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListClusterAdminCredentials prepares the ListClusterAdminCredentials request.
func (c ManagedClustersClient) preparerForListClusterAdminCredentials(ctx context.Context, id commonids.KubernetesClusterId, options ListClusterAdminCredentialsOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/listClusterAdminCredential", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListClusterAdminCredentials handles the response to the ListClusterAdminCredentials request. The method always
// closes the http.Response Body.
func (c ManagedClustersClient) responderForListClusterAdminCredentials(resp *http.Response) (result ListClusterAdminCredentialsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
