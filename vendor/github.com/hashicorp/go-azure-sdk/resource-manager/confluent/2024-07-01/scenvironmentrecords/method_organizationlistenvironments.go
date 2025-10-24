package scenvironmentrecords

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OrganizationListEnvironmentsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SCEnvironmentRecord
}

type OrganizationListEnvironmentsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SCEnvironmentRecord
}

type OrganizationListEnvironmentsOperationOptions struct {
	PageSize  *int64
	PageToken *string
}

func DefaultOrganizationListEnvironmentsOperationOptions() OrganizationListEnvironmentsOperationOptions {
	return OrganizationListEnvironmentsOperationOptions{}
}

func (o OrganizationListEnvironmentsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o OrganizationListEnvironmentsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o OrganizationListEnvironmentsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.PageSize != nil {
		out.Append("pageSize", fmt.Sprintf("%v", *o.PageSize))
	}
	if o.PageToken != nil {
		out.Append("pageToken", fmt.Sprintf("%v", *o.PageToken))
	}
	return &out
}

type OrganizationListEnvironmentsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *OrganizationListEnvironmentsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// OrganizationListEnvironments ...
func (c SCEnvironmentRecordsClient) OrganizationListEnvironments(ctx context.Context, id OrganizationId, options OrganizationListEnvironmentsOperationOptions) (result OrganizationListEnvironmentsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &OrganizationListEnvironmentsCustomPager{},
		Path:          fmt.Sprintf("%s/environments", id.ID()),
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
		Values *[]SCEnvironmentRecord `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// OrganizationListEnvironmentsComplete retrieves all the results into a single object
func (c SCEnvironmentRecordsClient) OrganizationListEnvironmentsComplete(ctx context.Context, id OrganizationId, options OrganizationListEnvironmentsOperationOptions) (OrganizationListEnvironmentsCompleteResult, error) {
	return c.OrganizationListEnvironmentsCompleteMatchingPredicate(ctx, id, options, SCEnvironmentRecordOperationPredicate{})
}

// OrganizationListEnvironmentsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c SCEnvironmentRecordsClient) OrganizationListEnvironmentsCompleteMatchingPredicate(ctx context.Context, id OrganizationId, options OrganizationListEnvironmentsOperationOptions, predicate SCEnvironmentRecordOperationPredicate) (result OrganizationListEnvironmentsCompleteResult, err error) {
	items := make([]SCEnvironmentRecord, 0)

	resp, err := c.OrganizationListEnvironments(ctx, id, options)
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

	result = OrganizationListEnvironmentsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
