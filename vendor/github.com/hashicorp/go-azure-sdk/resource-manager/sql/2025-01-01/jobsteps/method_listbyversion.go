package jobsteps

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByVersionOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]JobStep
}

type ListByVersionCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []JobStep
}

type ListByVersionCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByVersionCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByVersion ...
func (c JobStepsClient) ListByVersion(ctx context.Context, id VersionId) (result ListByVersionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByVersionCustomPager{},
		Path:       fmt.Sprintf("%s/steps", id.ID()),
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
		Values *[]JobStep `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByVersionComplete retrieves all the results into a single object
func (c JobStepsClient) ListByVersionComplete(ctx context.Context, id VersionId) (ListByVersionCompleteResult, error) {
	return c.ListByVersionCompleteMatchingPredicate(ctx, id, JobStepOperationPredicate{})
}

// ListByVersionCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c JobStepsClient) ListByVersionCompleteMatchingPredicate(ctx context.Context, id VersionId, predicate JobStepOperationPredicate) (result ListByVersionCompleteResult, err error) {
	items := make([]JobStep, 0)

	resp, err := c.ListByVersion(ctx, id)
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

	result = ListByVersionCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
