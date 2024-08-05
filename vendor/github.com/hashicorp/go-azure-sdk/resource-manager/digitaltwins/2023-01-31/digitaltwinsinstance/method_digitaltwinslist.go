package digitaltwinsinstance

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DigitalTwinsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DigitalTwinsDescription
}

type DigitalTwinsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DigitalTwinsDescription
}

type DigitalTwinsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *DigitalTwinsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// DigitalTwinsList ...
func (c DigitalTwinsInstanceClient) DigitalTwinsList(ctx context.Context, id commonids.SubscriptionId) (result DigitalTwinsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &DigitalTwinsListCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.DigitalTwins/digitalTwinsInstances", id.ID()),
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
		Values *[]DigitalTwinsDescription `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// DigitalTwinsListComplete retrieves all the results into a single object
func (c DigitalTwinsInstanceClient) DigitalTwinsListComplete(ctx context.Context, id commonids.SubscriptionId) (DigitalTwinsListCompleteResult, error) {
	return c.DigitalTwinsListCompleteMatchingPredicate(ctx, id, DigitalTwinsDescriptionOperationPredicate{})
}

// DigitalTwinsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DigitalTwinsInstanceClient) DigitalTwinsListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate DigitalTwinsDescriptionOperationPredicate) (result DigitalTwinsListCompleteResult, err error) {
	items := make([]DigitalTwinsDescription, 0)

	resp, err := c.DigitalTwinsList(ctx, id)
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

	result = DigitalTwinsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
