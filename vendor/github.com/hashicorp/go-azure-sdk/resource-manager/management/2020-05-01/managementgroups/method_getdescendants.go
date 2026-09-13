package managementgroups

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

type GetDescendantsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DescendantInfo
}

type GetDescendantsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DescendantInfo
}

type GetDescendantsOperationOptions struct {
	Top *int64
}

func DefaultGetDescendantsOperationOptions() GetDescendantsOperationOptions {
	return GetDescendantsOperationOptions{}
}

func (o GetDescendantsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o GetDescendantsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o GetDescendantsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type GetDescendantsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GetDescendantsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GetDescendants ...
func (c ManagementGroupsClient) GetDescendants(ctx context.Context, id commonids.ManagementGroupId, options GetDescendantsOperationOptions) (result GetDescendantsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &GetDescendantsCustomPager{},
		Path:          fmt.Sprintf("%s/descendants", id.ID()),
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
		Values *[]DescendantInfo `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GetDescendantsComplete retrieves all the results into a single object
func (c ManagementGroupsClient) GetDescendantsComplete(ctx context.Context, id commonids.ManagementGroupId, options GetDescendantsOperationOptions) (GetDescendantsCompleteResult, error) {
	return c.GetDescendantsCompleteMatchingPredicate(ctx, id, options, DescendantInfoOperationPredicate{})
}

// GetDescendantsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ManagementGroupsClient) GetDescendantsCompleteMatchingPredicate(ctx context.Context, id commonids.ManagementGroupId, options GetDescendantsOperationOptions, predicate DescendantInfoOperationPredicate) (result GetDescendantsCompleteResult, err error) {
	items := make([]DescendantInfo, 0)

	resp, err := c.GetDescendants(ctx, id, options)
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

	result = GetDescendantsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
