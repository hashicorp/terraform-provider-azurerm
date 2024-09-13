package arcsettings

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ArcSettingsListByClusterOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ArcSetting
}

type ArcSettingsListByClusterCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ArcSetting
}

type ArcSettingsListByClusterCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ArcSettingsListByClusterCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ArcSettingsListByCluster ...
func (c ArcSettingsClient) ArcSettingsListByCluster(ctx context.Context, id ClusterId) (result ArcSettingsListByClusterOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ArcSettingsListByClusterCustomPager{},
		Path:       fmt.Sprintf("%s/arcSettings", id.ID()),
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
		Values *[]ArcSetting `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ArcSettingsListByClusterComplete retrieves all the results into a single object
func (c ArcSettingsClient) ArcSettingsListByClusterComplete(ctx context.Context, id ClusterId) (ArcSettingsListByClusterCompleteResult, error) {
	return c.ArcSettingsListByClusterCompleteMatchingPredicate(ctx, id, ArcSettingOperationPredicate{})
}

// ArcSettingsListByClusterCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ArcSettingsClient) ArcSettingsListByClusterCompleteMatchingPredicate(ctx context.Context, id ClusterId, predicate ArcSettingOperationPredicate) (result ArcSettingsListByClusterCompleteResult, err error) {
	items := make([]ArcSetting, 0)

	resp, err := c.ArcSettingsListByCluster(ctx, id)
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

	result = ArcSettingsListByClusterCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
