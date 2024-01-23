package webapps

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

type InstallSiteExtensionSlotOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *SiteExtensionInfo
}

// InstallSiteExtensionSlot ...
func (c WebAppsClient) InstallSiteExtensionSlot(ctx context.Context, id SlotSiteExtensionId) (result InstallSiteExtensionSlotOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
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

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

<<<<<<< HEAD
=======
	if err = resp.Unmarshal(&result.Model); err != nil {
		return
	}

>>>>>>> bd17dc90f0 (finish refactor for go-azure-sdk first pass)
	result.Poller, err = resourcemanager.PollerFromResponse(resp, c.Client)
	if err != nil {
		return
	}

	return
}

// InstallSiteExtensionSlotThenPoll performs InstallSiteExtensionSlot then polls until it's completed
func (c WebAppsClient) InstallSiteExtensionSlotThenPoll(ctx context.Context, id SlotSiteExtensionId) error {
	result, err := c.InstallSiteExtensionSlot(ctx, id)
	if err != nil {
		return fmt.Errorf("performing InstallSiteExtensionSlot: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after InstallSiteExtensionSlot: %+v", err)
	}

	return nil
}