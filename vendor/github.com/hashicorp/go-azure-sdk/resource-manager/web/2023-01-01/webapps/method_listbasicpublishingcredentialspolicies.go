package webapps

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

type ListBasicPublishingCredentialsPoliciesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CsmPublishingCredentialsPoliciesEntity
}

type ListBasicPublishingCredentialsPoliciesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []CsmPublishingCredentialsPoliciesEntity
}

type ListBasicPublishingCredentialsPoliciesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListBasicPublishingCredentialsPoliciesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListBasicPublishingCredentialsPolicies ...
func (c WebAppsClient) ListBasicPublishingCredentialsPolicies(ctx context.Context, id commonids.AppServiceId) (result ListBasicPublishingCredentialsPoliciesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListBasicPublishingCredentialsPoliciesCustomPager{},
		Path:       fmt.Sprintf("%s/basicPublishingCredentialsPolicies", id.ID()),
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
		Values *[]CsmPublishingCredentialsPoliciesEntity `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListBasicPublishingCredentialsPoliciesComplete retrieves all the results into a single object
func (c WebAppsClient) ListBasicPublishingCredentialsPoliciesComplete(ctx context.Context, id commonids.AppServiceId) (ListBasicPublishingCredentialsPoliciesCompleteResult, error) {
	return c.ListBasicPublishingCredentialsPoliciesCompleteMatchingPredicate(ctx, id, CsmPublishingCredentialsPoliciesEntityOperationPredicate{})
}

// ListBasicPublishingCredentialsPoliciesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListBasicPublishingCredentialsPoliciesCompleteMatchingPredicate(ctx context.Context, id commonids.AppServiceId, predicate CsmPublishingCredentialsPoliciesEntityOperationPredicate) (result ListBasicPublishingCredentialsPoliciesCompleteResult, err error) {
	items := make([]CsmPublishingCredentialsPoliciesEntity, 0)

	resp, err := c.ListBasicPublishingCredentialsPolicies(ctx, id)
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

	result = ListBasicPublishingCredentialsPoliciesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
