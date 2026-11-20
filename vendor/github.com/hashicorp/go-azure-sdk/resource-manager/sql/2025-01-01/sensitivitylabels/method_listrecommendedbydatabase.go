package sensitivitylabels

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

type ListRecommendedByDatabaseOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SensitivityLabel
}

type ListRecommendedByDatabaseCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SensitivityLabel
}

type ListRecommendedByDatabaseOperationOptions struct {
	Filter                         *string
	IncludeDisabledRecommendations *bool
}

func DefaultListRecommendedByDatabaseOperationOptions() ListRecommendedByDatabaseOperationOptions {
	return ListRecommendedByDatabaseOperationOptions{}
}

func (o ListRecommendedByDatabaseOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListRecommendedByDatabaseOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListRecommendedByDatabaseOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.IncludeDisabledRecommendations != nil {
		out.Append("includeDisabledRecommendations", fmt.Sprintf("%v", *o.IncludeDisabledRecommendations))
	}
	return &out
}

type ListRecommendedByDatabaseCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListRecommendedByDatabaseCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListRecommendedByDatabase ...
func (c SensitivityLabelsClient) ListRecommendedByDatabase(ctx context.Context, id commonids.SqlDatabaseId, options ListRecommendedByDatabaseOperationOptions) (result ListRecommendedByDatabaseOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListRecommendedByDatabaseCustomPager{},
		Path:          fmt.Sprintf("%s/recommendedSensitivityLabels", id.ID()),
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
		Values *[]SensitivityLabel `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListRecommendedByDatabaseComplete retrieves all the results into a single object
func (c SensitivityLabelsClient) ListRecommendedByDatabaseComplete(ctx context.Context, id commonids.SqlDatabaseId, options ListRecommendedByDatabaseOperationOptions) (ListRecommendedByDatabaseCompleteResult, error) {
	return c.ListRecommendedByDatabaseCompleteMatchingPredicate(ctx, id, options, SensitivityLabelOperationPredicate{})
}

// ListRecommendedByDatabaseCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c SensitivityLabelsClient) ListRecommendedByDatabaseCompleteMatchingPredicate(ctx context.Context, id commonids.SqlDatabaseId, options ListRecommendedByDatabaseOperationOptions, predicate SensitivityLabelOperationPredicate) (result ListRecommendedByDatabaseCompleteResult, err error) {
	items := make([]SensitivityLabel, 0)

	resp, err := c.ListRecommendedByDatabase(ctx, id, options)
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

	result = ListRecommendedByDatabaseCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
