package gatewaycertificateauthority

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByServiceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]GatewayCertificateAuthorityContract
}

type ListByServiceCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []GatewayCertificateAuthorityContract
}

type ListByServiceOperationOptions struct {
	Filter *string
	Skip   *int64
	Top    *int64
}

func DefaultListByServiceOperationOptions() ListByServiceOperationOptions {
	return ListByServiceOperationOptions{}
}

func (o ListByServiceOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByServiceOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListByServiceOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// ListByService ...
func (c GatewayCertificateAuthorityClient) ListByService(ctx context.Context, id GatewayId, options ListByServiceOperationOptions) (result ListByServiceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/certificateAuthorities", id.ID()),
		OptionsObject: options,
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
		Values *[]GatewayCertificateAuthorityContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByServiceComplete retrieves all the results into a single object
func (c GatewayCertificateAuthorityClient) ListByServiceComplete(ctx context.Context, id GatewayId, options ListByServiceOperationOptions) (ListByServiceCompleteResult, error) {
	return c.ListByServiceCompleteMatchingPredicate(ctx, id, options, GatewayCertificateAuthorityContractOperationPredicate{})
}

// ListByServiceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c GatewayCertificateAuthorityClient) ListByServiceCompleteMatchingPredicate(ctx context.Context, id GatewayId, options ListByServiceOperationOptions, predicate GatewayCertificateAuthorityContractOperationPredicate) (result ListByServiceCompleteResult, err error) {
	items := make([]GatewayCertificateAuthorityContract, 0)

	resp, err := c.ListByService(ctx, id, options)
	if err != nil {
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

	result = ListByServiceCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
