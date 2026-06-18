package integrationruntime

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

type ObjectMetadataListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SsisObjectMetadata
}

type ObjectMetadataListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SsisObjectMetadata
}

type ObjectMetadataListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ObjectMetadataListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ObjectMetadataList ...
func (c IntegrationRuntimeClient) ObjectMetadataList(ctx context.Context, id IntegrationRuntimeId, input GetSsisObjectMetadataRequest) (result ObjectMetadataListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &ObjectMetadataListCustomPager{},
		Path:       fmt.Sprintf("%s/getObjectMetadata", id.ID()),
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

	temp := make([]SsisObjectMetadata, 0)
	if values.Values != nil {
		for i, v := range *values.Values {
			val, err := UnmarshalSsisObjectMetadataImplementation(v)
			if err != nil {
				err = fmt.Errorf("unmarshalling item %d for SsisObjectMetadata (%q): %+v", i, v, err)
				return result, err
			}
			temp = append(temp, val)
		}
	}
	result.Model = &temp

	return
}

// ObjectMetadataListComplete retrieves all the results into a single object
func (c IntegrationRuntimeClient) ObjectMetadataListComplete(ctx context.Context, id IntegrationRuntimeId, input GetSsisObjectMetadataRequest) (ObjectMetadataListCompleteResult, error) {
	return c.ObjectMetadataListCompleteMatchingPredicate(ctx, id, input, SsisObjectMetadataOperationPredicate{})
}

// ObjectMetadataListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c IntegrationRuntimeClient) ObjectMetadataListCompleteMatchingPredicate(ctx context.Context, id IntegrationRuntimeId, input GetSsisObjectMetadataRequest, predicate SsisObjectMetadataOperationPredicate) (result ObjectMetadataListCompleteResult, err error) {
	items := make([]SsisObjectMetadata, 0)

	resp, err := c.ObjectMetadataList(ctx, id, input)
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

	result = ObjectMetadataListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
