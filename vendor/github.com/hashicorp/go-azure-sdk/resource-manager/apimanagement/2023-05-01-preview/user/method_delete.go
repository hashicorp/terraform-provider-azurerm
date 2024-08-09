package user

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

type DeleteOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

type DeleteOperationOptions struct {
	AppType             *AppType
	DeleteSubscriptions *bool
	IfMatch             *string
	Notify              *bool
}

func DefaultDeleteOperationOptions() DeleteOperationOptions {
	return DeleteOperationOptions{}
}

func (o DeleteOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.IfMatch != nil {
		out.Append("If-Match", fmt.Sprintf("%v", *o.IfMatch))
	}
	return &out
}

func (o DeleteOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o DeleteOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.AppType != nil {
		out.Append("appType", fmt.Sprintf("%v", *o.AppType))
	}
	if o.DeleteSubscriptions != nil {
		out.Append("deleteSubscriptions", fmt.Sprintf("%v", *o.DeleteSubscriptions))
	}
	if o.Notify != nil {
		out.Append("notify", fmt.Sprintf("%v", *o.Notify))
	}
	return &out
}

// Delete ...
func (c UserClient) Delete(ctx context.Context, id UserId, options DeleteOperationOptions) (result DeleteOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusNoContent,
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

// DeleteThenPoll performs Delete then polls until it's completed
func (c UserClient) DeleteThenPoll(ctx context.Context, id UserId, options DeleteOperationOptions) error {
	result, err := c.Delete(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing Delete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after Delete: %+v", err)
	}

	return nil
}
