package containerinstance

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

type ContainerGroupsCreateOrUpdateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *ContainerGroup
}

// ContainerGroupsCreateOrUpdate ...
func (c ContainerInstanceClient) ContainerGroupsCreateOrUpdate(ctx context.Context, id ContainerGroupId, input ContainerGroup) (result ContainerGroupsCreateOrUpdateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusCreated,
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

// ContainerGroupsCreateOrUpdateThenPoll performs ContainerGroupsCreateOrUpdate then polls until it's completed
func (c ContainerInstanceClient) ContainerGroupsCreateOrUpdateThenPoll(ctx context.Context, id ContainerGroupId, input ContainerGroup) error {
	return c.ContainerGroupsCreateOrUpdateCallbackThenPoll(ctx, id, input, nil)
}

// ContainerGroupsCreateOrUpdateCallbackThenPoll performs ContainerGroupsCreateOrUpdate, runs the optional callback function, then polls until it's completed
func (c ContainerInstanceClient) ContainerGroupsCreateOrUpdateCallbackThenPoll(ctx context.Context, id ContainerGroupId, input ContainerGroup, callback func() error) error {
	result, err := c.ContainerGroupsCreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ContainerGroupsCreateOrUpdate: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after ContainerGroupsCreateOrUpdate: %+v", err)
	}

	return nil
}
