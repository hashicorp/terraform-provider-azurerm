package clusters

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

type ListCalloutPoliciesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CalloutPolicy
}

type ListCalloutPoliciesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []CalloutPolicy
}

type ListCalloutPoliciesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListCalloutPoliciesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListCalloutPolicies ...
func (c ClustersClient) ListCalloutPolicies(ctx context.Context, id commonids.KustoClusterId) (result ListCalloutPoliciesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &ListCalloutPoliciesCustomPager{},
		Path:       fmt.Sprintf("%s/listCalloutPolicies", id.ID()),
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
		Values *[]CalloutPolicy `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListCalloutPoliciesComplete retrieves all the results into a single object
func (c ClustersClient) ListCalloutPoliciesComplete(ctx context.Context, id commonids.KustoClusterId) (ListCalloutPoliciesCompleteResult, error) {
	return c.ListCalloutPoliciesCompleteMatchingPredicate(ctx, id, CalloutPolicyOperationPredicate{})
}

// ListCalloutPoliciesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ClustersClient) ListCalloutPoliciesCompleteMatchingPredicate(ctx context.Context, id commonids.KustoClusterId, predicate CalloutPolicyOperationPredicate) (result ListCalloutPoliciesCompleteResult, err error) {
	items := make([]CalloutPolicy, 0)

	resp, err := c.ListCalloutPolicies(ctx, id)
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

	result = ListCalloutPoliciesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
