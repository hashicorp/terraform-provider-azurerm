// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package quota

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	sdkclient "github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/resource-manager/quota/2025-07-15/groupquotalimits"
	"github.com/hashicorp/go-azure-sdk/resource-manager/quota/2025-07-15/subscriptionquotaallocation"
)

// listAllGroupQuotaLimits retrieves all quota limits for a (managementGroup, groupQuota,
// resourceProvider, location) scope, following NextLink pagination.
func listAllGroupQuotaLimits(ctx context.Context, c *groupquotalimits.GroupQuotaLimitsClient, id groupquotalimits.GroupQuotaLimitId) ([]groupquotalimits.GroupQuotaLimit, *http.Response, error) {
	resp, err := c.List(ctx, id)
	if err != nil {
		return nil, resp.HttpResponse, err
	}
	if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.Value == nil {
		return []groupquotalimits.GroupQuotaLimit{}, resp.HttpResponse, nil
	}

	items := make([]groupquotalimits.GroupQuotaLimit, len(*resp.Model.Properties.Value))
	copy(items, *resp.Model.Properties.Value)
	latestHTTPResp := resp.HttpResponse

	nextLink := resp.Model.Properties.NextLink
	for nextLink != nil && *nextLink != "" {
		nextURL, err := url.Parse(*nextLink)
		if err != nil {
			return items, latestHTTPResp, fmt.Errorf("parsing nextLink %q: %+v", *nextLink, err)
		}
		if !nextURL.IsAbs() {
			break
		}

		req, err := c.Client.NewRequest(ctx, sdkclient.RequestOptions{
			ContentType:         "application/json; charset=utf-8",
			ExpectedStatusCodes: []int{http.StatusOK},
			HttpMethod:          http.MethodGet,
			Path:                id.ID(),
		})
		if err != nil {
			return items, latestHTTPResp, fmt.Errorf("building next-page request: %+v", err)
		}
		req.URL = nextURL

		pageResp, err := req.Execute(ctx)
		if pageResp != nil {
			latestHTTPResp = pageResp.Response
		}
		if err != nil {
			return items, latestHTTPResp, fmt.Errorf("fetching next page: %+v", err)
		}

		var page groupquotalimits.GroupQuotaLimitList
		if err = pageResp.Unmarshal(&page); err != nil {
			return items, latestHTTPResp, fmt.Errorf("unmarshalling next-page response: %+v", err)
		}

		if page.Properties != nil && page.Properties.Value != nil {
			items = append(items, *page.Properties.Value...)
			nextLink = page.Properties.NextLink
		} else {
			break
		}
	}

	return items, latestHTTPResp, nil
}

// listAllSubscriptionAllocations retrieves all quota allocations for a
// (managementGroup, subscription, groupQuota, resourceProvider, location) scope,
// following NextLink pagination.
func listAllSubscriptionAllocations(ctx context.Context, c *subscriptionquotaallocation.SubscriptionQuotaAllocationClient, id subscriptionquotaallocation.QuotaAllocationId) ([]subscriptionquotaallocation.SubscriptionQuotaAllocations, *http.Response, error) {
	resp, err := c.GroupQuotaSubscriptionAllocationList(ctx, id)
	if err != nil {
		return nil, resp.HttpResponse, err
	}
	if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.Value == nil {
		return []subscriptionquotaallocation.SubscriptionQuotaAllocations{}, resp.HttpResponse, nil
	}

	items := make([]subscriptionquotaallocation.SubscriptionQuotaAllocations, len(*resp.Model.Properties.Value))
	copy(items, *resp.Model.Properties.Value)
	latestHTTPResp := resp.HttpResponse

	nextLink := resp.Model.Properties.NextLink
	for nextLink != nil && *nextLink != "" {
		nextURL, err := url.Parse(*nextLink)
		if err != nil {
			return items, latestHTTPResp, fmt.Errorf("parsing nextLink %q: %+v", *nextLink, err)
		}
		if !nextURL.IsAbs() {
			break
		}

		req, err := c.Client.NewRequest(ctx, sdkclient.RequestOptions{
			ContentType:         "application/json; charset=utf-8",
			ExpectedStatusCodes: []int{http.StatusOK},
			HttpMethod:          http.MethodGet,
			Path:                id.ID(),
		})
		if err != nil {
			return items, latestHTTPResp, fmt.Errorf("building next-page request: %+v", err)
		}
		req.URL = nextURL

		pageResp, err := req.Execute(ctx)
		if pageResp != nil {
			latestHTTPResp = pageResp.Response
		}
		if err != nil {
			return items, latestHTTPResp, fmt.Errorf("fetching next page: %+v", err)
		}

		var page subscriptionquotaallocation.SubscriptionQuotaAllocationsList
		if err = pageResp.Unmarshal(&page); err != nil {
			return items, latestHTTPResp, fmt.Errorf("unmarshalling next-page response: %+v", err)
		}

		if page.Properties != nil && page.Properties.Value != nil {
			items = append(items, *page.Properties.Value...)
			nextLink = page.Properties.NextLink
		} else {
			break
		}
	}

	return items, latestHTTPResp, nil
}
