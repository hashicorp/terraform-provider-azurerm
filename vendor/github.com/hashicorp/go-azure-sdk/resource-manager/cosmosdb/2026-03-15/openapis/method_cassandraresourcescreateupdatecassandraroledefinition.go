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

type CassandraResourcesCreateUpdateCassandraRoleDefinitionOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *CassandraRoleDefinitionResource
}

// CassandraResourcesCreateUpdateCassandraRoleDefinition ...
func (c OpenapisClient) CassandraResourcesCreateUpdateCassandraRoleDefinition(ctx context.Context, id CassandraRoleDefinitionId, input CassandraRoleDefinitionResource) (result CassandraResourcesCreateUpdateCassandraRoleDefinitionOperationResponse, err error) {
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

// CassandraResourcesCreateUpdateCassandraRoleDefinitionThenPoll performs CassandraResourcesCreateUpdateCassandraRoleDefinition then polls until it's completed
func (c OpenapisClient) CassandraResourcesCreateUpdateCassandraRoleDefinitionThenPoll(ctx context.Context, id CassandraRoleDefinitionId, input CassandraRoleDefinitionResource) error {
	return c.CassandraResourcesCreateUpdateCassandraRoleDefinitionCallbackThenPoll(ctx, id, input, nil)
}

// CassandraResourcesCreateUpdateCassandraRoleDefinitionCallbackThenPoll performs CassandraResourcesCreateUpdateCassandraRoleDefinition, runs the optional callback function, then polls until it's completed
func (c OpenapisClient) CassandraResourcesCreateUpdateCassandraRoleDefinitionCallbackThenPoll(ctx context.Context, id CassandraRoleDefinitionId, input CassandraRoleDefinitionResource, callback func() error) error {
	result, err := c.CassandraResourcesCreateUpdateCassandraRoleDefinition(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing CassandraResourcesCreateUpdateCassandraRoleDefinition: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after CassandraResourcesCreateUpdateCassandraRoleDefinition: %+v", err)
	}

	return nil
}
