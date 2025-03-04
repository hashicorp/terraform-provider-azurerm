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

type ListGeoRegionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]GeoRegion
}

type ListGeoRegionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []GeoRegion
}

type ListGeoRegionsOperationOptions struct {
	LinuxDynamicWorkersEnabled *bool
	LinuxWorkersEnabled        *bool
	Sku                        *SkuName
	XenonWorkersEnabled        *bool
}

func DefaultListGeoRegionsOperationOptions() ListGeoRegionsOperationOptions {
	return ListGeoRegionsOperationOptions{}
}

func (o ListGeoRegionsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListGeoRegionsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListGeoRegionsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.LinuxDynamicWorkersEnabled != nil {
		out.Append("linuxDynamicWorkersEnabled", fmt.Sprintf("%v", *o.LinuxDynamicWorkersEnabled))
	}
	if o.LinuxWorkersEnabled != nil {
		out.Append("linuxWorkersEnabled", fmt.Sprintf("%v", *o.LinuxWorkersEnabled))
	}
	if o.Sku != nil {
		out.Append("sku", fmt.Sprintf("%v", *o.Sku))
	}
	if o.XenonWorkersEnabled != nil {
		out.Append("xenonWorkersEnabled", fmt.Sprintf("%v", *o.XenonWorkersEnabled))
	}
	return &out
}

type ListGeoRegionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListGeoRegionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListGeoRegions ...
func (c ResourceProvidersClient) ListGeoRegions(ctx context.Context, id commonids.SubscriptionId, options ListGeoRegionsOperationOptions) (result ListGeoRegionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListGeoRegionsCustomPager{},
		Path:          fmt.Sprintf("%s/providers/Microsoft.Web/geoRegions", id.ID()),
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
		Values *[]GeoRegion `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListGeoRegionsComplete retrieves all the results into a single object
func (c ResourceProvidersClient) ListGeoRegionsComplete(ctx context.Context, id commonids.SubscriptionId, options ListGeoRegionsOperationOptions) (ListGeoRegionsCompleteResult, error) {
	return c.ListGeoRegionsCompleteMatchingPredicate(ctx, id, options, GeoRegionOperationPredicate{})
}

// ListGeoRegionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ResourceProvidersClient) ListGeoRegionsCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, options ListGeoRegionsOperationOptions, predicate GeoRegionOperationPredicate) (result ListGeoRegionsCompleteResult, err error) {
	items := make([]GeoRegion, 0)

	resp, err := c.ListGeoRegions(ctx, id, options)
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

	result = ListGeoRegionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
