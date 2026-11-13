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

type MongoMIResourcesListMongoMIRoleAssignmentsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]MongoMIRoleAssignmentResource
}

type MongoMIResourcesListMongoMIRoleAssignmentsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []MongoMIRoleAssignmentResource
}

type MongoMIResourcesListMongoMIRoleAssignmentsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *MongoMIResourcesListMongoMIRoleAssignmentsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// MongoMIResourcesListMongoMIRoleAssignments ...
func (c OpenapisClient) MongoMIResourcesListMongoMIRoleAssignments(ctx context.Context, id DatabaseAccountId) (result MongoMIResourcesListMongoMIRoleAssignmentsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &MongoMIResourcesListMongoMIRoleAssignmentsCustomPager{},
		Path:       fmt.Sprintf("%s/mongoMIRoleAssignments", id.ID()),
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
		Values *[]MongoMIRoleAssignmentResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// MongoMIResourcesListMongoMIRoleAssignmentsComplete retrieves all the results into a single object
func (c OpenapisClient) MongoMIResourcesListMongoMIRoleAssignmentsComplete(ctx context.Context, id DatabaseAccountId) (MongoMIResourcesListMongoMIRoleAssignmentsCompleteResult, error) {
	return c.MongoMIResourcesListMongoMIRoleAssignmentsCompleteMatchingPredicate(ctx, id, MongoMIRoleAssignmentResourceOperationPredicate{})
}

// MongoMIResourcesListMongoMIRoleAssignmentsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) MongoMIResourcesListMongoMIRoleAssignmentsCompleteMatchingPredicate(ctx context.Context, id DatabaseAccountId, predicate MongoMIRoleAssignmentResourceOperationPredicate) (result MongoMIResourcesListMongoMIRoleAssignmentsCompleteResult, err error) {
	items := make([]MongoMIRoleAssignmentResource, 0)

	resp, err := c.MongoMIResourcesListMongoMIRoleAssignments(ctx, id)
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

	result = MongoMIResourcesListMongoMIRoleAssignmentsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
