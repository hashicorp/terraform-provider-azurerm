package machinelearningcomputes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComputeListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ComputeResource
}

type ComputeListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ComputeResource
}

type ComputeListOperationOptions struct {
	Skip *string
}

func DefaultComputeListOperationOptions() ComputeListOperationOptions {
	return ComputeListOperationOptions{}
}

func (o ComputeListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ComputeListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ComputeListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	return &out
}

type ComputeListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ComputeListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ComputeList ...
func (c MachineLearningComputesClient) ComputeList(ctx context.Context, id WorkspaceId, options ComputeListOperationOptions) (result ComputeListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ComputeListCustomPager{},
		Path:          fmt.Sprintf("%s/computes", id.ID()),
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
		Values *[]ComputeResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ComputeListComplete retrieves all the results into a single object
func (c MachineLearningComputesClient) ComputeListComplete(ctx context.Context, id WorkspaceId, options ComputeListOperationOptions) (ComputeListCompleteResult, error) {
	return c.ComputeListCompleteMatchingPredicate(ctx, id, options, ComputeResourceOperationPredicate{})
}

// ComputeListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c MachineLearningComputesClient) ComputeListCompleteMatchingPredicate(ctx context.Context, id WorkspaceId, options ComputeListOperationOptions, predicate ComputeResourceOperationPredicate) (result ComputeListCompleteResult, err error) {
	items := make([]ComputeResource, 0)

	resp, err := c.ComputeList(ctx, id, options)
	if err != nil {
		result.LatestHttpResponse = resp.HttpResponse
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

	result = ComputeListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
