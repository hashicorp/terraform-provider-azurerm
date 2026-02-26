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

type OrganizationListSchemaRegistryClustersOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SchemaRegistryClusterRecord
}

type OrganizationListSchemaRegistryClustersCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SchemaRegistryClusterRecord
}

type OrganizationListSchemaRegistryClustersOperationOptions struct {
	PageSize  *int64
	PageToken *string
}

func DefaultOrganizationListSchemaRegistryClustersOperationOptions() OrganizationListSchemaRegistryClustersOperationOptions {
	return OrganizationListSchemaRegistryClustersOperationOptions{}
}

func (o OrganizationListSchemaRegistryClustersOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o OrganizationListSchemaRegistryClustersOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o OrganizationListSchemaRegistryClustersOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.PageSize != nil {
		out.Append("pageSize", fmt.Sprintf("%v", *o.PageSize))
	}
	if o.PageToken != nil {
		out.Append("pageToken", fmt.Sprintf("%v", *o.PageToken))
	}
	return &out
}

type OrganizationListSchemaRegistryClustersCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *OrganizationListSchemaRegistryClustersCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// OrganizationListSchemaRegistryClusters ...
func (c SCEnvironmentRecordsClient) OrganizationListSchemaRegistryClusters(ctx context.Context, id EnvironmentId, options OrganizationListSchemaRegistryClustersOperationOptions) (result OrganizationListSchemaRegistryClustersOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &OrganizationListSchemaRegistryClustersCustomPager{},
		Path:          fmt.Sprintf("%s/schemaRegistryClusters", id.ID()),
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
		Values *[]SchemaRegistryClusterRecord `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// OrganizationListSchemaRegistryClustersComplete retrieves all the results into a single object
func (c SCEnvironmentRecordsClient) OrganizationListSchemaRegistryClustersComplete(ctx context.Context, id EnvironmentId, options OrganizationListSchemaRegistryClustersOperationOptions) (OrganizationListSchemaRegistryClustersCompleteResult, error) {
	return c.OrganizationListSchemaRegistryClustersCompleteMatchingPredicate(ctx, id, options, SchemaRegistryClusterRecordOperationPredicate{})
}

// OrganizationListSchemaRegistryClustersCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c SCEnvironmentRecordsClient) OrganizationListSchemaRegistryClustersCompleteMatchingPredicate(ctx context.Context, id EnvironmentId, options OrganizationListSchemaRegistryClustersOperationOptions, predicate SchemaRegistryClusterRecordOperationPredicate) (result OrganizationListSchemaRegistryClustersCompleteResult, err error) {
	items := make([]SchemaRegistryClusterRecord, 0)

	resp, err := c.OrganizationListSchemaRegistryClusters(ctx, id, options)
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

	result = OrganizationListSchemaRegistryClustersCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
