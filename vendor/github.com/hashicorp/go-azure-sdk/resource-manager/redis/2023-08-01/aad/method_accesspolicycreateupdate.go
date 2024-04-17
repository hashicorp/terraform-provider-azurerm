package aad

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

type AccessPolicyCreateUpdateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *RedisCacheAccessPolicy
}

// AccessPolicyCreateUpdate ...
func (c AADClient) AccessPolicyCreateUpdate(ctx context.Context, id AccessPolicyId, input RedisCacheAccessPolicy) (result AccessPolicyCreateUpdateOperationResponse, err error) {
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

// AccessPolicyCreateUpdateThenPoll performs AccessPolicyCreateUpdate then polls until it's completed
func (c AADClient) AccessPolicyCreateUpdateThenPoll(ctx context.Context, id AccessPolicyId, input RedisCacheAccessPolicy) error {
	result, err := c.AccessPolicyCreateUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing AccessPolicyCreateUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after AccessPolicyCreateUpdate: %+v", err)
	}

	return nil
}
