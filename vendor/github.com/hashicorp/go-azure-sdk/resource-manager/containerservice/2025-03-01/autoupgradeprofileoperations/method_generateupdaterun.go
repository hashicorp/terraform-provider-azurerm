package autoupgradeprofileoperations

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

type GenerateUpdateRunOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *GenerateResponse
}

// GenerateUpdateRun ...
func (c AutoUpgradeProfileOperationsClient) GenerateUpdateRun(ctx context.Context, id AutoUpgradeProfileId) (result GenerateUpdateRunOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/generateUpdateRun", id.ID()),
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

// GenerateUpdateRunThenPoll performs GenerateUpdateRun then polls until it's completed
func (c AutoUpgradeProfileOperationsClient) GenerateUpdateRunThenPoll(ctx context.Context, id AutoUpgradeProfileId) error {
	result, err := c.GenerateUpdateRun(ctx, id)
	if err != nil {
		return fmt.Errorf("performing GenerateUpdateRun: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after GenerateUpdateRun: %+v", err)
	}

	return nil
}
