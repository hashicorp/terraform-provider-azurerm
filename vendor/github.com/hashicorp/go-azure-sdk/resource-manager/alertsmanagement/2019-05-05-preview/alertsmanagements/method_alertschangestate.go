package alertsmanagements

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertsChangeStateOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *Alert
}

type AlertsChangeStateOperationOptions struct {
	NewState *AlertState
}

func DefaultAlertsChangeStateOperationOptions() AlertsChangeStateOperationOptions {
	return AlertsChangeStateOperationOptions{}
}

func (o AlertsChangeStateOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o AlertsChangeStateOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o AlertsChangeStateOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.NewState != nil {
		out.Append("newState", fmt.Sprintf("%v", *o.NewState))
	}
	return &out
}

// AlertsChangeState ...
func (c AlertsManagementsClient) AlertsChangeState(ctx context.Context, id AlertId, input Comments, options AlertsChangeStateOperationOptions) (result AlertsChangeStateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		Path:          fmt.Sprintf("%s/changestate", id.ID()),
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

	var model Alert
	result.Model = &model

	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
