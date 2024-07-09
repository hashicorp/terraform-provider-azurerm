package managedenvironments

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedEnvironmentDiagnosticsListDetectorsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Diagnostics
}

type ManagedEnvironmentDiagnosticsListDetectorsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Diagnostics
}

type ManagedEnvironmentDiagnosticsListDetectorsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ManagedEnvironmentDiagnosticsListDetectorsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ManagedEnvironmentDiagnosticsListDetectors ...
func (c ManagedEnvironmentsClient) ManagedEnvironmentDiagnosticsListDetectors(ctx context.Context, id ManagedEnvironmentId) (result ManagedEnvironmentDiagnosticsListDetectorsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ManagedEnvironmentDiagnosticsListDetectorsCustomPager{},
		Path:       fmt.Sprintf("%s/detectors", id.ID()),
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
		Values *[]Diagnostics `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ManagedEnvironmentDiagnosticsListDetectorsComplete retrieves all the results into a single object
func (c ManagedEnvironmentsClient) ManagedEnvironmentDiagnosticsListDetectorsComplete(ctx context.Context, id ManagedEnvironmentId) (ManagedEnvironmentDiagnosticsListDetectorsCompleteResult, error) {
	return c.ManagedEnvironmentDiagnosticsListDetectorsCompleteMatchingPredicate(ctx, id, DiagnosticsOperationPredicate{})
}

// ManagedEnvironmentDiagnosticsListDetectorsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ManagedEnvironmentsClient) ManagedEnvironmentDiagnosticsListDetectorsCompleteMatchingPredicate(ctx context.Context, id ManagedEnvironmentId, predicate DiagnosticsOperationPredicate) (result ManagedEnvironmentDiagnosticsListDetectorsCompleteResult, err error) {
	items := make([]Diagnostics, 0)

	resp, err := c.ManagedEnvironmentDiagnosticsListDetectors(ctx, id)
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

	result = ManagedEnvironmentDiagnosticsListDetectorsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
