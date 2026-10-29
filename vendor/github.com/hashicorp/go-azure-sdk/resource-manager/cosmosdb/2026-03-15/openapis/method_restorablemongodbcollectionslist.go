package openapis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RestorableMongodbCollectionsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RestorableMongodbCollectionGetResult
}

type RestorableMongodbCollectionsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RestorableMongodbCollectionGetResult
}

type RestorableMongodbCollectionsListOperationOptions struct {
	EndTime                      *string
	RestorableMongodbDatabaseRid *string
	StartTime                    *string
}

func DefaultRestorableMongodbCollectionsListOperationOptions() RestorableMongodbCollectionsListOperationOptions {
	return RestorableMongodbCollectionsListOperationOptions{}
}

func (o RestorableMongodbCollectionsListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RestorableMongodbCollectionsListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o RestorableMongodbCollectionsListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.EndTime != nil {
		out.Append("endTime", fmt.Sprintf("%v", *o.EndTime))
	}
	if o.RestorableMongodbDatabaseRid != nil {
		out.Append("restorableMongodbDatabaseRid", fmt.Sprintf("%v", *o.RestorableMongodbDatabaseRid))
	}
	if o.StartTime != nil {
		out.Append("startTime", fmt.Sprintf("%v", *o.StartTime))
	}
	return &out
}

type RestorableMongodbCollectionsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *RestorableMongodbCollectionsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// RestorableMongodbCollectionsList ...
func (c OpenapisClient) RestorableMongodbCollectionsList(ctx context.Context, id RestorableDatabaseAccountId, options RestorableMongodbCollectionsListOperationOptions) (result RestorableMongodbCollectionsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &RestorableMongodbCollectionsListCustomPager{},
		Path:          fmt.Sprintf("%s/restorableMongodbCollections", id.ID()),
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
		Values *[]RestorableMongodbCollectionGetResult `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// RestorableMongodbCollectionsListComplete retrieves all the results into a single object
func (c OpenapisClient) RestorableMongodbCollectionsListComplete(ctx context.Context, id RestorableDatabaseAccountId, options RestorableMongodbCollectionsListOperationOptions) (RestorableMongodbCollectionsListCompleteResult, error) {
	return c.RestorableMongodbCollectionsListCompleteMatchingPredicate(ctx, id, options, RestorableMongodbCollectionGetResultOperationPredicate{})
}

// RestorableMongodbCollectionsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) RestorableMongodbCollectionsListCompleteMatchingPredicate(ctx context.Context, id RestorableDatabaseAccountId, options RestorableMongodbCollectionsListOperationOptions, predicate RestorableMongodbCollectionGetResultOperationPredicate) (result RestorableMongodbCollectionsListCompleteResult, err error) {
	items := make([]RestorableMongodbCollectionGetResult, 0)

	resp, err := c.RestorableMongodbCollectionsList(ctx, id, options)
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

	result = RestorableMongodbCollectionsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
