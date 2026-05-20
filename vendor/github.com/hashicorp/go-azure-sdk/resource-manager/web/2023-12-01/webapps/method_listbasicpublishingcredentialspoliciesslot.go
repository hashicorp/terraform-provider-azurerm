package webapps

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListBasicPublishingCredentialsPoliciesSlotOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CsmPublishingCredentialsPoliciesEntity
}

type ListBasicPublishingCredentialsPoliciesSlotCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []CsmPublishingCredentialsPoliciesEntity
}

type ListBasicPublishingCredentialsPoliciesSlotCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListBasicPublishingCredentialsPoliciesSlotCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListBasicPublishingCredentialsPoliciesSlot ...
func (c WebAppsClient) ListBasicPublishingCredentialsPoliciesSlot(ctx context.Context, id SlotId) (result ListBasicPublishingCredentialsPoliciesSlotOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListBasicPublishingCredentialsPoliciesSlotCustomPager{},
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

// ListBasicPublishingCredentialsPoliciesSlotComplete retrieves all the results into a single object
func (c WebAppsClient) ListBasicPublishingCredentialsPoliciesSlotComplete(ctx context.Context, id SlotId) (ListBasicPublishingCredentialsPoliciesSlotCompleteResult, error) {
	return c.ListBasicPublishingCredentialsPoliciesSlotCompleteMatchingPredicate(ctx, id, CsmPublishingCredentialsPoliciesEntityOperationPredicate{})
}

// ListBasicPublishingCredentialsPoliciesSlotCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListBasicPublishingCredentialsPoliciesSlotCompleteMatchingPredicate(ctx context.Context, id SlotId, predicate CsmPublishingCredentialsPoliciesEntityOperationPredicate) (result ListBasicPublishingCredentialsPoliciesSlotCompleteResult, err error) {
	items := make([]CsmPublishingCredentialsPoliciesEntity, 0)

	resp, err := c.ListBasicPublishingCredentialsPoliciesSlot(ctx, id)
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

	result = ListBasicPublishingCredentialsPoliciesSlotCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
