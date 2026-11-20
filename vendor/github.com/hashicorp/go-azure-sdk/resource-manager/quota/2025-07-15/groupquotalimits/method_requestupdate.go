package groupquotalimits

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

type RequestUpdateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *GroupQuotaLimitList
}

// RequestUpdate ...
func (c GroupQuotaLimitsClient) RequestUpdate(ctx context.Context, id GroupQuotaLimitId, input GroupQuotaLimitList) (result RequestUpdateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPatch,
		Path:       id.ID(),
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

// RequestUpdateThenPoll performs RequestUpdate then polls until it's completed
func (c GroupQuotaLimitsClient) RequestUpdateThenPoll(ctx context.Context, id GroupQuotaLimitId, input GroupQuotaLimitList) error {
	return c.RequestUpdateCallbackThenPoll(ctx, id, input, nil)
}

// RequestUpdateCallbackThenPoll performs RequestUpdate, runs the optional callback function, then polls until it's completed
func (c GroupQuotaLimitsClient) RequestUpdateCallbackThenPoll(ctx context.Context, id GroupQuotaLimitId, input GroupQuotaLimitList, callback func() error) error {
	result, err := c.RequestUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing RequestUpdate: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after RequestUpdate: %+v", err)
	}

	return nil
}
