package cloudendpointresource

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

type CloudEndpointsPostBackupOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// CloudEndpointsPostBackup ...
func (c CloudEndpointResourceClient) CloudEndpointsPostBackup(ctx context.Context, id CloudEndpointId, input BackupRequest) (result CloudEndpointsPostBackupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/postbackup", id.ID()),
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

// CloudEndpointsPostBackupThenPoll performs CloudEndpointsPostBackup then polls until it's completed
func (c CloudEndpointResourceClient) CloudEndpointsPostBackupThenPoll(ctx context.Context, id CloudEndpointId, input BackupRequest) error {
	result, err := c.CloudEndpointsPostBackup(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing CloudEndpointsPostBackup: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after CloudEndpointsPostBackup: %+v", err)
	}

	return nil
}
