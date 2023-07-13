// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azuresdkhacks

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

func DeleteFirewall(ctx context.Context, client *network.AzureFirewallsClient, resourceGroupName string, azureFirewallName string) (result network.AzureFirewallsDeleteFuture, err error) {
	req, err := firewallDeletePreparer(ctx, client, resourceGroupName, azureFirewallName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.AzureFirewallsClient", "Delete", nil, "Failure preparing request")
		return
	}

	result, err = client.DeleteSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.AzureFirewallsClient", "Delete", result.Response(), "Failure sending request")
		return
	}

	return
}

func firewallDeletePreparer(ctx context.Context, client *network.AzureFirewallsClient, resourceGroupName string, azureFirewallName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"azureFirewallName": autorest.Encode("path", azureFirewallName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2021-05-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		// The patch! See https://github.com/Azure/azure-sdk-for-go/issues/17013 for details
		autorest.AsContentType("application/json"),

		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/azureFirewalls/{azureFirewallName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}
