package serverdevopsaudit

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

type SettingsListByServerOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ServerDevOpsAuditingSettings
}

type SettingsListByServerCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ServerDevOpsAuditingSettings
}

type SettingsListByServerCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *SettingsListByServerCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// SettingsListByServer ...
func (c ServerDevOpsAuditClient) SettingsListByServer(ctx context.Context, id commonids.SqlServerId) (result SettingsListByServerOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &SettingsListByServerCustomPager{},
		Path:       fmt.Sprintf("%s/devOpsAuditingSettings", id.ID()),
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
		Values *[]ServerDevOpsAuditingSettings `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// SettingsListByServerComplete retrieves all the results into a single object
func (c ServerDevOpsAuditClient) SettingsListByServerComplete(ctx context.Context, id commonids.SqlServerId) (SettingsListByServerCompleteResult, error) {
	return c.SettingsListByServerCompleteMatchingPredicate(ctx, id, ServerDevOpsAuditingSettingsOperationPredicate{})
}

// SettingsListByServerCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ServerDevOpsAuditClient) SettingsListByServerCompleteMatchingPredicate(ctx context.Context, id commonids.SqlServerId, predicate ServerDevOpsAuditingSettingsOperationPredicate) (result SettingsListByServerCompleteResult, err error) {
	items := make([]ServerDevOpsAuditingSettings, 0)

	resp, err := c.SettingsListByServer(ctx, id)
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

	result = SettingsListByServerCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
