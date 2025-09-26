package managedclusters

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

type ResetServicePrincipalProfileOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// ResetServicePrincipalProfile ...
func (c ManagedClustersClient) ResetServicePrincipalProfile(ctx context.Context, id commonids.KubernetesClusterId, input ManagedClusterServicePrincipalProfile) (result ResetServicePrincipalProfileOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/resetServicePrincipalProfile", id.ID()),
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

// ResetServicePrincipalProfileThenPoll performs ResetServicePrincipalProfile then polls until it's completed
func (c ManagedClustersClient) ResetServicePrincipalProfileThenPoll(ctx context.Context, id commonids.KubernetesClusterId, input ManagedClusterServicePrincipalProfile) error {
	result, err := c.ResetServicePrincipalProfile(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ResetServicePrincipalProfile: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after ResetServicePrincipalProfile: %+v", err)
	}

	return nil
}
