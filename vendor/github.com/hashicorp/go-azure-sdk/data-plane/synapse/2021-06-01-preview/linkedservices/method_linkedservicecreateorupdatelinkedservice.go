package linkedservices

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinkedServiceCreateOrUpdateLinkedServiceOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *LinkedServiceResource
}

type LinkedServiceCreateOrUpdateLinkedServiceOperationOptions struct {
	IfMatch *string
}

func DefaultLinkedServiceCreateOrUpdateLinkedServiceOperationOptions() LinkedServiceCreateOrUpdateLinkedServiceOperationOptions {
	return LinkedServiceCreateOrUpdateLinkedServiceOperationOptions{}
}

func (o LinkedServiceCreateOrUpdateLinkedServiceOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.IfMatch != nil {
		out.Append("If-Match", fmt.Sprintf("%v", *o.IfMatch))
	}
	return &out
}

func (o LinkedServiceCreateOrUpdateLinkedServiceOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o LinkedServiceCreateOrUpdateLinkedServiceOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// LinkedServiceCreateOrUpdateLinkedService ...
func (c LinkedServicesClient) LinkedServiceCreateOrUpdateLinkedService(ctx context.Context, id LinkedServiceId, input LinkedServiceResource, options LinkedServiceCreateOrUpdateLinkedServiceOperationOptions) (result LinkedServiceCreateOrUpdateLinkedServiceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod:    http.MethodPut,
		OptionsObject: options,
		Path:          id.Path(),
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

	result.Poller, err = dataplane.PollerFromResponse(resp, c.Client)
	if err != nil {
		return
	}

	return
}

// LinkedServiceCreateOrUpdateLinkedServiceThenPoll performs LinkedServiceCreateOrUpdateLinkedService then polls until it's completed
func (c LinkedServicesClient) LinkedServiceCreateOrUpdateLinkedServiceThenPoll(ctx context.Context, id LinkedServiceId, input LinkedServiceResource, options LinkedServiceCreateOrUpdateLinkedServiceOperationOptions) error {
	return c.LinkedServiceCreateOrUpdateLinkedServiceCallbackThenPoll(ctx, id, input, options, nil)
}

// LinkedServiceCreateOrUpdateLinkedServiceCallbackThenPoll performs LinkedServiceCreateOrUpdateLinkedService, runs the optional callback function, then polls until it's completed
func (c LinkedServicesClient) LinkedServiceCreateOrUpdateLinkedServiceCallbackThenPoll(ctx context.Context, id LinkedServiceId, input LinkedServiceResource, options LinkedServiceCreateOrUpdateLinkedServiceOperationOptions, callback func() error) error {
	result, err := c.LinkedServiceCreateOrUpdateLinkedService(ctx, id, input, options)
	if err != nil {
		return fmt.Errorf("performing LinkedServiceCreateOrUpdateLinkedService: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after LinkedServiceCreateOrUpdateLinkedService: %+v", err)
	}

	return nil
}
