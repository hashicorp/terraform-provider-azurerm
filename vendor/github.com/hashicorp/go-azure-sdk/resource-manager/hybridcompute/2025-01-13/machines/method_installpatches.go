package machines

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

type InstallPatchesOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *MachineInstallPatchesResult
}

// InstallPatches ...
func (c MachinesClient) InstallPatches(ctx context.Context, id MachineId, input MachineInstallPatchesParameters) (result InstallPatchesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/installPatches", id.ID()),
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

// InstallPatchesThenPoll performs InstallPatches then polls until it's completed
func (c MachinesClient) InstallPatchesThenPoll(ctx context.Context, id MachineId, input MachineInstallPatchesParameters) error {
	return c.InstallPatchesCallbackThenPoll(ctx, id, input, nil)
}

// InstallPatchesCallbackThenPoll performs InstallPatches, runs the optional callback function, then polls until it's completed
func (c MachinesClient) InstallPatchesCallbackThenPoll(ctx context.Context, id MachineId, input MachineInstallPatchesParameters, callback func() error) error {
	result, err := c.InstallPatches(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing InstallPatches: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after InstallPatches: %+v", err)
	}

	return nil
}
