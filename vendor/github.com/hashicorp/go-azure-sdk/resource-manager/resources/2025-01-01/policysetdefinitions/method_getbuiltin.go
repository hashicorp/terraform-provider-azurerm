package policysetdefinitions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetBuiltInOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *PolicySetDefinition
}

type GetBuiltInOperationOptions struct {
	Expand *string
}

func DefaultGetBuiltInOperationOptions() GetBuiltInOperationOptions {
	return GetBuiltInOperationOptions{}
}

func (o GetBuiltInOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o GetBuiltInOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o GetBuiltInOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Expand != nil {
		out.Append("$expand", fmt.Sprintf("%v", *o.Expand))
	}
	return &out
}

// GetBuiltIn ...
func (c PolicySetDefinitionsClient) GetBuiltIn(ctx context.Context, id PolicySetDefinitionId, options GetBuiltInOperationOptions) (result GetBuiltInOperationResponse, err error) {
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

	var model PolicySetDefinition
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
