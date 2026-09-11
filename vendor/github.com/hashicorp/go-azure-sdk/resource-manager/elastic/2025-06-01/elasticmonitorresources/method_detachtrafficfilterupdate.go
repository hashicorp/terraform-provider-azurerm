package elasticmonitorresources

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

type DetachTrafficFilterUpdateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

type DetachTrafficFilterUpdateOperationOptions struct {
	RulesetId *string
}

func DefaultDetachTrafficFilterUpdateOperationOptions() DetachTrafficFilterUpdateOperationOptions {
	return DetachTrafficFilterUpdateOperationOptions{}
}

func (o DetachTrafficFilterUpdateOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o DetachTrafficFilterUpdateOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o DetachTrafficFilterUpdateOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.RulesetId != nil {
		out.Append("rulesetId", fmt.Sprintf("%v", *o.RulesetId))
	}
	return &out
}

// DetachTrafficFilterUpdate ...
func (c ElasticMonitorResourcesClient) DetachTrafficFilterUpdate(ctx context.Context, id MonitorId, options DetachTrafficFilterUpdateOperationOptions) (result DetachTrafficFilterUpdateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/detachTrafficFilter", id.ID()),
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

// DetachTrafficFilterUpdateThenPoll performs DetachTrafficFilterUpdate then polls until it's completed
func (c ElasticMonitorResourcesClient) DetachTrafficFilterUpdateThenPoll(ctx context.Context, id MonitorId, options DetachTrafficFilterUpdateOperationOptions) error {
	return c.DetachTrafficFilterUpdateCallbackThenPoll(ctx, id, options, nil)
}

// DetachTrafficFilterUpdateCallbackThenPoll performs DetachTrafficFilterUpdate, runs the optional callback function, then polls until it's completed
func (c ElasticMonitorResourcesClient) DetachTrafficFilterUpdateCallbackThenPoll(ctx context.Context, id MonitorId, options DetachTrafficFilterUpdateOperationOptions, callback func() error) error {
	result, err := c.DetachTrafficFilterUpdate(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing DetachTrafficFilterUpdate: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after DetachTrafficFilterUpdate: %+v", err)
	}

	return nil
}
