package threatintelligence

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IndicatorAppendTagsOperationResponse struct {
	HttpResponse *http.Response
}

// IndicatorAppendTags ...
func (c ThreatIntelligenceClient) IndicatorAppendTags(ctx context.Context, id IndicatorId, input ThreatIntelligenceAppendTags) (result IndicatorAppendTagsOperationResponse, err error) {
	req, err := c.preparerForIndicatorAppendTags(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "threatintelligence.ThreatIntelligenceClient", "IndicatorAppendTags", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "threatintelligence.ThreatIntelligenceClient", "IndicatorAppendTags", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForIndicatorAppendTags(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "threatintelligence.ThreatIntelligenceClient", "IndicatorAppendTags", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForIndicatorAppendTags prepares the IndicatorAppendTags request.
func (c ThreatIntelligenceClient) preparerForIndicatorAppendTags(ctx context.Context, id IndicatorId, input ThreatIntelligenceAppendTags) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/appendTags", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForIndicatorAppendTags handles the response to the IndicatorAppendTags request. The method always
// closes the http.Response Body.
func (c ThreatIntelligenceClient) responderForIndicatorAppendTags(resp *http.Response) (result IndicatorAppendTagsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
