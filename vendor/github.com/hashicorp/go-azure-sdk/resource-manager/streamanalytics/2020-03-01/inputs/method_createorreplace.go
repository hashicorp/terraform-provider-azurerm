package inputs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CreateOrReplaceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *Input
}

type CreateOrReplaceOperationOptions struct {
	IfMatch     *string
	IfNoneMatch *string
}

func DefaultCreateOrReplaceOperationOptions() CreateOrReplaceOperationOptions {
	return CreateOrReplaceOperationOptions{}
}

func (o CreateOrReplaceOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.IfMatch != nil {
		out.Append("If-Match", fmt.Sprintf("%v", *o.IfMatch))
	}
	if o.IfNoneMatch != nil {
		out.Append("If-None-Match", fmt.Sprintf("%v", *o.IfNoneMatch))
	}
	return &out
}

func (o CreateOrReplaceOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o CreateOrReplaceOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// CreateOrReplace ...
func (c InputsClient) CreateOrReplace(ctx context.Context, id InputId, input Input, options CreateOrReplaceOperationOptions) (result CreateOrReplaceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusCreated,
			http.StatusOK,
		},
		HttpMethod:    http.MethodPut,
		Path:          id.ID(),
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

	var model Input
	result.Model = &model

	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
