package deploymentsettings

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByClustersOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DeploymentSetting
}

type ListByClustersCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DeploymentSetting
}

type ListByClustersCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByClustersCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByClusters ...
func (c DeploymentSettingsClient) ListByClusters(ctx context.Context, id ClusterId) (result ListByClustersOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByClustersCustomPager{},
		Path:       fmt.Sprintf("%s/deploymentSettings", id.ID()),
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
		Values *[]DeploymentSetting `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByClustersComplete retrieves all the results into a single object
func (c DeploymentSettingsClient) ListByClustersComplete(ctx context.Context, id ClusterId) (ListByClustersCompleteResult, error) {
	return c.ListByClustersCompleteMatchingPredicate(ctx, id, DeploymentSettingOperationPredicate{})
}

// ListByClustersCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DeploymentSettingsClient) ListByClustersCompleteMatchingPredicate(ctx context.Context, id ClusterId, predicate DeploymentSettingOperationPredicate) (result ListByClustersCompleteResult, err error) {
	items := make([]DeploymentSetting, 0)

	resp, err := c.ListByClusters(ctx, id)
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

	result = ListByClustersCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
