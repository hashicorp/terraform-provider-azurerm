package replicationrecoveryplans

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

type PlannedFailoverOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *RecoveryPlan
}

// PlannedFailover ...
func (c ReplicationRecoveryPlansClient) PlannedFailover(ctx context.Context, id ReplicationRecoveryPlanId, input RecoveryPlanPlannedFailoverInput) (result PlannedFailoverOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/plannedFailover", id.ID()),
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

// PlannedFailoverThenPoll performs PlannedFailover then polls until it's completed
func (c ReplicationRecoveryPlansClient) PlannedFailoverThenPoll(ctx context.Context, id ReplicationRecoveryPlanId, input RecoveryPlanPlannedFailoverInput) error {
	result, err := c.PlannedFailover(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing PlannedFailover: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after PlannedFailover: %+v", err)
	}

	return nil
}
