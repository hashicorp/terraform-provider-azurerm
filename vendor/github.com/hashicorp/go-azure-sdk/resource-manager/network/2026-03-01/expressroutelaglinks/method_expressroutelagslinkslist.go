package expressroutelaglinks

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteLagsLinksListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ExpressRouteLagLink
}

type ExpressRouteLagsLinksListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ExpressRouteLagLink
}

type ExpressRouteLagsLinksListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ExpressRouteLagsLinksListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ExpressRouteLagsLinksList ...
func (c ExpressRouteLagLinksClient) ExpressRouteLagsLinksList(ctx context.Context, id ExpressRouteLagId) (result ExpressRouteLagsLinksListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ExpressRouteLagsLinksListCustomPager{},
		Path:       fmt.Sprintf("%s/links", id.ID()),
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
		Values *[]ExpressRouteLagLink `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ExpressRouteLagsLinksListComplete retrieves all the results into a single object
func (c ExpressRouteLagLinksClient) ExpressRouteLagsLinksListComplete(ctx context.Context, id ExpressRouteLagId) (ExpressRouteLagsLinksListCompleteResult, error) {
	return c.ExpressRouteLagsLinksListCompleteMatchingPredicate(ctx, id, ExpressRouteLagLinkOperationPredicate{})
}

// ExpressRouteLagsLinksListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ExpressRouteLagLinksClient) ExpressRouteLagsLinksListCompleteMatchingPredicate(ctx context.Context, id ExpressRouteLagId, predicate ExpressRouteLagLinkOperationPredicate) (result ExpressRouteLagsLinksListCompleteResult, err error) {
	items := make([]ExpressRouteLagLink, 0)

	resp, err := c.ExpressRouteLagsLinksList(ctx, id)
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

	result = ExpressRouteLagsLinksListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
