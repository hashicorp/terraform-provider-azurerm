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

type AccessPolicyAssignmentDeleteOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// AccessPolicyAssignmentDelete ...
func (c AADClient) AccessPolicyAssignmentDelete(ctx context.Context, id AccessPolicyAssignmentId) (result AccessPolicyAssignmentDeleteOperationResponse, err error) {
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

// AccessPolicyAssignmentDeleteThenPoll performs AccessPolicyAssignmentDelete then polls until it's completed
func (c AADClient) AccessPolicyAssignmentDeleteThenPoll(ctx context.Context, id AccessPolicyAssignmentId) error {
	result, err := c.AccessPolicyAssignmentDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing AccessPolicyAssignmentDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after AccessPolicyAssignmentDelete: %+v", err)
	}

	return nil
}
