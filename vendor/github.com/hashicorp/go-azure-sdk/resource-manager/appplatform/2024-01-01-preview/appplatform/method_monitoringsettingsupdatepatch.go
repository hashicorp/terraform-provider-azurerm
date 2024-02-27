package appplatform

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

type MonitoringSettingsUpdatePatchOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *MonitoringSettingResource
}

// MonitoringSettingsUpdatePatch ...
func (c AppPlatformClient) MonitoringSettingsUpdatePatch(ctx context.Context, id commonids.SpringCloudServiceId, input MonitoringSettingResource) (result MonitoringSettingsUpdatePatchOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPatch,
		Path:       fmt.Sprintf("%s/monitoringSettings/default", id.ID()),
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

// MonitoringSettingsUpdatePatchThenPoll performs MonitoringSettingsUpdatePatch then polls until it's completed
func (c AppPlatformClient) MonitoringSettingsUpdatePatchThenPoll(ctx context.Context, id commonids.SpringCloudServiceId, input MonitoringSettingResource) error {
	result, err := c.MonitoringSettingsUpdatePatch(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing MonitoringSettingsUpdatePatch: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after MonitoringSettingsUpdatePatch: %+v", err)
	}

	return nil
}
