package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CheckNameAvailabilityOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *CheckNameAvailabilityOutput
}

type CheckNameAvailabilityOperationOptions struct {
	XMsClientRequestId *string
}

func DefaultCheckNameAvailabilityOperationOptions() CheckNameAvailabilityOperationOptions {
	return CheckNameAvailabilityOperationOptions{}
}

func (o CheckNameAvailabilityOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.XMsClientRequestId != nil {
		out.Append("x-ms-client-request-id", fmt.Sprintf("%v", *o.XMsClientRequestId))
	}
	return &out
}

func (o CheckNameAvailabilityOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o CheckNameAvailabilityOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// CheckNameAvailability ...
func (c ServicesClient) CheckNameAvailability(ctx context.Context, id commonids.SubscriptionId, input CheckNameAvailabilityInput, options CheckNameAvailabilityOperationOptions) (result CheckNameAvailabilityOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		Path:          fmt.Sprintf("%s/providers/Microsoft.Search/checkNameAvailability", id.ID()),
		OptionsObject: options,
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

	var model CheckNameAvailabilityOutput
	result.Model = &model

	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
