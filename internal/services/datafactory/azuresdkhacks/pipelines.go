// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azuresdkhacks

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/kermit/sdk/datafactory/2018-06-01/datafactory"
)

// TODO4.0: check if the workaround could be removed.
// Workaround for https://github.com/hashicorp/terraform-provider-azurerm/issues/24758
// Tracked on https://github.com/Azure/azure-rest-api-specs/issues/27816
// This file is almost copied from https://github.com/tombuildsstuff/kermit/blob/main/sdk/datafactory/2018-06-01/datafactory/pipelines.go
// Added a custom client to use custom `PipelineResource`.

type PipelinesClient struct {
	OriginalClient *datafactory.PipelinesClient
}

func (client PipelinesClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, factoryName string, pipelineName string, pipeline PipelineResource, ifMatch string) (result PipelineResource, err error) {
	if err := validation.Validate([]validation.Validation{
		{
			TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{
				{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil},
			},
		},
		{
			TargetValue: factoryName,
			Constraints: []validation.Constraint{
				{Target: "factoryName", Name: validation.MaxLength, Rule: 63, Chain: nil},
				{Target: "factoryName", Name: validation.MinLength, Rule: 3, Chain: nil},
				{Target: "factoryName", Name: validation.Pattern, Rule: `^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$`, Chain: nil},
			},
		},
		{
			TargetValue: pipelineName,
			Constraints: []validation.Constraint{
				{Target: "pipelineName", Name: validation.MaxLength, Rule: 260, Chain: nil},
				{Target: "pipelineName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "pipelineName", Name: validation.Pattern, Rule: `^[A-Za-z0-9_][^<>*#.%&:\\+?/]*$`, Chain: nil},
			},
		},
		{
			TargetValue: pipeline,
			Constraints: []validation.Constraint{{
				Target: "pipeline.Pipeline", Name: validation.Null, Rule: true,
				Chain: []validation.Constraint{
					{
						Target: "pipeline.Pipeline.Concurrency", Name: validation.Null, Rule: false,
						Chain: []validation.Constraint{{Target: "pipeline.Pipeline.Concurrency", Name: validation.InclusiveMinimum, Rule: int64(1), Chain: nil}},
					},
				},
			}},
		},
	}); err != nil {
		return result, validation.NewError("datafactory.PipelinesClient", "CreateOrUpdate", err.Error())
	}

	req, err := client.CreateOrUpdatePreparer(ctx, resourceGroupName, factoryName, pipelineName, pipeline, ifMatch)
	if err != nil {
		err = autorest.NewErrorWithError(err, "datafactory.PipelinesClient", "CreateOrUpdate", nil, "Failure preparing request")
		return
	}

	resp, err := client.OriginalClient.CreateOrUpdateSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "datafactory.PipelinesClient", "CreateOrUpdate", resp, "Failure sending request")
		return
	}

	result, err = client.CreateOrUpdateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "datafactory.PipelinesClient", "CreateOrUpdate", resp, "Failure responding to request")
		return
	}

	return
}

// CreateOrUpdatePreparer prepares the CreateOrUpdate request.
func (client PipelinesClient) CreateOrUpdatePreparer(ctx context.Context, resourceGroupName string, factoryName string, pipelineName string, pipeline PipelineResource, ifMatch string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"factoryName":       autorest.Encode("path", factoryName),
		"pipelineName":      autorest.Encode("path", pipelineName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.OriginalClient.SubscriptionID),
	}

	const APIVersion = "2018-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.OriginalClient.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataFactory/factories/{factoryName}/pipelines/{pipelineName}", pathParameters),
		autorest.WithJSON(pipeline),
		autorest.WithQueryParameters(queryParameters))
	if len(ifMatch) > 0 {
		preparer = autorest.DecoratePreparer(preparer,
			autorest.WithHeader("If-Match", autorest.String(ifMatch)))
	}
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

func (client PipelinesClient) CreateOrUpdateResponder(resp *http.Response) (result PipelineResource, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

func (client PipelinesClient) Get(ctx context.Context, resourceGroupName string, factoryName string, pipelineName string, ifNoneMatch string) (result PipelineResource, err error) {
	if err := validation.Validate([]validation.Validation{
		{
			TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{
				{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil},
			},
		},
		{
			TargetValue: factoryName,
			Constraints: []validation.Constraint{
				{Target: "factoryName", Name: validation.MaxLength, Rule: 63, Chain: nil},
				{Target: "factoryName", Name: validation.MinLength, Rule: 3, Chain: nil},
				{Target: "factoryName", Name: validation.Pattern, Rule: `^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$`, Chain: nil},
			},
		},
		{
			TargetValue: pipelineName,
			Constraints: []validation.Constraint{
				{Target: "pipelineName", Name: validation.MaxLength, Rule: 260, Chain: nil},
				{Target: "pipelineName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "pipelineName", Name: validation.Pattern, Rule: `^[A-Za-z0-9_][^<>*#.%&:\\+?/]*$`, Chain: nil},
			},
		},
	}); err != nil {
		return result, validation.NewError("datafactory.PipelinesClient", "Get", err.Error())
	}

	req, err := client.GetPreparer(ctx, resourceGroupName, factoryName, pipelineName, ifNoneMatch)
	if err != nil {
		err = autorest.NewErrorWithError(err, "datafactory.PipelinesClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.OriginalClient.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "datafactory.PipelinesClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "datafactory.PipelinesClient", "Get", resp, "Failure responding to request")
		return
	}

	return
}

func (client PipelinesClient) GetPreparer(ctx context.Context, resourceGroupName string, factoryName string, pipelineName string, ifNoneMatch string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"factoryName":       autorest.Encode("path", factoryName),
		"pipelineName":      autorest.Encode("path", pipelineName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.OriginalClient.SubscriptionID),
	}

	const APIVersion = "2018-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.OriginalClient.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataFactory/factories/{factoryName}/pipelines/{pipelineName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	if len(ifNoneMatch) > 0 {
		preparer = autorest.DecoratePreparer(preparer,
			autorest.WithHeader("If-None-Match", autorest.String(ifNoneMatch)))
	}
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

func (client PipelinesClient) GetResponder(resp *http.Response) (result PipelineResource, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNotModified),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
