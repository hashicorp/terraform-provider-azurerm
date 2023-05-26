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

type IndicatorMetricsListOperationResponse struct {
	HttpResponse *http.Response
	Model        *ThreatIntelligenceMetricsList
}

// IndicatorMetricsList ...
func (c ThreatIntelligenceClient) IndicatorMetricsList(ctx context.Context, id WorkspaceId) (result IndicatorMetricsListOperationResponse, err error) {
	req, err := c.preparerForIndicatorMetricsList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "threatintelligence.ThreatIntelligenceClient", "IndicatorMetricsList", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "threatintelligence.ThreatIntelligenceClient", "IndicatorMetricsList", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForIndicatorMetricsList(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "threatintelligence.ThreatIntelligenceClient", "IndicatorMetricsList", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForIndicatorMetricsList prepares the IndicatorMetricsList request.
func (c ThreatIntelligenceClient) preparerForIndicatorMetricsList(ctx context.Context, id WorkspaceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.SecurityInsights/threatIntelligence/main/metrics", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForIndicatorMetricsList handles the response to the IndicatorMetricsList request. The method always
// closes the http.Response Body.
func (c ThreatIntelligenceClient) responderForIndicatorMetricsList(resp *http.Response) (result IndicatorMetricsListOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
