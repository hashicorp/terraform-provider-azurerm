package profiles

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LogAnalyticsGetLogAnalyticsRankingsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *RankingsResponse
}

type LogAnalyticsGetLogAnalyticsRankingsOperationOptions struct {
	CustomDomains *[]string
	DateTimeBegin *string
	DateTimeEnd   *string
	MaxRanking    *int64
	Metrics       *[]string
	Rankings      *[]string
}

func DefaultLogAnalyticsGetLogAnalyticsRankingsOperationOptions() LogAnalyticsGetLogAnalyticsRankingsOperationOptions {
	return LogAnalyticsGetLogAnalyticsRankingsOperationOptions{}
}

func (o LogAnalyticsGetLogAnalyticsRankingsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o LogAnalyticsGetLogAnalyticsRankingsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o LogAnalyticsGetLogAnalyticsRankingsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.CustomDomains != nil {
		out.Append("customDomains", fmt.Sprintf("%v", *o.CustomDomains))
	}
	if o.DateTimeBegin != nil {
		out.Append("dateTimeBegin", fmt.Sprintf("%v", *o.DateTimeBegin))
	}
	if o.DateTimeEnd != nil {
		out.Append("dateTimeEnd", fmt.Sprintf("%v", *o.DateTimeEnd))
	}
	if o.MaxRanking != nil {
		out.Append("maxRanking", fmt.Sprintf("%v", *o.MaxRanking))
	}
	if o.Metrics != nil {
		out.Append("metrics", fmt.Sprintf("%v", *o.Metrics))
	}
	if o.Rankings != nil {
		out.Append("rankings", fmt.Sprintf("%v", *o.Rankings))
	}
	return &out
}

// LogAnalyticsGetLogAnalyticsRankings ...
func (c ProfilesClient) LogAnalyticsGetLogAnalyticsRankings(ctx context.Context, id ProfileId, options LogAnalyticsGetLogAnalyticsRankingsOperationOptions) (result LogAnalyticsGetLogAnalyticsRankingsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/getLogAnalyticsRankings", id.ID()),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	var model RankingsResponse
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
