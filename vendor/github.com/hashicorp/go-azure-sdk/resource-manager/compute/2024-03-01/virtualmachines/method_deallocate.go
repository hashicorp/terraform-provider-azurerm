package virtualmachines

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

type DeallocateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

type DeallocateOperationOptions struct {
	Hibernate *bool
}

func DefaultDeallocateOperationOptions() DeallocateOperationOptions {
	return DeallocateOperationOptions{}
}

func (o DeallocateOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o DeallocateOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o DeallocateOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Hibernate != nil {
		out.Append("hibernate", fmt.Sprintf("%v", *o.Hibernate))
	}
	return &out
}

// Deallocate ...
func (c VirtualMachinesClient) Deallocate(ctx context.Context, id VirtualMachineId, options DeallocateOperationOptions) (result DeallocateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/deallocate", id.ID()),
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

// DeallocateThenPoll performs Deallocate then polls until it's completed
func (c VirtualMachinesClient) DeallocateThenPoll(ctx context.Context, id VirtualMachineId, options DeallocateOperationOptions) error {
	result, err := c.Deallocate(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing Deallocate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after Deallocate: %+v", err)
	}

	return nil
}
