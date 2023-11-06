package backupshorttermretentionpolicies

import (
	"context"
	"fmt"
	"net/http"

<<<<<<< HEAD
=======
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
>>>>>>> 5e957238fca9519400c2479c7d1f73e3d1b0871c
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// Update ...
<<<<<<< HEAD
func (c BackupShortTermRetentionPoliciesClient) Update(ctx context.Context, id DatabaseId, input BackupShortTermRetentionPolicy) (result UpdateOperationResponse, err error) {
=======
func (c BackupShortTermRetentionPoliciesClient) Update(ctx context.Context, id commonids.SqlDatabaseId, input BackupShortTermRetentionPolicy) (result UpdateOperationResponse, err error) {
>>>>>>> 5e957238fca9519400c2479c7d1f73e3d1b0871c
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPatch,
		Path:       fmt.Sprintf("%s/backupShortTermRetentionPolicies/default", id.ID()),
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

// UpdateThenPoll performs Update then polls until it's completed
<<<<<<< HEAD
func (c BackupShortTermRetentionPoliciesClient) UpdateThenPoll(ctx context.Context, id DatabaseId, input BackupShortTermRetentionPolicy) error {
=======
func (c BackupShortTermRetentionPoliciesClient) UpdateThenPoll(ctx context.Context, id commonids.SqlDatabaseId, input BackupShortTermRetentionPolicy) error {
>>>>>>> 5e957238fca9519400c2479c7d1f73e3d1b0871c
	result, err := c.Update(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing Update: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after Update: %+v", err)
	}

	return nil
}
