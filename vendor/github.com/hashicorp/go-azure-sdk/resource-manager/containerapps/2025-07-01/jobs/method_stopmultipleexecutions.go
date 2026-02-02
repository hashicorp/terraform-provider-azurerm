package jobs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StopMultipleExecutionsOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]JobExecution
}

type StopMultipleExecutionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []JobExecution
}

type StopMultipleExecutionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *StopMultipleExecutionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// StopMultipleExecutions ...
func (c JobsClient) StopMultipleExecutions(ctx context.Context, id JobId) (result StopMultipleExecutionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &StopMultipleExecutionsCustomPager{},
		Path:       fmt.Sprintf("%s/stop", id.ID()),
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

	result.Poller, err = resourcemanager.PollerFromResponse(resp, c.Client)
	if err != nil {
		return
	}

	return
}

// StopMultipleExecutionsThenPoll performs StopMultipleExecutions then polls until it's completed
func (c JobsClient) StopMultipleExecutionsThenPoll(ctx context.Context, id JobId) error {
	result, err := c.StopMultipleExecutions(ctx, id)
	if err != nil {
		return fmt.Errorf("performing StopMultipleExecutions: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after StopMultipleExecutions: %+v", err)
	}

	return nil
}
