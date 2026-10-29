package openapis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TableResourcesCreateUpdateTableRoleDefinitionOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *TableRoleDefinitionResource
}

// TableResourcesCreateUpdateTableRoleDefinition ...
func (c OpenapisClient) TableResourcesCreateUpdateTableRoleDefinition(ctx context.Context, id TableRoleDefinitionId, input TableRoleDefinitionResource) (result TableResourcesCreateUpdateTableRoleDefinitionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPut,
		Path:       id.ID(),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	if err = req.Marshal(input); err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	result.Poller, err = resourcemanager.PollerFromResponse(resp, c.Client)
	if err != nil {
		return
	}

	return
}

// TableResourcesCreateUpdateTableRoleDefinitionThenPoll performs TableResourcesCreateUpdateTableRoleDefinition then polls until it's completed
func (c OpenapisClient) TableResourcesCreateUpdateTableRoleDefinitionThenPoll(ctx context.Context, id TableRoleDefinitionId, input TableRoleDefinitionResource) error {
	return c.TableResourcesCreateUpdateTableRoleDefinitionCallbackThenPoll(ctx, id, input, nil)
}

// TableResourcesCreateUpdateTableRoleDefinitionCallbackThenPoll performs TableResourcesCreateUpdateTableRoleDefinition, runs the optional callback function, then polls until it's completed
func (c OpenapisClient) TableResourcesCreateUpdateTableRoleDefinitionCallbackThenPoll(ctx context.Context, id TableRoleDefinitionId, input TableRoleDefinitionResource, callback func() error) error {
	result, err := c.TableResourcesCreateUpdateTableRoleDefinition(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing TableResourcesCreateUpdateTableRoleDefinition: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after TableResourcesCreateUpdateTableRoleDefinition: %+v", err)
	}

	return nil
}
