package apimanagementservice

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

type ApplyNetworkConfigurationUpdatesOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// ApplyNetworkConfigurationUpdates ...
func (c ApiManagementServiceClient) ApplyNetworkConfigurationUpdates(ctx context.Context, id ServiceId, input ApiManagementServiceApplyNetworkConfigurationParameters) (result ApplyNetworkConfigurationUpdatesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/applynetworkconfigurationupdates", id.ID()),
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

// ApplyNetworkConfigurationUpdatesThenPoll performs ApplyNetworkConfigurationUpdates then polls until it's completed
func (c ApiManagementServiceClient) ApplyNetworkConfigurationUpdatesThenPoll(ctx context.Context, id ServiceId, input ApiManagementServiceApplyNetworkConfigurationParameters) error {
	result, err := c.ApplyNetworkConfigurationUpdates(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ApplyNetworkConfigurationUpdates: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after ApplyNetworkConfigurationUpdates: %+v", err)
	}

	return nil
}
