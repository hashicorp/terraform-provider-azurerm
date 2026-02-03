package virtualwans

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

type VpnSitesConfigurationDownloadOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// VpnSitesConfigurationDownload ...
func (c VirtualWANsClient) VpnSitesConfigurationDownload(ctx context.Context, id VirtualWANId, input GetVpnSitesConfigurationRequest) (result VpnSitesConfigurationDownloadOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/vpnConfiguration", id.ID()),
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

// VpnSitesConfigurationDownloadThenPoll performs VpnSitesConfigurationDownload then polls until it's completed
func (c VirtualWANsClient) VpnSitesConfigurationDownloadThenPoll(ctx context.Context, id VirtualWANId, input GetVpnSitesConfigurationRequest) error {
	result, err := c.VpnSitesConfigurationDownload(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing VpnSitesConfigurationDownload: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after VpnSitesConfigurationDownload: %+v", err)
	}

	return nil
}
