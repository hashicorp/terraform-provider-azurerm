package staticsites

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListStaticSiteUsersOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]StaticSiteUserARMResource
}

type ListStaticSiteUsersCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []StaticSiteUserARMResource
}

// ListStaticSiteUsers ...
func (c StaticSitesClient) ListStaticSiteUsers(ctx context.Context, id AuthProviderId) (result ListStaticSiteUsersOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/listUsers", id.ID()),
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
		Values *[]StaticSiteUserARMResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListStaticSiteUsersComplete retrieves all the results into a single object
func (c StaticSitesClient) ListStaticSiteUsersComplete(ctx context.Context, id AuthProviderId) (ListStaticSiteUsersCompleteResult, error) {
	return c.ListStaticSiteUsersCompleteMatchingPredicate(ctx, id, StaticSiteUserARMResourceOperationPredicate{})
}

// ListStaticSiteUsersCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c StaticSitesClient) ListStaticSiteUsersCompleteMatchingPredicate(ctx context.Context, id AuthProviderId, predicate StaticSiteUserARMResourceOperationPredicate) (result ListStaticSiteUsersCompleteResult, err error) {
	items := make([]StaticSiteUserARMResource, 0)

	resp, err := c.ListStaticSiteUsers(ctx, id)
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

	result = ListStaticSiteUsersCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
