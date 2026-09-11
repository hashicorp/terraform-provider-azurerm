package servicegateways

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

type UpdateServicesOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// UpdateServices ...
func (c ServiceGatewaysClient) UpdateServices(ctx context.Context, id ServiceGatewayId, input ServiceGatewayUpdateServicesRequest) (result UpdateServicesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusNoContent,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/updateServices", id.ID()),
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

// UpdateServicesThenPoll performs UpdateServices then polls until it's completed
func (c ServiceGatewaysClient) UpdateServicesThenPoll(ctx context.Context, id ServiceGatewayId, input ServiceGatewayUpdateServicesRequest) error {
	return c.UpdateServicesCallbackThenPoll(ctx, id, input, nil)
}

// UpdateServicesCallbackThenPoll performs UpdateServices, runs the optional callback function, then polls until it's completed
func (c ServiceGatewaysClient) UpdateServicesCallbackThenPoll(ctx context.Context, id ServiceGatewayId, input ServiceGatewayUpdateServicesRequest, callback func() error) error {
	result, err := c.UpdateServices(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing UpdateServices: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after UpdateServices: %+v", err)
	}

	return nil
}
