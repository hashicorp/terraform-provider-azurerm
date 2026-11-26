package documents

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SearchGetOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *SearchDocumentsResult
}

type SearchGetOperationOptions struct {
	Answers                       *QueryAnswerType
	Captions                      *QueryCaptionType
	Count                         *bool
	Debug                         *QueryDebugMode
	Facet                         *[]string
	Filter                        *string
	Highlight                     *[]string
	HighlightPostTag              *string
	HighlightPreTag               *string
	MinimumCoverage               *float64
	Orderby                       *[]string
	QueryType                     *QueryType
	ScoringParameter              *[]string
	ScoringProfile                *string
	ScoringStatistics             *ScoringStatistics
	Search                        *string
	SearchFields                  *[]string
	SearchMode                    *SearchMode
	Select                        *[]string
	SemanticConfiguration         *string
	SemanticErrorHandling         *SemanticErrorMode
	SemanticMaxWaitInMilliseconds *int64
	SemanticQuery                 *string
	SessionId                     *string
	Skip                          *int64
	Top                           *int64
	XMsClientRequestId            *string
}

func DefaultSearchGetOperationOptions() SearchGetOperationOptions {
	return SearchGetOperationOptions{}
}

func (o SearchGetOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.XMsClientRequestId != nil {
		out.Append("x-ms-client-request-id", fmt.Sprintf("%v", *o.XMsClientRequestId))
	}
	return &out
}

func (o SearchGetOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o SearchGetOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Answers != nil {
		out.Append("answers", fmt.Sprintf("%v", *o.Answers))
	}
	if o.Captions != nil {
		out.Append("captions", fmt.Sprintf("%v", *o.Captions))
	}
	if o.Count != nil {
		out.Append("$count", fmt.Sprintf("%v", *o.Count))
	}
	if o.Debug != nil {
		out.Append("debug", fmt.Sprintf("%v", *o.Debug))
	}
	if o.Facet != nil {
		out.Append("facet", fmt.Sprintf("%v", *o.Facet))
	}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Highlight != nil {
		out.Append("highlight", fmt.Sprintf("%v", *o.Highlight))
	}
	if o.HighlightPostTag != nil {
		out.Append("highlightPostTag", fmt.Sprintf("%v", *o.HighlightPostTag))
	}
	if o.HighlightPreTag != nil {
		out.Append("highlightPreTag", fmt.Sprintf("%v", *o.HighlightPreTag))
	}
	if o.MinimumCoverage != nil {
		out.Append("minimumCoverage", fmt.Sprintf("%v", *o.MinimumCoverage))
	}
	if o.Orderby != nil {
		out.Append("$orderby", fmt.Sprintf("%v", *o.Orderby))
	}
	if o.QueryType != nil {
		out.Append("queryType", fmt.Sprintf("%v", *o.QueryType))
	}
	if o.ScoringParameter != nil {
		out.Append("scoringParameter", fmt.Sprintf("%v", *o.ScoringParameter))
	}
	if o.ScoringProfile != nil {
		out.Append("scoringProfile", fmt.Sprintf("%v", *o.ScoringProfile))
	}
	if o.ScoringStatistics != nil {
		out.Append("scoringStatistics", fmt.Sprintf("%v", *o.ScoringStatistics))
	}
	if o.Search != nil {
		out.Append("search", fmt.Sprintf("%v", *o.Search))
	}
	if o.SearchFields != nil {
		out.Append("searchFields", fmt.Sprintf("%v", *o.SearchFields))
	}
	if o.SearchMode != nil {
		out.Append("searchMode", fmt.Sprintf("%v", *o.SearchMode))
	}
	if o.Select != nil {
		out.Append("$select", fmt.Sprintf("%v", *o.Select))
	}
	if o.SemanticConfiguration != nil {
		out.Append("semanticConfiguration", fmt.Sprintf("%v", *o.SemanticConfiguration))
	}
	if o.SemanticErrorHandling != nil {
		out.Append("semanticErrorHandling", fmt.Sprintf("%v", *o.SemanticErrorHandling))
	}
	if o.SemanticMaxWaitInMilliseconds != nil {
		out.Append("semanticMaxWaitInMilliseconds", fmt.Sprintf("%v", *o.SemanticMaxWaitInMilliseconds))
	}
	if o.SemanticQuery != nil {
		out.Append("semanticQuery", fmt.Sprintf("%v", *o.SemanticQuery))
	}
	if o.SessionId != nil {
		out.Append("sessionId", fmt.Sprintf("%v", *o.SessionId))
	}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// SearchGet ...
func (c DocumentsClient) SearchGet(ctx context.Context, options SearchGetOperationOptions) (result SearchGetOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
			http.StatusPartialContent,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Path:          "/docs",
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

	var model SearchDocumentsResult
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
