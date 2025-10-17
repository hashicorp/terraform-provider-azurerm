package topicrecords

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TopicsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]TopicRecord
}

type TopicsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []TopicRecord
}

type TopicsListOperationOptions struct {
	PageSize  *int64
	PageToken *string
}

func DefaultTopicsListOperationOptions() TopicsListOperationOptions {
	return TopicsListOperationOptions{}
}

func (o TopicsListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o TopicsListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o TopicsListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.PageSize != nil {
		out.Append("pageSize", fmt.Sprintf("%v", *o.PageSize))
	}
	if o.PageToken != nil {
		out.Append("pageToken", fmt.Sprintf("%v", *o.PageToken))
	}
	return &out
}

type TopicsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *TopicsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// TopicsList ...
func (c TopicRecordsClient) TopicsList(ctx context.Context, id ClusterId, options TopicsListOperationOptions) (result TopicsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &TopicsListCustomPager{},
		Path:          fmt.Sprintf("%s/topics", id.ID()),
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
		Values *[]TopicRecord `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// TopicsListComplete retrieves all the results into a single object
func (c TopicRecordsClient) TopicsListComplete(ctx context.Context, id ClusterId, options TopicsListOperationOptions) (TopicsListCompleteResult, error) {
	return c.TopicsListCompleteMatchingPredicate(ctx, id, options, TopicRecordOperationPredicate{})
}

// TopicsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c TopicRecordsClient) TopicsListCompleteMatchingPredicate(ctx context.Context, id ClusterId, options TopicsListOperationOptions, predicate TopicRecordOperationPredicate) (result TopicsListCompleteResult, err error) {
	items := make([]TopicRecord, 0)

	resp, err := c.TopicsList(ctx, id, options)
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

	result = TopicsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
