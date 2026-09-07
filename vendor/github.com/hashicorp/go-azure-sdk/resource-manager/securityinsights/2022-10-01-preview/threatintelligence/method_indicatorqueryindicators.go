package threatintelligence

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IndicatorQueryIndicatorsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ThreatIntelligenceInformation
}

type IndicatorQueryIndicatorsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ThreatIntelligenceInformation
}

type IndicatorQueryIndicatorsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *IndicatorQueryIndicatorsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// IndicatorQueryIndicators ...
func (c ThreatIntelligenceClient) IndicatorQueryIndicators(ctx context.Context, id WorkspaceId, input ThreatIntelligenceFilteringCriteria) (result IndicatorQueryIndicatorsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &IndicatorQueryIndicatorsCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.SecurityInsights/threatIntelligence/main/queryIndicators", id.ID()),
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
		Values *[]json.RawMessage `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	temp := make([]ThreatIntelligenceInformation, 0)
	if values.Values != nil {
		for i, v := range *values.Values {
			val, err := UnmarshalThreatIntelligenceInformationImplementation(v)
			if err != nil {
				err = fmt.Errorf("unmarshalling item %d for ThreatIntelligenceInformation (%q): %+v", i, v, err)
				return result, err
			}
			temp = append(temp, val)
		}
	}
	result.Model = &temp

	return
}

// IndicatorQueryIndicatorsComplete retrieves all the results into a single object
func (c ThreatIntelligenceClient) IndicatorQueryIndicatorsComplete(ctx context.Context, id WorkspaceId, input ThreatIntelligenceFilteringCriteria) (IndicatorQueryIndicatorsCompleteResult, error) {
	return c.IndicatorQueryIndicatorsCompleteMatchingPredicate(ctx, id, input, ThreatIntelligenceInformationOperationPredicate{})
}

// IndicatorQueryIndicatorsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ThreatIntelligenceClient) IndicatorQueryIndicatorsCompleteMatchingPredicate(ctx context.Context, id WorkspaceId, input ThreatIntelligenceFilteringCriteria, predicate ThreatIntelligenceInformationOperationPredicate) (result IndicatorQueryIndicatorsCompleteResult, err error) {
	items := make([]ThreatIntelligenceInformation, 0)

	resp, err := c.IndicatorQueryIndicators(ctx, id, input)
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

	result = IndicatorQueryIndicatorsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
