package sqlvirtualmachines

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

type SqlVirtualMachineTroubleshootTroubleshootOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *SqlVMTroubleshooting
}

// SqlVirtualMachineTroubleshootTroubleshoot ...
func (c SqlVirtualMachinesClient) SqlVirtualMachineTroubleshootTroubleshoot(ctx context.Context, id SqlVirtualMachineId, input SqlVMTroubleshooting) (result SqlVirtualMachineTroubleshootTroubleshootOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/troubleshoot", id.ID()),
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

// SqlVirtualMachineTroubleshootTroubleshootThenPoll performs SqlVirtualMachineTroubleshootTroubleshoot then polls until it's completed
func (c SqlVirtualMachinesClient) SqlVirtualMachineTroubleshootTroubleshootThenPoll(ctx context.Context, id SqlVirtualMachineId, input SqlVMTroubleshooting) error {
	result, err := c.SqlVirtualMachineTroubleshootTroubleshoot(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing SqlVirtualMachineTroubleshootTroubleshoot: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after SqlVirtualMachineTroubleshootTroubleshoot: %+v", err)
	}

	return nil
}
