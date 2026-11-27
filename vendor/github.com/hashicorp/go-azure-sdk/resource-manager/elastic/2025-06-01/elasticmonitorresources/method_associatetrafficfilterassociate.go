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

type AssociateTrafficFilterAssociateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

type AssociateTrafficFilterAssociateOperationOptions struct {
	RulesetId *string
}

func DefaultAssociateTrafficFilterAssociateOperationOptions() AssociateTrafficFilterAssociateOperationOptions {
	return AssociateTrafficFilterAssociateOperationOptions{}
}

func (o AssociateTrafficFilterAssociateOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o AssociateTrafficFilterAssociateOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o AssociateTrafficFilterAssociateOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.RulesetId != nil {
		out.Append("rulesetId", fmt.Sprintf("%v", *o.RulesetId))
	}
	return &out
}

// AssociateTrafficFilterAssociate ...
func (c ElasticMonitorResourcesClient) AssociateTrafficFilterAssociate(ctx context.Context, id MonitorId, options AssociateTrafficFilterAssociateOperationOptions) (result AssociateTrafficFilterAssociateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/associateTrafficFilter", id.ID()),
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

// AssociateTrafficFilterAssociateThenPoll performs AssociateTrafficFilterAssociate then polls until it's completed
func (c ElasticMonitorResourcesClient) AssociateTrafficFilterAssociateThenPoll(ctx context.Context, id MonitorId, options AssociateTrafficFilterAssociateOperationOptions) error {
	return c.AssociateTrafficFilterAssociateCallbackThenPoll(ctx, id, options, nil)
}

// AssociateTrafficFilterAssociateCallbackThenPoll performs AssociateTrafficFilterAssociate, runs the optional callback function, then polls until it's completed
func (c ElasticMonitorResourcesClient) AssociateTrafficFilterAssociateCallbackThenPoll(ctx context.Context, id MonitorId, options AssociateTrafficFilterAssociateOperationOptions, callback func() error) error {
	result, err := c.AssociateTrafficFilterAssociate(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing AssociateTrafficFilterAssociate: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after AssociateTrafficFilterAssociate: %+v", err)
	}

	return nil
}
