package expressroutelagmembers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteLagsMembersListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ExpressRouteLagMember
}

type ExpressRouteLagsMembersListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ExpressRouteLagMember
}

type ExpressRouteLagsMembersListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ExpressRouteLagsMembersListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ExpressRouteLagsMembersList ...
func (c ExpressRouteLagMembersClient) ExpressRouteLagsMembersList(ctx context.Context, id LinkId) (result ExpressRouteLagsMembersListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ExpressRouteLagsMembersListCustomPager{},
		Path:       fmt.Sprintf("%s/members", id.ID()),
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
		Values *[]ExpressRouteLagMember `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ExpressRouteLagsMembersListComplete retrieves all the results into a single object
func (c ExpressRouteLagMembersClient) ExpressRouteLagsMembersListComplete(ctx context.Context, id LinkId) (ExpressRouteLagsMembersListCompleteResult, error) {
	return c.ExpressRouteLagsMembersListCompleteMatchingPredicate(ctx, id, ExpressRouteLagMemberOperationPredicate{})
}

// ExpressRouteLagsMembersListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ExpressRouteLagMembersClient) ExpressRouteLagsMembersListCompleteMatchingPredicate(ctx context.Context, id LinkId, predicate ExpressRouteLagMemberOperationPredicate) (result ExpressRouteLagsMembersListCompleteResult, err error) {
	items := make([]ExpressRouteLagMember, 0)

	resp, err := c.ExpressRouteLagsMembersList(ctx, id)
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

	result = ExpressRouteLagsMembersListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
