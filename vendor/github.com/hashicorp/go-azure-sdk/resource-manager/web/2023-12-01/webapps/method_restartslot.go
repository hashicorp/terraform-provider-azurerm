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

type RestartSlotOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
}

type RestartSlotOperationOptions struct {
	SoftRestart *bool
	Synchronous *bool
}

func DefaultRestartSlotOperationOptions() RestartSlotOperationOptions {
	return RestartSlotOperationOptions{}
}

func (o RestartSlotOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RestartSlotOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o RestartSlotOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.SoftRestart != nil {
		out.Append("softRestart", fmt.Sprintf("%v", *o.SoftRestart))
	}
	if o.Synchronous != nil {
		out.Append("synchronous", fmt.Sprintf("%v", *o.Synchronous))
	}
	return &out
}

// RestartSlot ...
func (c WebAppsClient) RestartSlot(ctx context.Context, id SlotId, options RestartSlotOperationOptions) (result RestartSlotOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/restart", id.ID()),
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

	return
}
