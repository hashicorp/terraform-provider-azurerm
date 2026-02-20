package guestconfigurationhcrpassignments

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GuestConfigurationHCRPAssignmentReportsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]GuestConfigurationAssignmentReport
}

type GuestConfigurationHCRPAssignmentReportsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []GuestConfigurationAssignmentReport
}

type GuestConfigurationHCRPAssignmentReportsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GuestConfigurationHCRPAssignmentReportsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GuestConfigurationHCRPAssignmentReportsList ...
func (c GuestConfigurationHCRPAssignmentsClient) GuestConfigurationHCRPAssignmentReportsList(ctx context.Context, id Providers2GuestConfigurationAssignmentId) (result GuestConfigurationHCRPAssignmentReportsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &GuestConfigurationHCRPAssignmentReportsListCustomPager{},
		Path:       fmt.Sprintf("%s/reports", id.ID()),
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
		Values *[]GuestConfigurationAssignmentReport `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GuestConfigurationHCRPAssignmentReportsListComplete retrieves all the results into a single object
func (c GuestConfigurationHCRPAssignmentsClient) GuestConfigurationHCRPAssignmentReportsListComplete(ctx context.Context, id Providers2GuestConfigurationAssignmentId) (GuestConfigurationHCRPAssignmentReportsListCompleteResult, error) {
	return c.GuestConfigurationHCRPAssignmentReportsListCompleteMatchingPredicate(ctx, id, GuestConfigurationAssignmentReportOperationPredicate{})
}

// GuestConfigurationHCRPAssignmentReportsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c GuestConfigurationHCRPAssignmentsClient) GuestConfigurationHCRPAssignmentReportsListCompleteMatchingPredicate(ctx context.Context, id Providers2GuestConfigurationAssignmentId, predicate GuestConfigurationAssignmentReportOperationPredicate) (result GuestConfigurationHCRPAssignmentReportsListCompleteResult, err error) {
	items := make([]GuestConfigurationAssignmentReport, 0)

	resp, err := c.GuestConfigurationHCRPAssignmentReportsList(ctx, id)
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

	result = GuestConfigurationHCRPAssignmentReportsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
