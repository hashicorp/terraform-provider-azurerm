package hcpopenshiftclusters

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

type RequestAdminCredentialOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *HcpOpenShiftClusterAdminCredential
}

// RequestAdminCredential ...
func (c HcpOpenShiftClustersClient) RequestAdminCredential(ctx context.Context, id HcpOpenShiftClusterId, input HcpOpenShiftClusterAdminCredentialRequest) (result RequestAdminCredentialOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/requestAdminCredential", id.ID()),
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

// RequestAdminCredentialThenPoll performs RequestAdminCredential then polls until it's completed
func (c HcpOpenShiftClustersClient) RequestAdminCredentialThenPoll(ctx context.Context, id HcpOpenShiftClusterId, input HcpOpenShiftClusterAdminCredentialRequest) error {
	return c.RequestAdminCredentialCallbackThenPoll(ctx, id, input, nil)
}

// RequestAdminCredentialCallbackThenPoll performs RequestAdminCredential, runs the optional callback function, then polls until it's completed
func (c HcpOpenShiftClustersClient) RequestAdminCredentialCallbackThenPoll(ctx context.Context, id HcpOpenShiftClusterId, input HcpOpenShiftClusterAdminCredentialRequest, callback func() error) error {
	result, err := c.RequestAdminCredential(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing RequestAdminCredential: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after RequestAdminCredential: %+v", err)
	}

	return nil
}
