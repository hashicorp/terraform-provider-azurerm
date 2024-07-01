package providers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProviderPermissionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ProviderPermission
}

type ProviderPermissionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ProviderPermission
}

type ProviderPermissionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ProviderPermissionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ProviderPermissions ...
func (c ProvidersClient) ProviderPermissions(ctx context.Context, id SubscriptionProviderId) (result ProviderPermissionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ProviderPermissionsCustomPager{},
		Path:       fmt.Sprintf("%s/providerPermissions", id.ID()),
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
		Values *[]ProviderPermission `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ProviderPermissionsComplete retrieves all the results into a single object
func (c ProvidersClient) ProviderPermissionsComplete(ctx context.Context, id SubscriptionProviderId) (ProviderPermissionsCompleteResult, error) {
	return c.ProviderPermissionsCompleteMatchingPredicate(ctx, id, ProviderPermissionOperationPredicate{})
}

// ProviderPermissionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ProvidersClient) ProviderPermissionsCompleteMatchingPredicate(ctx context.Context, id SubscriptionProviderId, predicate ProviderPermissionOperationPredicate) (result ProviderPermissionsCompleteResult, err error) {
	items := make([]ProviderPermission, 0)

	resp, err := c.ProviderPermissions(ctx, id)
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

	result = ProviderPermissionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
