package machinelearningcomputes

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

type ComputeDeleteOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

type ComputeDeleteOperationOptions struct {
	UnderlyingResourceAction *UnderlyingResourceAction
}

func DefaultComputeDeleteOperationOptions() ComputeDeleteOperationOptions {
	return ComputeDeleteOperationOptions{}
}

func (o ComputeDeleteOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ComputeDeleteOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ComputeDeleteOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.UnderlyingResourceAction != nil {
		out.Append("underlyingResourceAction", fmt.Sprintf("%v", *o.UnderlyingResourceAction))
	}
	return &out
}

// ComputeDelete ...
func (c MachineLearningComputesClient) ComputeDelete(ctx context.Context, id ComputeId, options ComputeDeleteOperationOptions) (result ComputeDeleteOperationResponse, err error) {
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

// ComputeDeleteThenPoll performs ComputeDelete then polls until it's completed
func (c MachineLearningComputesClient) ComputeDeleteThenPoll(ctx context.Context, id ComputeId, options ComputeDeleteOperationOptions) error {
	result, err := c.ComputeDelete(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing ComputeDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after ComputeDelete: %+v", err)
	}

	return nil
}
