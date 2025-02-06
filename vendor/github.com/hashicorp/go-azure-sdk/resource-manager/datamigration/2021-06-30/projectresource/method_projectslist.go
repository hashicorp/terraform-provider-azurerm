package projectresource

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProjectsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Project
}

type ProjectsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Project
}

type ProjectsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ProjectsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ProjectsList ...
func (c ProjectResourceClient) ProjectsList(ctx context.Context, id ServiceId) (result ProjectsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ProjectsListCustomPager{},
		Path:       fmt.Sprintf("%s/projects", id.ID()),
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
		Values *[]Project `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ProjectsListComplete retrieves all the results into a single object
func (c ProjectResourceClient) ProjectsListComplete(ctx context.Context, id ServiceId) (ProjectsListCompleteResult, error) {
	return c.ProjectsListCompleteMatchingPredicate(ctx, id, ProjectOperationPredicate{})
}

// ProjectsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ProjectResourceClient) ProjectsListCompleteMatchingPredicate(ctx context.Context, id ServiceId, predicate ProjectOperationPredicate) (result ProjectsListCompleteResult, err error) {
	items := make([]Project, 0)

	resp, err := c.ProjectsList(ctx, id)
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

	result = ProjectsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
