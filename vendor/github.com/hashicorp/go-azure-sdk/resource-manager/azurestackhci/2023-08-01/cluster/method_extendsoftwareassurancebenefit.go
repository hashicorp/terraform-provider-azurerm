package cluster

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

type ExtendSoftwareAssuranceBenefitOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// ExtendSoftwareAssuranceBenefit ...
func (c ClusterClient) ExtendSoftwareAssuranceBenefit(ctx context.Context, id ClusterId, input SoftwareAssuranceChangeRequest) (result ExtendSoftwareAssuranceBenefitOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/extendSoftwareAssuranceBenefit", id.ID()),
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

// ExtendSoftwareAssuranceBenefitThenPoll performs ExtendSoftwareAssuranceBenefit then polls until it's completed
func (c ClusterClient) ExtendSoftwareAssuranceBenefitThenPoll(ctx context.Context, id ClusterId, input SoftwareAssuranceChangeRequest) error {
	result, err := c.ExtendSoftwareAssuranceBenefit(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ExtendSoftwareAssuranceBenefit: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after ExtendSoftwareAssuranceBenefit: %+v", err)
	}

	return nil
}
