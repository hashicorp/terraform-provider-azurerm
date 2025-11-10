// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azuresdkhacks

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/Azure/go-autorest/autorest/validation"
	securityinsight "github.com/jackofallops/kermit/sdk/securityinsights/2022-10-01-preview/securityinsights"
)

type ThreatIntelligenceIndicatorClient struct {
	securityinsight.BaseClient
}

func (client ThreatIntelligenceIndicatorClient) Get(ctx context.Context, resourceGroupName string, workspaceName string, name string) (result ThreatIntelligenceInformationModel, err error) {
	if err := validation.Validate([]validation.Validation{
		{
			TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}},
		},
		{
			TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{
				{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
			},
		},
		{
			TargetValue: workspaceName,
			Constraints: []validation.Constraint{
				{Target: "workspaceName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "workspaceName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "workspaceName", Name: validation.Pattern, Rule: `^[A-Za-z0-9][A-Za-z0-9-]+[A-Za-z0-9]$`, Chain: nil},
			},
		},
	}); err != nil {
		return result, validation.NewError("securityinsight.ThreatIntelligenceIndicatorClient", "Get", "%+v", err.Error())
	}

	req, err := client.GetPreparer(ctx, resourceGroupName, workspaceName, name)
	if err != nil {
		err = autorest.NewErrorWithError(err, "securityinsight.ThreatIntelligenceIndicatorClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "securityinsight.ThreatIntelligenceIndicatorClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "securityinsight.ThreatIntelligenceIndicatorClient", "Get", resp, "Failure responding to request")
		return
	}

	return
}

func (client ThreatIntelligenceIndicatorClient) GetPreparer(ctx context.Context, resourceGroupName string, workspaceName string, name string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"name":              autorest.Encode("path", name),
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
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.OperationalInsights/workspaces/{workspaceName}/providers/Microsoft.SecurityInsights/threatIntelligence/main/indicators/{name}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

func (client ThreatIntelligenceIndicatorClient) GetSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

func (client ThreatIntelligenceIndicatorClient) GetResponder(resp *http.Response) (result ThreatIntelligenceInformationModel, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

func (client ThreatIntelligenceIndicatorClient) CreateIndicator(ctx context.Context, resourceGroupName string, workspaceName string, threatIntelligenceProperties ThreatIntelligenceIndicatorModel) (result ThreatIntelligenceInformationModel, err error) {
	if err := validation.Validate([]validation.Validation{
		{
			TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}},
		},
		{
			TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{
				{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
			},
		},
		{
			TargetValue: workspaceName,
			Constraints: []validation.Constraint{
				{Target: "workspaceName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "workspaceName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "workspaceName", Name: validation.Pattern, Rule: `^[A-Za-z0-9][A-Za-z0-9-]+[A-Za-z0-9]$`, Chain: nil},
			},
		},
	}); err != nil {
		return result, validation.NewError("securityinsight.ThreatIntelligenceIndicatorClient", "CreateIndicator", "%+v", err.Error())
	}

	req, err := client.CreateIndicatorPreparer(ctx, resourceGroupName, workspaceName, threatIntelligenceProperties)
	if err != nil {
		err = autorest.NewErrorWithError(err, "securityinsight.ThreatIntelligenceIndicatorClient", "CreateIndicator", nil, "Failure preparing request")
		return
	}

	resp, err := client.CreateIndicatorSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "securityinsight.ThreatIntelligenceIndicatorClient", "CreateIndicator", resp, "Failure sending request")
		return
	}

	result, err = client.CreateIndicatorResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "securityinsight.ThreatIntelligenceIndicatorClient", "CreateIndicator", resp, "Failure responding to request")
		return
	}

	return
}

// CreateIndicatorPreparer prepares the CreateIndicator request.
func (client ThreatIntelligenceIndicatorClient) CreateIndicatorPreparer(ctx context.Context, resourceGroupName string, workspaceName string, threatIntelligenceProperties ThreatIntelligenceIndicatorModel) (*http.Request, error) {
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
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.OperationalInsights/workspaces/{workspaceName}/providers/Microsoft.SecurityInsights/threatIntelligence/main/createIndicator", pathParameters),
		autorest.WithJSON(threatIntelligenceProperties),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

func (client ThreatIntelligenceIndicatorClient) CreateIndicatorSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

func (client ThreatIntelligenceIndicatorClient) CreateIndicatorResponder(resp *http.Response) (result ThreatIntelligenceInformationModel, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

func (client ThreatIntelligenceIndicatorClient) QueryIndicators(ctx context.Context, resourceGroupName string, workspaceName string, threatIntelligenceFilteringCriteria securityinsight.ThreatIntelligenceFilteringCriteria) (result ThreatIntelligenceInformationListPage, err error) {
	if err := validation.Validate([]validation.Validation{
		{
			TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}},
		},
		{
			TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{
				{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
			},
		},
		{
			TargetValue: workspaceName,
			Constraints: []validation.Constraint{
				{Target: "workspaceName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "workspaceName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "workspaceName", Name: validation.Pattern, Rule: `^[A-Za-z0-9][A-Za-z0-9-]+[A-Za-z0-9]$`, Chain: nil},
			},
		},
	}); err != nil {
		return result, validation.NewError("securityinsight.ThreatIntelligenceIndicatorClient", "QueryIndicators", "%+v", err.Error())
	}

	result.fn = client.queryIndicatorsNextResults
	req, err := client.QueryIndicatorsPreparer(ctx, resourceGroupName, workspaceName, threatIntelligenceFilteringCriteria)
	if err != nil {
		err = autorest.NewErrorWithError(err, "securityinsight.ThreatIntelligenceIndicatorClient", "QueryIndicators", nil, "Failure preparing request")
		return
	}

	resp, err := client.QueryIndicatorsSender(req)
	if err != nil {
		result.tiil.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "securityinsight.ThreatIntelligenceIndicatorClient", "QueryIndicators", resp, "Failure sending request")
		return
	}

	result.tiil, err = client.QueryIndicatorsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "securityinsight.ThreatIntelligenceIndicatorClient", "QueryIndicators", resp, "Failure responding to request")
		return
	}
	if result.tiil.hasNextLink() && result.tiil.IsEmpty() {
		err = result.NextWithContext(ctx)
		return
	}

	return
}

func (client ThreatIntelligenceIndicatorClient) queryIndicatorsNextResults(ctx context.Context, lastResults ThreatIntelligenceInformationList) (result ThreatIntelligenceInformationList, err error) {
	req, err := lastResults.threatIntelligenceInformationListPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "securityinsight.ThreatIntelligenceIndicatorClient", "queryIndicatorsNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.QueryIndicatorsSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "securityinsight.ThreatIntelligenceIndicatorClient", "queryIndicatorsNextResults", resp, "Failure sending next results request")
	}
	result, err = client.QueryIndicatorsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "securityinsight.ThreatIntelligenceIndicatorClient", "queryIndicatorsNextResults", resp, "Failure responding to next results request")
	}
	return
}

func (tiil ThreatIntelligenceInformationList) threatIntelligenceInformationListPreparer(ctx context.Context) (*http.Request, error) {
	if !tiil.hasNextLink() {
		return nil, nil
	}
	return autorest.Prepare((&http.Request{}).WithContext(ctx),
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(tiil.NextLink)))
}

func (client ThreatIntelligenceIndicatorClient) QueryIndicatorsPreparer(ctx context.Context, resourceGroupName string, workspaceName string, threatIntelligenceFilteringCriteria securityinsight.ThreatIntelligenceFilteringCriteria) (*http.Request, error) {
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
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.OperationalInsights/workspaces/{workspaceName}/providers/Microsoft.SecurityInsights/threatIntelligence/main/queryIndicators", pathParameters),
		autorest.WithJSON(threatIntelligenceFilteringCriteria),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// QueryIndicatorsSender sends the QueryIndicators request. The method will close the
// http.Response Body if it receives an error.
func (client ThreatIntelligenceIndicatorClient) QueryIndicatorsSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// QueryIndicatorsResponder handles the response to the QueryIndicators request. The method always
// closes the http.Response Body.
func (client ThreatIntelligenceIndicatorClient) QueryIndicatorsResponder(resp *http.Response) (result ThreatIntelligenceInformationList, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

func (client ThreatIntelligenceIndicatorClient) Create(ctx context.Context, resourceGroupName string, workspaceName string, name string, threatIntelligenceProperties ThreatIntelligenceIndicatorModel) (result ThreatIntelligenceInformationModel, err error) {
	if err := validation.Validate([]validation.Validation{
		{
			TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}},
		},
		{
			TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{
				{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
			},
		},
		{
			TargetValue: workspaceName,
			Constraints: []validation.Constraint{
				{Target: "workspaceName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "workspaceName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "workspaceName", Name: validation.Pattern, Rule: `^[A-Za-z0-9][A-Za-z0-9-]+[A-Za-z0-9]$`, Chain: nil},
			},
		},
	}); err != nil {
		return result, validation.NewError("securityinsight.ThreatIntelligenceIndicatorClient", "Create", "%+v", err.Error())
	}

	req, err := client.CreatePreparer(ctx, resourceGroupName, workspaceName, name, threatIntelligenceProperties)
	if err != nil {
		err = autorest.NewErrorWithError(err, "securityinsight.ThreatIntelligenceIndicatorClient", "Create", nil, "Failure preparing request")
		return
	}

	resp, err := client.CreateSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "securityinsight.ThreatIntelligenceIndicatorClient", "Create", resp, "Failure sending request")
		return
	}

	result, err = client.CreateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "securityinsight.ThreatIntelligenceIndicatorClient", "Create", resp, "Failure responding to request")
		return
	}

	return
}

func (client ThreatIntelligenceIndicatorClient) CreatePreparer(ctx context.Context, resourceGroupName string, workspaceName string, name string, threatIntelligenceProperties ThreatIntelligenceIndicatorModel) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"name":              autorest.Encode("path", name),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"workspaceName":     autorest.Encode("path", workspaceName),
	}

	const APIVersion = "2022-10-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.OperationalInsights/workspaces/{workspaceName}/providers/Microsoft.SecurityInsights/threatIntelligence/main/indicators/{name}", pathParameters),
		autorest.WithJSON(threatIntelligenceProperties),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

func (client ThreatIntelligenceIndicatorClient) CreateSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

func (client ThreatIntelligenceIndicatorClient) CreateResponder(resp *http.Response) (result ThreatIntelligenceInformationModel, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
