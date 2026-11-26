package diagnosticsettingscategories

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

type DiagnosticSettingsCategoryListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DiagnosticSettingsCategoryResource
}

type DiagnosticSettingsCategoryListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DiagnosticSettingsCategoryResource
}

type DiagnosticSettingsCategoryListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *DiagnosticSettingsCategoryListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// DiagnosticSettingsCategoryList ...
func (c DiagnosticSettingsCategoriesClient) DiagnosticSettingsCategoryList(ctx context.Context, id commonids.ScopeId) (result DiagnosticSettingsCategoryListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &DiagnosticSettingsCategoryListCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.Insights/diagnosticSettingsCategories", id.ID()),
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
		Values *[]DiagnosticSettingsCategoryResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// DiagnosticSettingsCategoryListComplete retrieves all the results into a single object
func (c DiagnosticSettingsCategoriesClient) DiagnosticSettingsCategoryListComplete(ctx context.Context, id commonids.ScopeId) (DiagnosticSettingsCategoryListCompleteResult, error) {
	return c.DiagnosticSettingsCategoryListCompleteMatchingPredicate(ctx, id, DiagnosticSettingsCategoryResourceOperationPredicate{})
}

// DiagnosticSettingsCategoryListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DiagnosticSettingsCategoriesClient) DiagnosticSettingsCategoryListCompleteMatchingPredicate(ctx context.Context, id commonids.ScopeId, predicate DiagnosticSettingsCategoryResourceOperationPredicate) (result DiagnosticSettingsCategoryListCompleteResult, err error) {
	items := make([]DiagnosticSettingsCategoryResource, 0)

	resp, err := c.DiagnosticSettingsCategoryList(ctx, id)
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

	result = DiagnosticSettingsCategoryListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
