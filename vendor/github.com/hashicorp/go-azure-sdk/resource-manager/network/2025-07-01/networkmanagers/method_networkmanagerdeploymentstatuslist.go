package networkmanagers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkManagerDeploymentStatusListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *NetworkManagerDeploymentStatusListResult
}

type NetworkManagerDeploymentStatusListOperationOptions struct {
	Top *int64
}

func DefaultNetworkManagerDeploymentStatusListOperationOptions() NetworkManagerDeploymentStatusListOperationOptions {
	return NetworkManagerDeploymentStatusListOperationOptions{}
}

func (o NetworkManagerDeploymentStatusListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o NetworkManagerDeploymentStatusListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o NetworkManagerDeploymentStatusListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// NetworkManagerDeploymentStatusList ...
func (c NetworkManagersClient) NetworkManagerDeploymentStatusList(ctx context.Context, id NetworkManagerId, input NetworkManagerDeploymentStatusParameter, options NetworkManagerDeploymentStatusListOperationOptions) (result NetworkManagerDeploymentStatusListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/listDeploymentStatus", id.ID()),
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

	var model NetworkManagerDeploymentStatusListResult
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
