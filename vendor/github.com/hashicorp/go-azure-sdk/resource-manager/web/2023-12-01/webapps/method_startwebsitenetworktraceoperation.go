package webapps

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

type StartWebSiteNetworkTraceOperationOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]NetworkTrace
}

type StartWebSiteNetworkTraceOperationOperationOptions struct {
	DurationInSeconds *int64
	MaxFrameLength    *int64
	SasURL            *string
}

func DefaultStartWebSiteNetworkTraceOperationOperationOptions() StartWebSiteNetworkTraceOperationOperationOptions {
	return StartWebSiteNetworkTraceOperationOperationOptions{}
}

func (o StartWebSiteNetworkTraceOperationOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o StartWebSiteNetworkTraceOperationOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o StartWebSiteNetworkTraceOperationOperationOptions) ToQuery() *client.QueryParams {
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

// StartWebSiteNetworkTraceOperation ...
func (c WebAppsClient) StartWebSiteNetworkTraceOperation(ctx context.Context, id commonids.AppServiceId, options StartWebSiteNetworkTraceOperationOperationOptions) (result StartWebSiteNetworkTraceOperationOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/networkTrace/startOperation", id.ID()),
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

// StartWebSiteNetworkTraceOperationThenPoll performs StartWebSiteNetworkTraceOperation then polls until it's completed
func (c WebAppsClient) StartWebSiteNetworkTraceOperationThenPoll(ctx context.Context, id commonids.AppServiceId, options StartWebSiteNetworkTraceOperationOperationOptions) error {
	result, err := c.StartWebSiteNetworkTraceOperation(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing StartWebSiteNetworkTraceOperation: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after StartWebSiteNetworkTraceOperation: %+v", err)
	}

	return nil
}
