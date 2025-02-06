package appplatform

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

type DeploymentsDisableRemoteDebuggingOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *RemoteDebugging
}

// DeploymentsDisableRemoteDebugging ...
func (c AppPlatformClient) DeploymentsDisableRemoteDebugging(ctx context.Context, id DeploymentId) (result DeploymentsDisableRemoteDebuggingOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/disableRemoteDebugging", id.ID()),
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

// DeploymentsDisableRemoteDebuggingThenPoll performs DeploymentsDisableRemoteDebugging then polls until it's completed
func (c AppPlatformClient) DeploymentsDisableRemoteDebuggingThenPoll(ctx context.Context, id DeploymentId) error {
	result, err := c.DeploymentsDisableRemoteDebugging(ctx, id)
	if err != nil {
		return fmt.Errorf("performing DeploymentsDisableRemoteDebugging: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after DeploymentsDisableRemoteDebugging: %+v", err)
	}

	return nil
}
