package containerservices

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListOrchestratorsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *OrchestratorVersionProfileListResult
}

type ListOrchestratorsOperationOptions struct {
	ResourceType *string
}

func DefaultListOrchestratorsOperationOptions() ListOrchestratorsOperationOptions {
	return ListOrchestratorsOperationOptions{}
}

func (o ListOrchestratorsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListOrchestratorsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListOrchestratorsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.ResourceType != nil {
		out.Append("resource-type", fmt.Sprintf("%v", *o.ResourceType))
	}
	return &out
}

// ListOrchestrators ...
func (c ContainerServicesClient) ListOrchestrators(ctx context.Context, id LocationId, options ListOrchestratorsOperationOptions) (result ListOrchestratorsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/orchestrators", id.ID()),
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

	var model OrchestratorVersionProfileListResult
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
