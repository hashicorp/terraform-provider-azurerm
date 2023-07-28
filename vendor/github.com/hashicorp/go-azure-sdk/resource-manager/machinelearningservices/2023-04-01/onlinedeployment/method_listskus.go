package onlinedeployment

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListSkusOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SkuResource
}

type ListSkusCompleteResult struct {
	Items []SkuResource
}

type ListSkusOperationOptions struct {
	Count *int64
	Skip  *string
}

func DefaultListSkusOperationOptions() ListSkusOperationOptions {
	return ListSkusOperationOptions{}
}

func (o ListSkusOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListSkusOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListSkusOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Count != nil {
		out.Append("count", fmt.Sprintf("%v", *o.Count))
	}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	return &out
}

// ListSkus ...
func (c OnlineDeploymentClient) ListSkus(ctx context.Context, id OnlineEndpointDeploymentId, options ListSkusOperationOptions) (result ListSkusOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/skus", id.ID()),
		OptionsObject: options,
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.ExecutePaged(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	var values struct {
		Values *[]SkuResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListSkusComplete retrieves all the results into a single object
func (c OnlineDeploymentClient) ListSkusComplete(ctx context.Context, id OnlineEndpointDeploymentId, options ListSkusOperationOptions) (ListSkusCompleteResult, error) {
	return c.ListSkusCompleteMatchingPredicate(ctx, id, options, SkuResourceOperationPredicate{})
}

// ListSkusCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OnlineDeploymentClient) ListSkusCompleteMatchingPredicate(ctx context.Context, id OnlineEndpointDeploymentId, options ListSkusOperationOptions, predicate SkuResourceOperationPredicate) (result ListSkusCompleteResult, err error) {
	items := make([]SkuResource, 0)

	resp, err := c.ListSkus(ctx, id, options)
	if err != nil {
		err = fmt.Errorf("loading results: %+v", err)
		return
	}
	if resp.Model != nil {
		for _, v := range *resp.Model {
			if predicate.Matches(v) {
				items = append(items, v)
			}
		}
	}

	result = ListSkusCompleteResult{
		Items: items,
	}
	return
}
