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

type ListPredefinedUrlCategoriesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *PredefinedUrlCategoriesResponse
}

type ListPredefinedUrlCategoriesOperationOptions struct {
	Skip *string
	Top  *int64
}

func DefaultListPredefinedUrlCategoriesOperationOptions() ListPredefinedUrlCategoriesOperationOptions {
	return ListPredefinedUrlCategoriesOperationOptions{}
}

func (o ListPredefinedUrlCategoriesOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListPredefinedUrlCategoriesOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListPredefinedUrlCategoriesOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Skip != nil {
		out.Append("skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// ListPredefinedUrlCategories ...
func (c LocalRulestacksClient) ListPredefinedUrlCategories(ctx context.Context, id LocalRulestackId, options ListPredefinedUrlCategoriesOperationOptions) (result ListPredefinedUrlCategoriesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		Path:          fmt.Sprintf("%s/listPredefinedUrlCategories", id.ID()),
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
