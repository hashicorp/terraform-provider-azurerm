package snapshotpolicy

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

type SnapshotPoliciesDeleteOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// SnapshotPoliciesDelete ...
func (c SnapshotPolicyClient) SnapshotPoliciesDelete(ctx context.Context, id SnapshotPolicyId) (result SnapshotPoliciesDeleteOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusNoContent,
			http.StatusOK,
		},
		HttpMethod: http.MethodDelete,
		Path:       id.ID(),
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

// SnapshotPoliciesDeleteThenPoll performs SnapshotPoliciesDelete then polls until it's completed
func (c SnapshotPolicyClient) SnapshotPoliciesDeleteThenPoll(ctx context.Context, id SnapshotPolicyId) error {
	result, err := c.SnapshotPoliciesDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing SnapshotPoliciesDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after SnapshotPoliciesDelete: %+v", err)
	}

	return nil
}
