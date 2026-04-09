package appserviceplans

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

type RestartWebAppsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
}

type RestartWebAppsOperationOptions struct {
	SoftRestart *bool
}

func DefaultRestartWebAppsOperationOptions() RestartWebAppsOperationOptions {
	return RestartWebAppsOperationOptions{}
}

func (o RestartWebAppsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RestartWebAppsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o RestartWebAppsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.SoftRestart != nil {
		out.Append("softRestart", fmt.Sprintf("%v", *o.SoftRestart))
	}
	return &out
}

// RestartWebApps ...
func (c AppServicePlansClient) RestartWebApps(ctx context.Context, id commonids.AppServicePlanId, options RestartWebAppsOperationOptions) (result RestartWebAppsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusNoContent,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/restartSites", id.ID()),
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
