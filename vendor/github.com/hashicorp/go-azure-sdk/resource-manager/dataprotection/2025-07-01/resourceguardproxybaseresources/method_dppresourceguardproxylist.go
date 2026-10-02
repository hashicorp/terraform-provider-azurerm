package resourceguardproxybaseresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DppResourceGuardProxyListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ResourceGuardProxyBaseResource
}

type DppResourceGuardProxyListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ResourceGuardProxyBaseResource
}

type DppResourceGuardProxyListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *DppResourceGuardProxyListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// DppResourceGuardProxyList ...
func (c ResourceGuardProxyBaseResourcesClient) DppResourceGuardProxyList(ctx context.Context, id BackupVaultId) (result DppResourceGuardProxyListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &DppResourceGuardProxyListCustomPager{},
		Path:       fmt.Sprintf("%s/backupResourceGuardProxies", id.ID()),
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
		Values *[]ResourceGuardProxyBaseResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// DppResourceGuardProxyListComplete retrieves all the results into a single object
func (c ResourceGuardProxyBaseResourcesClient) DppResourceGuardProxyListComplete(ctx context.Context, id BackupVaultId) (DppResourceGuardProxyListCompleteResult, error) {
	return c.DppResourceGuardProxyListCompleteMatchingPredicate(ctx, id, ResourceGuardProxyBaseResourceOperationPredicate{})
}

// DppResourceGuardProxyListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ResourceGuardProxyBaseResourcesClient) DppResourceGuardProxyListCompleteMatchingPredicate(ctx context.Context, id BackupVaultId, predicate ResourceGuardProxyBaseResourceOperationPredicate) (result DppResourceGuardProxyListCompleteResult, err error) {
	items := make([]ResourceGuardProxyBaseResource, 0)

	resp, err := c.DppResourceGuardProxyList(ctx, id)
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

	result = DppResourceGuardProxyListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
