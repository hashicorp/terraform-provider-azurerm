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

type CreateAndAssociateIPFilterCreateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

type CreateAndAssociateIPFilterCreateOperationOptions struct {
	IPs  *string
	Name *string
}

func DefaultCreateAndAssociateIPFilterCreateOperationOptions() CreateAndAssociateIPFilterCreateOperationOptions {
	return CreateAndAssociateIPFilterCreateOperationOptions{}
}

func (o CreateAndAssociateIPFilterCreateOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o CreateAndAssociateIPFilterCreateOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o CreateAndAssociateIPFilterCreateOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.IPs != nil {
		out.Append("ips", fmt.Sprintf("%v", *o.IPs))
	}
	if o.Name != nil {
		out.Append("name", fmt.Sprintf("%v", *o.Name))
	}
	return &out
}

// CreateAndAssociateIPFilterCreate ...
func (c ElasticMonitorResourcesClient) CreateAndAssociateIPFilterCreate(ctx context.Context, id MonitorId, options CreateAndAssociateIPFilterCreateOperationOptions) (result CreateAndAssociateIPFilterCreateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/createAndAssociateIPFilter", id.ID()),
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

// CreateAndAssociateIPFilterCreateThenPoll performs CreateAndAssociateIPFilterCreate then polls until it's completed
func (c ElasticMonitorResourcesClient) CreateAndAssociateIPFilterCreateThenPoll(ctx context.Context, id MonitorId, options CreateAndAssociateIPFilterCreateOperationOptions) error {
	return c.CreateAndAssociateIPFilterCreateCallbackThenPoll(ctx, id, options, nil)
}

// CreateAndAssociateIPFilterCreateCallbackThenPoll performs CreateAndAssociateIPFilterCreate, runs the optional callback function, then polls until it's completed
func (c ElasticMonitorResourcesClient) CreateAndAssociateIPFilterCreateCallbackThenPoll(ctx context.Context, id MonitorId, options CreateAndAssociateIPFilterCreateOperationOptions, callback func() error) error {
	result, err := c.CreateAndAssociateIPFilterCreate(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing CreateAndAssociateIPFilterCreate: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after CreateAndAssociateIPFilterCreate: %+v", err)
	}

	return nil
}
