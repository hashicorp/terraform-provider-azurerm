package ehnamespaces

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

type NamespacesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]EHNamespace
}

type NamespacesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []EHNamespace
}

type NamespacesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *NamespacesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// NamespacesList ...
func (c EHNamespacesClient) NamespacesList(ctx context.Context, id commonids.SubscriptionId) (result NamespacesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &NamespacesListCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.EventHub/namespaces", id.ID()),
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
		Values *[]EHNamespace `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// NamespacesListComplete retrieves all the results into a single object
func (c EHNamespacesClient) NamespacesListComplete(ctx context.Context, id commonids.SubscriptionId) (NamespacesListCompleteResult, error) {
	return c.NamespacesListCompleteMatchingPredicate(ctx, id, EHNamespaceOperationPredicate{})
}

// NamespacesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c EHNamespacesClient) NamespacesListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate EHNamespaceOperationPredicate) (result NamespacesListCompleteResult, err error) {
	items := make([]EHNamespace, 0)

	resp, err := c.NamespacesList(ctx, id)
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

	result = NamespacesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
