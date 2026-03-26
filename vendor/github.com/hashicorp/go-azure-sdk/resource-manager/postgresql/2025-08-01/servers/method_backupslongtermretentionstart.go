package servers

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

type BackupsLongTermRetentionStartOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *BackupsLongTermRetentionResponse
}

// BackupsLongTermRetentionStart ...
func (c ServersClient) BackupsLongTermRetentionStart(ctx context.Context, id FlexibleServerId, input BackupsLongTermRetentionRequest) (result BackupsLongTermRetentionStartOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/startLtrBackup", id.ID()),
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

// BackupsLongTermRetentionStartThenPoll performs BackupsLongTermRetentionStart then polls until it's completed
func (c ServersClient) BackupsLongTermRetentionStartThenPoll(ctx context.Context, id FlexibleServerId, input BackupsLongTermRetentionRequest) error {
	result, err := c.BackupsLongTermRetentionStart(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing BackupsLongTermRetentionStart: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after BackupsLongTermRetentionStart: %+v", err)
	}

	return nil
}
