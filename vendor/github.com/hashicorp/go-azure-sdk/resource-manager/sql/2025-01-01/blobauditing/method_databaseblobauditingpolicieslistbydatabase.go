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

type DatabaseBlobAuditingPoliciesListByDatabaseOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DatabaseBlobAuditingPolicy
}

type DatabaseBlobAuditingPoliciesListByDatabaseCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DatabaseBlobAuditingPolicy
}

type DatabaseBlobAuditingPoliciesListByDatabaseCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *DatabaseBlobAuditingPoliciesListByDatabaseCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// DatabaseBlobAuditingPoliciesListByDatabase ...
func (c BlobAuditingClient) DatabaseBlobAuditingPoliciesListByDatabase(ctx context.Context, id commonids.SqlDatabaseId) (result DatabaseBlobAuditingPoliciesListByDatabaseOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &DatabaseBlobAuditingPoliciesListByDatabaseCustomPager{},
		Path:       fmt.Sprintf("%s/auditingSettings", id.ID()),
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
		Values *[]DatabaseBlobAuditingPolicy `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// DatabaseBlobAuditingPoliciesListByDatabaseComplete retrieves all the results into a single object
func (c BlobAuditingClient) DatabaseBlobAuditingPoliciesListByDatabaseComplete(ctx context.Context, id commonids.SqlDatabaseId) (DatabaseBlobAuditingPoliciesListByDatabaseCompleteResult, error) {
	return c.DatabaseBlobAuditingPoliciesListByDatabaseCompleteMatchingPredicate(ctx, id, DatabaseBlobAuditingPolicyOperationPredicate{})
}

// DatabaseBlobAuditingPoliciesListByDatabaseCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c BlobAuditingClient) DatabaseBlobAuditingPoliciesListByDatabaseCompleteMatchingPredicate(ctx context.Context, id commonids.SqlDatabaseId, predicate DatabaseBlobAuditingPolicyOperationPredicate) (result DatabaseBlobAuditingPoliciesListByDatabaseCompleteResult, err error) {
	items := make([]DatabaseBlobAuditingPolicy, 0)

	resp, err := c.DatabaseBlobAuditingPoliciesListByDatabase(ctx, id)
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

	result = DatabaseBlobAuditingPoliciesListByDatabaseCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
