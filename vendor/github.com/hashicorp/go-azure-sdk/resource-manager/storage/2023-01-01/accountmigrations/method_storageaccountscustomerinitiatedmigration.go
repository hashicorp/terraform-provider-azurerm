package accountmigrations

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageAccountsCustomerInitiatedMigrationOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// StorageAccountsCustomerInitiatedMigration ...
func (c AccountMigrationsClient) StorageAccountsCustomerInitiatedMigration(ctx context.Context, id commonids.StorageAccountId, input StorageAccountMigration) (result StorageAccountsCustomerInitiatedMigrationOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/startAccountMigration", id.ID()),
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

// StorageAccountsCustomerInitiatedMigrationThenPoll performs StorageAccountsCustomerInitiatedMigration then polls until it's completed
func (c AccountMigrationsClient) StorageAccountsCustomerInitiatedMigrationThenPoll(ctx context.Context, id commonids.StorageAccountId, input StorageAccountMigration) error {
	result, err := c.StorageAccountsCustomerInitiatedMigration(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing StorageAccountsCustomerInitiatedMigration: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after StorageAccountsCustomerInitiatedMigration: %+v", err)
	}

	return nil
}
