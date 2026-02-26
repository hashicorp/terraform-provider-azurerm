package containerinstance

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CGProfileListAllRevisionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ContainerGroupProfile
}

type CGProfileListAllRevisionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ContainerGroupProfile
}

type CGProfileListAllRevisionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *CGProfileListAllRevisionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// CGProfileListAllRevisions ...
func (c ContainerInstanceClient) CGProfileListAllRevisions(ctx context.Context, id ContainerGroupProfileId) (result CGProfileListAllRevisionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &CGProfileListAllRevisionsCustomPager{},
		Path:       fmt.Sprintf("%s/revisions", id.ID()),
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
		Values *[]ContainerGroupProfile `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// CGProfileListAllRevisionsComplete retrieves all the results into a single object
func (c ContainerInstanceClient) CGProfileListAllRevisionsComplete(ctx context.Context, id ContainerGroupProfileId) (CGProfileListAllRevisionsCompleteResult, error) {
	return c.CGProfileListAllRevisionsCompleteMatchingPredicate(ctx, id, ContainerGroupProfileOperationPredicate{})
}

// CGProfileListAllRevisionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ContainerInstanceClient) CGProfileListAllRevisionsCompleteMatchingPredicate(ctx context.Context, id ContainerGroupProfileId, predicate ContainerGroupProfileOperationPredicate) (result CGProfileListAllRevisionsCompleteResult, err error) {
	items := make([]ContainerGroupProfile, 0)

	resp, err := c.CGProfileListAllRevisions(ctx, id)
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

	result = CGProfileListAllRevisionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
