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

type RotateServiceAccountSigningKeysOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// RotateServiceAccountSigningKeys ...
func (c ManagedClustersClient) RotateServiceAccountSigningKeys(ctx context.Context, id commonids.KubernetesClusterId) (result RotateServiceAccountSigningKeysOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusNoContent,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/rotateServiceAccountSigningKeys", id.ID()),
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

// RotateServiceAccountSigningKeysThenPoll performs RotateServiceAccountSigningKeys then polls until it's completed
func (c ManagedClustersClient) RotateServiceAccountSigningKeysThenPoll(ctx context.Context, id commonids.KubernetesClusterId) error {
	result, err := c.RotateServiceAccountSigningKeys(ctx, id)
	if err != nil {
		return fmt.Errorf("performing RotateServiceAccountSigningKeys: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after RotateServiceAccountSigningKeys: %+v", err)
	}

	return nil
}
