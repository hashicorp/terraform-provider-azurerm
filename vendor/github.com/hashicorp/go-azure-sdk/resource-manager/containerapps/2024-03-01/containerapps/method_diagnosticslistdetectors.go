package containerapps

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiagnosticsListDetectorsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Diagnostics
}

type DiagnosticsListDetectorsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Diagnostics
}

type DiagnosticsListDetectorsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *DiagnosticsListDetectorsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// DiagnosticsListDetectors ...
func (c ContainerAppsClient) DiagnosticsListDetectors(ctx context.Context, id ContainerAppId) (result DiagnosticsListDetectorsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &DiagnosticsListDetectorsCustomPager{},
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

// DiagnosticsListDetectorsComplete retrieves all the results into a single object
func (c ContainerAppsClient) DiagnosticsListDetectorsComplete(ctx context.Context, id ContainerAppId) (DiagnosticsListDetectorsCompleteResult, error) {
	return c.DiagnosticsListDetectorsCompleteMatchingPredicate(ctx, id, DiagnosticsOperationPredicate{})
}

// DiagnosticsListDetectorsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ContainerAppsClient) DiagnosticsListDetectorsCompleteMatchingPredicate(ctx context.Context, id ContainerAppId, predicate DiagnosticsOperationPredicate) (result DiagnosticsListDetectorsCompleteResult, err error) {
	items := make([]Diagnostics, 0)

	resp, err := c.DiagnosticsListDetectors(ctx, id)
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

	result = DiagnosticsListDetectorsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
