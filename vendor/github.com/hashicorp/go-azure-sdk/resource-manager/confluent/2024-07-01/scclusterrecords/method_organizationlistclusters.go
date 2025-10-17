package scclusterrecords

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OrganizationListClustersOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SCClusterRecord
}

type OrganizationListClustersCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SCClusterRecord
}

type OrganizationListClustersOperationOptions struct {
	PageSize  *int64
	PageToken *string
}

func DefaultOrganizationListClustersOperationOptions() OrganizationListClustersOperationOptions {
	return OrganizationListClustersOperationOptions{}
}

func (o OrganizationListClustersOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o OrganizationListClustersOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o OrganizationListClustersOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.PageSize != nil {
		out.Append("pageSize", fmt.Sprintf("%v", *o.PageSize))
	}
	if o.PageToken != nil {
		out.Append("pageToken", fmt.Sprintf("%v", *o.PageToken))
	}
	return &out
}

type OrganizationListClustersCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *OrganizationListClustersCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// OrganizationListClusters ...
func (c SCClusterRecordsClient) OrganizationListClusters(ctx context.Context, id EnvironmentId, options OrganizationListClustersOperationOptions) (result OrganizationListClustersOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &OrganizationListClustersCustomPager{},
		Path:          fmt.Sprintf("%s/clusters", id.ID()),
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
		Values *[]SCClusterRecord `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// OrganizationListClustersComplete retrieves all the results into a single object
func (c SCClusterRecordsClient) OrganizationListClustersComplete(ctx context.Context, id EnvironmentId, options OrganizationListClustersOperationOptions) (OrganizationListClustersCompleteResult, error) {
	return c.OrganizationListClustersCompleteMatchingPredicate(ctx, id, options, SCClusterRecordOperationPredicate{})
}

// OrganizationListClustersCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c SCClusterRecordsClient) OrganizationListClustersCompleteMatchingPredicate(ctx context.Context, id EnvironmentId, options OrganizationListClustersOperationOptions, predicate SCClusterRecordOperationPredicate) (result OrganizationListClustersCompleteResult, err error) {
	items := make([]SCClusterRecord, 0)

	resp, err := c.OrganizationListClusters(ctx, id, options)
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

	result = OrganizationListClustersCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
