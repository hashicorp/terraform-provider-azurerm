package webapps

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

type MigrateStorageOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *StorageMigrationResponse
}

type MigrateStorageOperationOptions struct {
	SubscriptionName *string
}

func DefaultMigrateStorageOperationOptions() MigrateStorageOperationOptions {
	return MigrateStorageOperationOptions{}
}

func (o MigrateStorageOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o MigrateStorageOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o MigrateStorageOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.SubscriptionName != nil {
		out.Append("subscriptionName", fmt.Sprintf("%v", *o.SubscriptionName))
	}
	return &out
}

// MigrateStorage ...
func (c WebAppsClient) MigrateStorage(ctx context.Context, id commonids.AppServiceId, input StorageMigrationOptions, options MigrateStorageOperationOptions) (result MigrateStorageOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPut,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/migrate", id.ID()),
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

// MigrateStorageThenPoll performs MigrateStorage then polls until it's completed
func (c WebAppsClient) MigrateStorageThenPoll(ctx context.Context, id commonids.AppServiceId, input StorageMigrationOptions, options MigrateStorageOperationOptions) error {
	result, err := c.MigrateStorage(ctx, id, input, options)
	if err != nil {
		return fmt.Errorf("performing MigrateStorage: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after MigrateStorage: %+v", err)
	}

	return nil
}
