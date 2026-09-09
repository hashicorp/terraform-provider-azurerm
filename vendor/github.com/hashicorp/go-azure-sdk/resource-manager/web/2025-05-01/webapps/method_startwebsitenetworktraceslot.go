package webapps

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StartWebSiteNetworkTraceSlotOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *string
}

type StartWebSiteNetworkTraceSlotOperationOptions struct {
	DurationInSeconds *int64
	MaxFrameLength    *int64
	SasURL            *string
}

func DefaultStartWebSiteNetworkTraceSlotOperationOptions() StartWebSiteNetworkTraceSlotOperationOptions {
	return StartWebSiteNetworkTraceSlotOperationOptions{}
}

func (o StartWebSiteNetworkTraceSlotOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o StartWebSiteNetworkTraceSlotOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o StartWebSiteNetworkTraceSlotOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.DurationInSeconds != nil {
		out.Append("durationInSeconds", fmt.Sprintf("%v", *o.DurationInSeconds))
	}
	if o.MaxFrameLength != nil {
		out.Append("maxFrameLength", fmt.Sprintf("%v", *o.MaxFrameLength))
	}
	if o.SasURL != nil {
		out.Append("sasUrl", fmt.Sprintf("%v", *o.SasURL))
	}
	return &out
}

// StartWebSiteNetworkTraceSlot ...
func (c WebAppsClient) StartWebSiteNetworkTraceSlot(ctx context.Context, id SlotId, options StartWebSiteNetworkTraceSlotOperationOptions) (result StartWebSiteNetworkTraceSlotOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/networkTrace/start", id.ID()),
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

	var model string
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
