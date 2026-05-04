package appplatform

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BuildpackBindingListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]BuildpackBindingResource
}

type BuildpackBindingListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []BuildpackBindingResource
}

type BuildpackBindingListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *BuildpackBindingListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// BuildpackBindingList ...
func (c AppPlatformClient) BuildpackBindingList(ctx context.Context, id BuilderId) (result BuildpackBindingListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &BuildpackBindingListCustomPager{},
		Path:       fmt.Sprintf("%s/buildPackBindings", id.ID()),
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
		Values *[]BuildpackBindingResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// BuildpackBindingListComplete retrieves all the results into a single object
func (c AppPlatformClient) BuildpackBindingListComplete(ctx context.Context, id BuilderId) (BuildpackBindingListCompleteResult, error) {
	return c.BuildpackBindingListCompleteMatchingPredicate(ctx, id, BuildpackBindingResourceOperationPredicate{})
}

// BuildpackBindingListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppPlatformClient) BuildpackBindingListCompleteMatchingPredicate(ctx context.Context, id BuilderId, predicate BuildpackBindingResourceOperationPredicate) (result BuildpackBindingListCompleteResult, err error) {
	items := make([]BuildpackBindingResource, 0)

	resp, err := c.BuildpackBindingList(ctx, id)
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

	result = BuildpackBindingListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
