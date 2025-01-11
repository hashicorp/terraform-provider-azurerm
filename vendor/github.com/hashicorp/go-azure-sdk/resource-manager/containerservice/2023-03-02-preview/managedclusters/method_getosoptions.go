package managedclusters

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetOSOptionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *OSOptionProfile
}

type GetOSOptionsOperationOptions struct {
	ResourceType *string
}

func DefaultGetOSOptionsOperationOptions() GetOSOptionsOperationOptions {
	return GetOSOptionsOperationOptions{}
}

func (o GetOSOptionsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o GetOSOptionsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o GetOSOptionsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.ResourceType != nil {
		out.Append("resource-type", fmt.Sprintf("%v", *o.ResourceType))
	}
	return &out
}

// GetOSOptions ...
func (c ManagedClustersClient) GetOSOptions(ctx context.Context, id LocationId, options GetOSOptionsOperationOptions) (result GetOSOptionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/osOptions/default", id.ID()),
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

	var model OSOptionProfile
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
