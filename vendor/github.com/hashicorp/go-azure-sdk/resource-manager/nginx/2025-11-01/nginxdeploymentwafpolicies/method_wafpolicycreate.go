package nginxdeploymentwafpolicies

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

type WafPolicyCreateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *NginxDeploymentWafPolicy
}

// WafPolicyCreate ...
func (c NginxDeploymentWafPoliciesClient) WafPolicyCreate(ctx context.Context, id WafPolicyId, input NginxDeploymentWafPolicy) (result WafPolicyCreateOperationResponse, err error) {
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

// WafPolicyCreateThenPoll performs WafPolicyCreate then polls until it's completed
func (c NginxDeploymentWafPoliciesClient) WafPolicyCreateThenPoll(ctx context.Context, id WafPolicyId, input NginxDeploymentWafPolicy) error {
	return c.WafPolicyCreateCallbackThenPoll(ctx, id, input, nil)
}

// WafPolicyCreateCallbackThenPoll performs WafPolicyCreate, runs the optional callback function, then polls until it's completed
func (c NginxDeploymentWafPoliciesClient) WafPolicyCreateCallbackThenPoll(ctx context.Context, id WafPolicyId, input NginxDeploymentWafPolicy, callback func() error) error {
	result, err := c.WafPolicyCreate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing WafPolicyCreate: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after WafPolicyCreate: %+v", err)
	}

	return nil
}
