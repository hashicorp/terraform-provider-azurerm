package storageaccounts

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

type HierarchicalNamespaceMigrationOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

type HierarchicalNamespaceMigrationOperationOptions struct {
	RequestType *string
}

func DefaultHierarchicalNamespaceMigrationOperationOptions() HierarchicalNamespaceMigrationOperationOptions {
	return HierarchicalNamespaceMigrationOperationOptions{}
}

func (o HierarchicalNamespaceMigrationOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o HierarchicalNamespaceMigrationOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o HierarchicalNamespaceMigrationOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.RequestType != nil {
		out.Append("requestType", fmt.Sprintf("%v", *o.RequestType))
	}
	return &out
}

// HierarchicalNamespaceMigration ...
func (c StorageAccountsClient) HierarchicalNamespaceMigration(ctx context.Context, id commonids.StorageAccountId, options HierarchicalNamespaceMigrationOperationOptions) (result HierarchicalNamespaceMigrationOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/hnsonmigration", id.ID()),
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

// HierarchicalNamespaceMigrationThenPoll performs HierarchicalNamespaceMigration then polls until it's completed
func (c StorageAccountsClient) HierarchicalNamespaceMigrationThenPoll(ctx context.Context, id commonids.StorageAccountId, options HierarchicalNamespaceMigrationOperationOptions) error {
	result, err := c.HierarchicalNamespaceMigration(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing HierarchicalNamespaceMigration: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after HierarchicalNamespaceMigration: %+v", err)
	}

	return nil
}
