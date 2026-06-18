package linkedservices

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinkedServiceGetLinkedServiceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *LinkedServiceResource
}

type LinkedServiceGetLinkedServiceOperationOptions struct {
	IfNoneMatch *string
}

func DefaultLinkedServiceGetLinkedServiceOperationOptions() LinkedServiceGetLinkedServiceOperationOptions {
	return LinkedServiceGetLinkedServiceOperationOptions{}
}

func (o LinkedServiceGetLinkedServiceOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.IfNoneMatch != nil {
		out.Append("If-None-Match", fmt.Sprintf("%v", *o.IfNoneMatch))
	}
	return &out
}

func (o LinkedServiceGetLinkedServiceOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o LinkedServiceGetLinkedServiceOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// LinkedServiceGetLinkedService ...
func (c LinkedServicesClient) LinkedServiceGetLinkedService(ctx context.Context, id LinkedServiceId, options LinkedServiceGetLinkedServiceOperationOptions) (result LinkedServiceGetLinkedServiceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Path:          id.Path(),
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

	var model LinkedServiceResource
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
