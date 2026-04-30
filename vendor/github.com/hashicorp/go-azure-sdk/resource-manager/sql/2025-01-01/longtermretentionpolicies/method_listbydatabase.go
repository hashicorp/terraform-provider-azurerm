package longtermretentionpolicies

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

type ListByDatabaseOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]LongTermRetentionPolicy
}

type ListByDatabaseCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []LongTermRetentionPolicy
}

type ListByDatabaseCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByDatabaseCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByDatabase ...
func (c LongTermRetentionPoliciesClient) ListByDatabase(ctx context.Context, id commonids.SqlDatabaseId) (result ListByDatabaseOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByDatabaseCustomPager{},
		Path:       fmt.Sprintf("%s/backupLongTermRetentionPolicies", id.ID()),
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
		Values *[]LongTermRetentionPolicy `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByDatabaseComplete retrieves all the results into a single object
func (c LongTermRetentionPoliciesClient) ListByDatabaseComplete(ctx context.Context, id commonids.SqlDatabaseId) (ListByDatabaseCompleteResult, error) {
	return c.ListByDatabaseCompleteMatchingPredicate(ctx, id, LongTermRetentionPolicyOperationPredicate{})
}

// ListByDatabaseCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c LongTermRetentionPoliciesClient) ListByDatabaseCompleteMatchingPredicate(ctx context.Context, id commonids.SqlDatabaseId, predicate LongTermRetentionPolicyOperationPredicate) (result ListByDatabaseCompleteResult, err error) {
	items := make([]LongTermRetentionPolicy, 0)

	resp, err := c.ListByDatabase(ctx, id)
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

	result = ListByDatabaseCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
