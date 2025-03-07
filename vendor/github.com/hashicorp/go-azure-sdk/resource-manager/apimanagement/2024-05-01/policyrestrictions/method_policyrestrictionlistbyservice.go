package policyrestrictions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PolicyRestrictionListByServiceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PolicyRestrictionContract
}

type PolicyRestrictionListByServiceCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []PolicyRestrictionContract
}

type PolicyRestrictionListByServiceCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *PolicyRestrictionListByServiceCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// PolicyRestrictionListByService ...
func (c PolicyRestrictionsClient) PolicyRestrictionListByService(ctx context.Context, id ServiceId) (result PolicyRestrictionListByServiceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &PolicyRestrictionListByServiceCustomPager{},
		Path:       fmt.Sprintf("%s/policyRestrictions", id.ID()),
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
		Values *[]PolicyRestrictionContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// PolicyRestrictionListByServiceComplete retrieves all the results into a single object
func (c PolicyRestrictionsClient) PolicyRestrictionListByServiceComplete(ctx context.Context, id ServiceId) (PolicyRestrictionListByServiceCompleteResult, error) {
	return c.PolicyRestrictionListByServiceCompleteMatchingPredicate(ctx, id, PolicyRestrictionContractOperationPredicate{})
}

// PolicyRestrictionListByServiceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c PolicyRestrictionsClient) PolicyRestrictionListByServiceCompleteMatchingPredicate(ctx context.Context, id ServiceId, predicate PolicyRestrictionContractOperationPredicate) (result PolicyRestrictionListByServiceCompleteResult, err error) {
	items := make([]PolicyRestrictionContract, 0)

	resp, err := c.PolicyRestrictionListByService(ctx, id)
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

	result = PolicyRestrictionListByServiceCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
