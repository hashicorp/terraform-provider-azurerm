package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeleteOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
}

type DeleteOperationOptions struct {
	AppType             *AppType
	DeleteSubscriptions *bool
	IfMatch             *string
	Notify              *bool
}

func DefaultDeleteOperationOptions() DeleteOperationOptions {
	return DeleteOperationOptions{}
}

func (o DeleteOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.IfMatch != nil {
		out.Append("If-Match", fmt.Sprintf("%v", *o.IfMatch))
	}
	return &out
}

func (o DeleteOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o DeleteOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.AppType != nil {
		out.Append("appType", fmt.Sprintf("%v", *o.AppType))
	}
	if o.DeleteSubscriptions != nil {
		out.Append("deleteSubscriptions", fmt.Sprintf("%v", *o.DeleteSubscriptions))
	}
	if o.Notify != nil {
		out.Append("notify", fmt.Sprintf("%v", *o.Notify))
	}
	return &out
}

// Delete ...
func (c UserClient) Delete(ctx context.Context, id UserId, options DeleteOperationOptions) (result DeleteOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusNoContent,
			http.StatusOK,
		},
		HttpMethod:    http.MethodDelete,
		Path:          id.ID(),
		OptionsObject: options,
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
