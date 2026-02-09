package apioperationpolicy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceApiOperationPolicyCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *PolicyContract
}

type WorkspaceApiOperationPolicyCreateOrUpdateOperationOptions struct {
	IfMatch *string
}

func DefaultWorkspaceApiOperationPolicyCreateOrUpdateOperationOptions() WorkspaceApiOperationPolicyCreateOrUpdateOperationOptions {
	return WorkspaceApiOperationPolicyCreateOrUpdateOperationOptions{}
}

func (o WorkspaceApiOperationPolicyCreateOrUpdateOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.IfMatch != nil {
		out.Append("If-Match", fmt.Sprintf("%v", *o.IfMatch))
	}
	return &out
}

func (o WorkspaceApiOperationPolicyCreateOrUpdateOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o WorkspaceApiOperationPolicyCreateOrUpdateOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// WorkspaceApiOperationPolicyCreateOrUpdate ...
func (c ApiOperationPolicyClient) WorkspaceApiOperationPolicyCreateOrUpdate(ctx context.Context, id ApiOperationId, input PolicyContract, options WorkspaceApiOperationPolicyCreateOrUpdateOperationOptions) (result WorkspaceApiOperationPolicyCreateOrUpdateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusCreated,
			http.StatusOK,
		},
		HttpMethod:    http.MethodPut,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/policies/policy", id.ID()),
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

	var model PolicyContract
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
