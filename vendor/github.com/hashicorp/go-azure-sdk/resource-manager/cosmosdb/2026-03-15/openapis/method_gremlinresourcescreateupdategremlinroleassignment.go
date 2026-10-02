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

type GremlinResourcesCreateUpdateGremlinRoleAssignmentOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *GremlinRoleAssignmentResource
}

// GremlinResourcesCreateUpdateGremlinRoleAssignment ...
func (c OpenapisClient) GremlinResourcesCreateUpdateGremlinRoleAssignment(ctx context.Context, id GremlinRoleAssignmentId, input GremlinRoleAssignmentResource) (result GremlinResourcesCreateUpdateGremlinRoleAssignmentOperationResponse, err error) {
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

// GremlinResourcesCreateUpdateGremlinRoleAssignmentThenPoll performs GremlinResourcesCreateUpdateGremlinRoleAssignment then polls until it's completed
func (c OpenapisClient) GremlinResourcesCreateUpdateGremlinRoleAssignmentThenPoll(ctx context.Context, id GremlinRoleAssignmentId, input GremlinRoleAssignmentResource) error {
	return c.GremlinResourcesCreateUpdateGremlinRoleAssignmentCallbackThenPoll(ctx, id, input, nil)
}

// GremlinResourcesCreateUpdateGremlinRoleAssignmentCallbackThenPoll performs GremlinResourcesCreateUpdateGremlinRoleAssignment, runs the optional callback function, then polls until it's completed
func (c OpenapisClient) GremlinResourcesCreateUpdateGremlinRoleAssignmentCallbackThenPoll(ctx context.Context, id GremlinRoleAssignmentId, input GremlinRoleAssignmentResource, callback func() error) error {
	result, err := c.GremlinResourcesCreateUpdateGremlinRoleAssignment(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing GremlinResourcesCreateUpdateGremlinRoleAssignment: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after GremlinResourcesCreateUpdateGremlinRoleAssignment: %+v", err)
	}

	return nil
}
