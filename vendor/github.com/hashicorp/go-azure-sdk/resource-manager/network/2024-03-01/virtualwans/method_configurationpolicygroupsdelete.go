package virtualwans

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

type ConfigurationPolicyGroupsDeleteOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// ConfigurationPolicyGroupsDelete ...
func (c VirtualWANsClient) ConfigurationPolicyGroupsDelete(ctx context.Context, id ConfigurationPolicyGroupId) (result ConfigurationPolicyGroupsDeleteOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusNoContent,
			http.StatusOK,
		},
		HttpMethod: http.MethodDelete,
		Path:       id.ID(),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
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

// ConfigurationPolicyGroupsDeleteThenPoll performs ConfigurationPolicyGroupsDelete then polls until it's completed
func (c VirtualWANsClient) ConfigurationPolicyGroupsDeleteThenPoll(ctx context.Context, id ConfigurationPolicyGroupId) error {
	result, err := c.ConfigurationPolicyGroupsDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ConfigurationPolicyGroupsDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after ConfigurationPolicyGroupsDelete: %+v", err)
	}

	return nil
}
