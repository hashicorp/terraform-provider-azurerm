package blobcontainers

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

type GetImmutabilityPolicyOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *ImmutabilityPolicy
}

type GetImmutabilityPolicyOperationOptions struct {
	IfMatch *string
}

func DefaultGetImmutabilityPolicyOperationOptions() GetImmutabilityPolicyOperationOptions {
	return GetImmutabilityPolicyOperationOptions{}
}

func (o GetImmutabilityPolicyOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.IfMatch != nil {
		out.Append("If-Match", fmt.Sprintf("%v", *o.IfMatch))
	}
	return &out
}

func (o GetImmutabilityPolicyOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o GetImmutabilityPolicyOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// GetImmutabilityPolicy ...
func (c BlobContainersClient) GetImmutabilityPolicy(ctx context.Context, id commonids.StorageContainerId, options GetImmutabilityPolicyOperationOptions) (result GetImmutabilityPolicyOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/immutabilityPolicies/default", id.ID()),
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

	var model ImmutabilityPolicy
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
