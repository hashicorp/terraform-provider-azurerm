package serviceresource

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

type ServicesDeleteOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

type ServicesDeleteOperationOptions struct {
	DeleteRunningTasks *bool
}

func DefaultServicesDeleteOperationOptions() ServicesDeleteOperationOptions {
	return ServicesDeleteOperationOptions{}
}

func (o ServicesDeleteOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ServicesDeleteOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ServicesDeleteOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.DeleteRunningTasks != nil {
		out.Append("deleteRunningTasks", fmt.Sprintf("%v", *o.DeleteRunningTasks))
	}
	return &out
}

// ServicesDelete ...
func (c ServiceResourceClient) ServicesDelete(ctx context.Context, id ServiceId, options ServicesDeleteOperationOptions) (result ServicesDeleteOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusNoContent,
			http.StatusOK,
		},
		HttpMethod:    http.MethodDelete,
		OptionsObject: options,
		Path:          id.ID(),
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

// ServicesDeleteThenPoll performs ServicesDelete then polls until it's completed
func (c ServiceResourceClient) ServicesDeleteThenPoll(ctx context.Context, id ServiceId, options ServicesDeleteOperationOptions) error {
	result, err := c.ServicesDelete(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing ServicesDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after ServicesDelete: %+v", err)
	}

	return nil
}
