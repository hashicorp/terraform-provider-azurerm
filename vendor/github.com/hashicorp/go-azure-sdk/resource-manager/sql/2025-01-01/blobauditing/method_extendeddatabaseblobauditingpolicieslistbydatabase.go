package blobauditing

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

type ExtendedDatabaseBlobAuditingPoliciesListByDatabaseOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ExtendedDatabaseBlobAuditingPolicy
}

type ExtendedDatabaseBlobAuditingPoliciesListByDatabaseCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ExtendedDatabaseBlobAuditingPolicy
}

type ExtendedDatabaseBlobAuditingPoliciesListByDatabaseCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ExtendedDatabaseBlobAuditingPoliciesListByDatabaseCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ExtendedDatabaseBlobAuditingPoliciesListByDatabase ...
func (c BlobAuditingClient) ExtendedDatabaseBlobAuditingPoliciesListByDatabase(ctx context.Context, id commonids.SqlDatabaseId) (result ExtendedDatabaseBlobAuditingPoliciesListByDatabaseOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ExtendedDatabaseBlobAuditingPoliciesListByDatabaseCustomPager{},
		Path:       fmt.Sprintf("%s/extendedAuditingSettings", id.ID()),
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
		Values *[]ExtendedDatabaseBlobAuditingPolicy `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ExtendedDatabaseBlobAuditingPoliciesListByDatabaseComplete retrieves all the results into a single object
func (c BlobAuditingClient) ExtendedDatabaseBlobAuditingPoliciesListByDatabaseComplete(ctx context.Context, id commonids.SqlDatabaseId) (ExtendedDatabaseBlobAuditingPoliciesListByDatabaseCompleteResult, error) {
	return c.ExtendedDatabaseBlobAuditingPoliciesListByDatabaseCompleteMatchingPredicate(ctx, id, ExtendedDatabaseBlobAuditingPolicyOperationPredicate{})
}

// ExtendedDatabaseBlobAuditingPoliciesListByDatabaseCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c BlobAuditingClient) ExtendedDatabaseBlobAuditingPoliciesListByDatabaseCompleteMatchingPredicate(ctx context.Context, id commonids.SqlDatabaseId, predicate ExtendedDatabaseBlobAuditingPolicyOperationPredicate) (result ExtendedDatabaseBlobAuditingPoliciesListByDatabaseCompleteResult, err error) {
	items := make([]ExtendedDatabaseBlobAuditingPolicy, 0)

	resp, err := c.ExtendedDatabaseBlobAuditingPoliciesListByDatabase(ctx, id)
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

	result = ExtendedDatabaseBlobAuditingPoliciesListByDatabaseCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
