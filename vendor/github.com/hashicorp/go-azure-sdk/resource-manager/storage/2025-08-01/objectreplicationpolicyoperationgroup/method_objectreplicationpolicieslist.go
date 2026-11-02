package objectreplicationpolicyoperationgroup

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

type ObjectReplicationPoliciesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ObjectReplicationPolicy
}

type ObjectReplicationPoliciesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ObjectReplicationPolicy
}

type ObjectReplicationPoliciesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ObjectReplicationPoliciesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ObjectReplicationPoliciesList ...
func (c ObjectReplicationPolicyOperationGroupClient) ObjectReplicationPoliciesList(ctx context.Context, id commonids.StorageAccountId) (result ObjectReplicationPoliciesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ObjectReplicationPoliciesListCustomPager{},
		Path:       fmt.Sprintf("%s/objectReplicationPolicies", id.ID()),
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
		Values *[]ObjectReplicationPolicy `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ObjectReplicationPoliciesListComplete retrieves all the results into a single object
func (c ObjectReplicationPolicyOperationGroupClient) ObjectReplicationPoliciesListComplete(ctx context.Context, id commonids.StorageAccountId) (ObjectReplicationPoliciesListCompleteResult, error) {
	return c.ObjectReplicationPoliciesListCompleteMatchingPredicate(ctx, id, ObjectReplicationPolicyOperationPredicate{})
}

// ObjectReplicationPoliciesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ObjectReplicationPolicyOperationGroupClient) ObjectReplicationPoliciesListCompleteMatchingPredicate(ctx context.Context, id commonids.StorageAccountId, predicate ObjectReplicationPolicyOperationPredicate) (result ObjectReplicationPoliciesListCompleteResult, err error) {
	items := make([]ObjectReplicationPolicy, 0)

	resp, err := c.ObjectReplicationPoliciesList(ctx, id)
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

	result = ObjectReplicationPoliciesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
