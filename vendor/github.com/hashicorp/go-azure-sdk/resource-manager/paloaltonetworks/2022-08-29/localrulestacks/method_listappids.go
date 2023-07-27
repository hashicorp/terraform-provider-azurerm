package localrulestacks

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListAppIdsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *ListAppIdResponse
}

type ListAppIdsOperationOptions struct {
	AppIdVersion *string
	AppPrefix    *string
	Skip         *string
	Top          *int64
}

func DefaultListAppIdsOperationOptions() ListAppIdsOperationOptions {
	return ListAppIdsOperationOptions{}
}

func (o ListAppIdsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListAppIdsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListAppIdsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.AppIdVersion != nil {
		out.Append("appIdVersion", fmt.Sprintf("%v", *o.AppIdVersion))
	}
	if o.AppPrefix != nil {
		out.Append("appPrefix", fmt.Sprintf("%v", *o.AppPrefix))
	}
	if o.Skip != nil {
		out.Append("skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// ListAppIds ...
func (c LocalRulestacksClient) ListAppIds(ctx context.Context, id LocalRulestackId, options ListAppIdsOperationOptions) (result ListAppIdsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		Path:          fmt.Sprintf("%s/listAppIds", id.ID()),
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

	if err = resp.Unmarshal(&result.Model); err != nil {
		return
	}

	return
}
