package providers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetAtTenantScopeOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *Provider
}

type GetAtTenantScopeOperationOptions struct {
	Expand *string
}

func DefaultGetAtTenantScopeOperationOptions() GetAtTenantScopeOperationOptions {
	return GetAtTenantScopeOperationOptions{}
}

func (o GetAtTenantScopeOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o GetAtTenantScopeOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o GetAtTenantScopeOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Expand != nil {
		out.Append("$expand", fmt.Sprintf("%v", *o.Expand))
	}
	return &out
}

// GetAtTenantScope ...
func (c ProvidersClient) GetAtTenantScope(ctx context.Context, id ProviderId, options GetAtTenantScopeOperationOptions) (result GetAtTenantScopeOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
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

	var model Provider
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
