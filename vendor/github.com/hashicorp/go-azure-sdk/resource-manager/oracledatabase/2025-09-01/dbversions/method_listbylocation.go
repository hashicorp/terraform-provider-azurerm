package dbversions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByLocationOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DbVersion
}

type ListByLocationCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DbVersion
}

type ListByLocationOperationOptions struct {
	DbSystemId                       *string
	DbSystemShape                    *BaseDbSystemShapes
	IsDatabaseSoftwareImageSupported *bool
	IsUpgradeSupported               *bool
	ShapeFamily                      *ShapeFamilyType
	StorageManagement                *StorageManagementType
}

func DefaultListByLocationOperationOptions() ListByLocationOperationOptions {
	return ListByLocationOperationOptions{}
}

func (o ListByLocationOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByLocationOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListByLocationOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.DbSystemId != nil {
		out.Append("dbSystemId", fmt.Sprintf("%v", *o.DbSystemId))
	}
	if o.DbSystemShape != nil {
		out.Append("dbSystemShape", fmt.Sprintf("%v", *o.DbSystemShape))
	}
	if o.IsDatabaseSoftwareImageSupported != nil {
		out.Append("isDatabaseSoftwareImageSupported", fmt.Sprintf("%v", *o.IsDatabaseSoftwareImageSupported))
	}
	if o.IsUpgradeSupported != nil {
		out.Append("isUpgradeSupported", fmt.Sprintf("%v", *o.IsUpgradeSupported))
	}
	if o.ShapeFamily != nil {
		out.Append("shapeFamily", fmt.Sprintf("%v", *o.ShapeFamily))
	}
	if o.StorageManagement != nil {
		out.Append("storageManagement", fmt.Sprintf("%v", *o.StorageManagement))
	}
	return &out
}

type ListByLocationCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByLocationCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByLocation ...
func (c DbVersionsClient) ListByLocation(ctx context.Context, id LocationId, options ListByLocationOperationOptions) (result ListByLocationOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListByLocationCustomPager{},
		Path:          fmt.Sprintf("%s/dbSystemDbVersions", id.ID()),
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
		Values *[]DbVersion `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByLocationComplete retrieves all the results into a single object
func (c DbVersionsClient) ListByLocationComplete(ctx context.Context, id LocationId, options ListByLocationOperationOptions) (ListByLocationCompleteResult, error) {
	return c.ListByLocationCompleteMatchingPredicate(ctx, id, options, DbVersionOperationPredicate{})
}

// ListByLocationCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DbVersionsClient) ListByLocationCompleteMatchingPredicate(ctx context.Context, id LocationId, options ListByLocationOperationOptions, predicate DbVersionOperationPredicate) (result ListByLocationCompleteResult, err error) {
	items := make([]DbVersion, 0)

	resp, err := c.ListByLocation(ctx, id, options)
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

	result = ListByLocationCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
