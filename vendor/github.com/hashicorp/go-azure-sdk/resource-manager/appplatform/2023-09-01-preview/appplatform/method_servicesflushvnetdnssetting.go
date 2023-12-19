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

type ServicesFlushVnetDnsSettingOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// ServicesFlushVnetDnsSetting ...
func (c AppPlatformClient) ServicesFlushVnetDnsSetting(ctx context.Context, id commonids.SpringCloudServiceId) (result ServicesFlushVnetDnsSettingOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/flushVirtualNetworkDnsSettings", id.ID()),
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

// ServicesFlushVnetDnsSettingThenPoll performs ServicesFlushVnetDnsSetting then polls until it's completed
func (c AppPlatformClient) ServicesFlushVnetDnsSettingThenPoll(ctx context.Context, id commonids.SpringCloudServiceId) error {
	result, err := c.ServicesFlushVnetDnsSetting(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ServicesFlushVnetDnsSetting: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after ServicesFlushVnetDnsSetting: %+v", err)
	}

	return nil
}
