package provisionedclusterinstances

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

type ProvisionedClusterInstancesDeleteOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// ProvisionedClusterInstancesDelete ...
func (c ProvisionedClusterInstancesClient) ProvisionedClusterInstancesDelete(ctx context.Context, id commonids.ScopeId) (result ProvisionedClusterInstancesDeleteOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusNoContent,
		},
		HttpMethod: http.MethodDelete,
		Path:       fmt.Sprintf("%s/providers/Microsoft.HybridContainerService/provisionedClusterInstances/default", id.ID()),
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

// ProvisionedClusterInstancesDeleteThenPoll performs ProvisionedClusterInstancesDelete then polls until it's completed
func (c ProvisionedClusterInstancesClient) ProvisionedClusterInstancesDeleteThenPoll(ctx context.Context, id commonids.ScopeId) error {
	result, err := c.ProvisionedClusterInstancesDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ProvisionedClusterInstancesDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after ProvisionedClusterInstancesDelete: %+v", err)
	}

	return nil
}
