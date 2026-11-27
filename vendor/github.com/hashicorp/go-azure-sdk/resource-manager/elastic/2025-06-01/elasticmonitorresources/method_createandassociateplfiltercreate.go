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

type CreateAndAssociatePLFilterCreateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

type CreateAndAssociatePLFilterCreateOperationOptions struct {
	Name                *string
	PrivateEndpointGuid *string
	PrivateEndpointName *string
}

func DefaultCreateAndAssociatePLFilterCreateOperationOptions() CreateAndAssociatePLFilterCreateOperationOptions {
	return CreateAndAssociatePLFilterCreateOperationOptions{}
}

func (o CreateAndAssociatePLFilterCreateOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o CreateAndAssociatePLFilterCreateOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o CreateAndAssociatePLFilterCreateOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Name != nil {
		out.Append("name", fmt.Sprintf("%v", *o.Name))
	}
	if o.PrivateEndpointGuid != nil {
		out.Append("privateEndpointGuid", fmt.Sprintf("%v", *o.PrivateEndpointGuid))
	}
	if o.PrivateEndpointName != nil {
		out.Append("privateEndpointName", fmt.Sprintf("%v", *o.PrivateEndpointName))
	}
	return &out
}

// CreateAndAssociatePLFilterCreate ...
func (c ElasticMonitorResourcesClient) CreateAndAssociatePLFilterCreate(ctx context.Context, id MonitorId, options CreateAndAssociatePLFilterCreateOperationOptions) (result CreateAndAssociatePLFilterCreateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/createAndAssociatePLFilter", id.ID()),
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

// CreateAndAssociatePLFilterCreateThenPoll performs CreateAndAssociatePLFilterCreate then polls until it's completed
func (c ElasticMonitorResourcesClient) CreateAndAssociatePLFilterCreateThenPoll(ctx context.Context, id MonitorId, options CreateAndAssociatePLFilterCreateOperationOptions) error {
	return c.CreateAndAssociatePLFilterCreateCallbackThenPoll(ctx, id, options, nil)
}

// CreateAndAssociatePLFilterCreateCallbackThenPoll performs CreateAndAssociatePLFilterCreate, runs the optional callback function, then polls until it's completed
func (c ElasticMonitorResourcesClient) CreateAndAssociatePLFilterCreateCallbackThenPoll(ctx context.Context, id MonitorId, options CreateAndAssociatePLFilterCreateOperationOptions, callback func() error) error {
	result, err := c.CreateAndAssociatePLFilterCreate(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing CreateAndAssociatePLFilterCreate: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after CreateAndAssociatePLFilterCreate: %+v", err)
	}

	return nil
}
