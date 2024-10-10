package servicelinker

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

type LinkerListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]LinkerResource
}

type LinkerListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []LinkerResource
}

type LinkerListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *LinkerListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// LinkerList ...
func (c ServicelinkerClient) LinkerList(ctx context.Context, id commonids.ScopeId) (result LinkerListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &LinkerListCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.ServiceLinker/linkers", id.ID()),
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
		Values *[]LinkerResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// LinkerListComplete retrieves all the results into a single object
func (c ServicelinkerClient) LinkerListComplete(ctx context.Context, id commonids.ScopeId) (LinkerListCompleteResult, error) {
	return c.LinkerListCompleteMatchingPredicate(ctx, id, LinkerResourceOperationPredicate{})
}

// LinkerListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ServicelinkerClient) LinkerListCompleteMatchingPredicate(ctx context.Context, id commonids.ScopeId, predicate LinkerResourceOperationPredicate) (result LinkerListCompleteResult, err error) {
	items := make([]LinkerResource, 0)

	resp, err := c.LinkerList(ctx, id)
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

	result = LinkerListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
