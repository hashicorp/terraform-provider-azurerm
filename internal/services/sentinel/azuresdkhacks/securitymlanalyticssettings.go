// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azuresdkhacks

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	securityinsight "github.com/tombuildsstuff/kermit/sdk/securityinsights/2022-10-01-preview/securityinsights"
)

type SecurityMLAnalyticsSettingsClient struct {
	securityinsight.BaseClient
}

func (client SecurityMLAnalyticsSettingsClient) List(ctx context.Context, resourceGroupName string, workspaceName string) (result SecurityMLAnalyticsSettingsListPage, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}}},
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil}}},
		{TargetValue: workspaceName,
			Constraints: []validation.Constraint{{Target: "workspaceName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "workspaceName", Name: validation.MinLength, Rule: 1, Chain: nil}}}}); err != nil {
		return result, validation.NewError("securityinsight.SecurityMLAnalyticsSettingsClient", "List", err.Error())
	}

	result.fn = client.listNextResults
	req, err := client.ListPreparer(ctx, resourceGroupName, workspaceName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "securityinsight.SecurityMLAnalyticsSettingsClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.smasl.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "securityinsight.SecurityMLAnalyticsSettingsClient", "List", resp, "Failure sending request")
		return
	}

	result.smasl, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "securityinsight.SecurityMLAnalyticsSettingsClient", "List", resp, "Failure responding to request")
		return
	}
	if result.smasl.hasNextLink() && result.smasl.IsEmpty() {
		err = result.NextWithContext(ctx)
		return
	}

	return
}

// ListPreparer prepares the List request.
func (client SecurityMLAnalyticsSettingsClient) ListPreparer(ctx context.Context, resourceGroupName string, workspaceName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"workspaceName":     autorest.Encode("path", workspaceName),
	}

	const APIVersion = "2022-10-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.OperationalInsights/workspaces/{workspaceName}/providers/Microsoft.SecurityInsights/securityMLAnalyticsSettings", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client SecurityMLAnalyticsSettingsClient) ListSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client SecurityMLAnalyticsSettingsClient) ListResponder(resp *http.Response) (result SecurityMLAnalyticsSettingsList, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listNextResults retrieves the next set of results, if any.
func (client SecurityMLAnalyticsSettingsClient) listNextResults(ctx context.Context, lastResults SecurityMLAnalyticsSettingsList) (result SecurityMLAnalyticsSettingsList, err error) {
	req, err := lastResults.securityMLAnalyticsSettingsListPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "securityinsight.SecurityMLAnalyticsSettingsClient", "listNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "securityinsight.SecurityMLAnalyticsSettingsClient", "listNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "securityinsight.SecurityMLAnalyticsSettingsClient", "listNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListComplete enumerates all values, automatically crossing page boundaries as required.
func (client SecurityMLAnalyticsSettingsClient) ListComplete(ctx context.Context, resourceGroupName string, workspaceName string) (result SecurityMLAnalyticsSettingsListIterator, err error) {
	result.page, err = client.List(ctx, resourceGroupName, workspaceName)
	return
}
