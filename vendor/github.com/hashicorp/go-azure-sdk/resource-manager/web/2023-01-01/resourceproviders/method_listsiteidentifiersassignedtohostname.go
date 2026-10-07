package resourceproviders

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

type ListSiteIdentifiersAssignedToHostNameOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Identifier
}

type ListSiteIdentifiersAssignedToHostNameCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Identifier
}

type ListSiteIdentifiersAssignedToHostNameCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListSiteIdentifiersAssignedToHostNameCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListSiteIdentifiersAssignedToHostName ...
func (c ResourceProvidersClient) ListSiteIdentifiersAssignedToHostName(ctx context.Context, id commonids.SubscriptionId, input NameIdentifier) (result ListSiteIdentifiersAssignedToHostNameOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &ListSiteIdentifiersAssignedToHostNameCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.Web/listSitesAssignedToHostName", id.ID()),
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
		Values *[]Identifier `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListSiteIdentifiersAssignedToHostNameComplete retrieves all the results into a single object
func (c ResourceProvidersClient) ListSiteIdentifiersAssignedToHostNameComplete(ctx context.Context, id commonids.SubscriptionId, input NameIdentifier) (ListSiteIdentifiersAssignedToHostNameCompleteResult, error) {
	return c.ListSiteIdentifiersAssignedToHostNameCompleteMatchingPredicate(ctx, id, input, IdentifierOperationPredicate{})
}

// ListSiteIdentifiersAssignedToHostNameCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ResourceProvidersClient) ListSiteIdentifiersAssignedToHostNameCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, input NameIdentifier, predicate IdentifierOperationPredicate) (result ListSiteIdentifiersAssignedToHostNameCompleteResult, err error) {
	items := make([]Identifier, 0)

	resp, err := c.ListSiteIdentifiersAssignedToHostName(ctx, id, input)
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

	result = ListSiteIdentifiersAssignedToHostNameCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
