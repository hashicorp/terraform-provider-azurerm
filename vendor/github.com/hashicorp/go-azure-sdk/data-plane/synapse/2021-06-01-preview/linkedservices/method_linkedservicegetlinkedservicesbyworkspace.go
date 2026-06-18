package linkedservices

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinkedServiceGetLinkedServicesByWorkspaceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]LinkedServiceResource
}

type LinkedServiceGetLinkedServicesByWorkspaceCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []LinkedServiceResource
}

type LinkedServiceGetLinkedServicesByWorkspaceCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *LinkedServiceGetLinkedServicesByWorkspaceCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// LinkedServiceGetLinkedServicesByWorkspace ...
func (c LinkedServicesClient) LinkedServiceGetLinkedServicesByWorkspace(ctx context.Context) (result LinkedServiceGetLinkedServicesByWorkspaceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &LinkedServiceGetLinkedServicesByWorkspaceCustomPager{},
		Path:       "/linkedServices",
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
		Values *[]LinkedServiceResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// LinkedServiceGetLinkedServicesByWorkspaceComplete retrieves all the results into a single object
func (c LinkedServicesClient) LinkedServiceGetLinkedServicesByWorkspaceComplete(ctx context.Context) (LinkedServiceGetLinkedServicesByWorkspaceCompleteResult, error) {
	return c.LinkedServiceGetLinkedServicesByWorkspaceCompleteMatchingPredicate(ctx, LinkedServiceResourceOperationPredicate{})
}

// LinkedServiceGetLinkedServicesByWorkspaceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c LinkedServicesClient) LinkedServiceGetLinkedServicesByWorkspaceCompleteMatchingPredicate(ctx context.Context, predicate LinkedServiceResourceOperationPredicate) (result LinkedServiceGetLinkedServicesByWorkspaceCompleteResult, err error) {
	items := make([]LinkedServiceResource, 0)

	resp, err := c.LinkedServiceGetLinkedServicesByWorkspace(ctx)
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

	result = LinkedServiceGetLinkedServicesByWorkspaceCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
