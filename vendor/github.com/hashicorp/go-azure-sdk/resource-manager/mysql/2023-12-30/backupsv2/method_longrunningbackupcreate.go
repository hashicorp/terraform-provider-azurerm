package backupsv2

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

type LongRunningBackupCreateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *ServerBackupV2
}

// LongRunningBackupCreate ...
func (c BackupsV2Client) LongRunningBackupCreate(ctx context.Context, id BackupsV2Id, input ServerBackupV2) (result LongRunningBackupCreateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
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

// LongRunningBackupCreateThenPoll performs LongRunningBackupCreate then polls until it's completed
func (c BackupsV2Client) LongRunningBackupCreateThenPoll(ctx context.Context, id BackupsV2Id, input ServerBackupV2) error {
	result, err := c.LongRunningBackupCreate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing LongRunningBackupCreate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after LongRunningBackupCreate: %+v", err)
	}

	return nil
}
