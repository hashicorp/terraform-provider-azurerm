// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

var _ pollers.PollerType = &apiManagementAPIPoller{}

type apiManagementAPIPoller struct {
	client  *api.ApiClient
	id      api.ApiId
	asyncID string
}

var (
	pollingSuccess = pollers.PollResult{
		Status:       pollers.PollingStatusSucceeded,
		PollInterval: 10 * time.Second,
	}
	pollingInProgress = pollers.PollResult{
		Status:       pollers.PollingStatusInProgress,
		PollInterval: 10 * time.Second,
	}
)

// NewAPIManagementAPIPoller - creates a new poller for API Management API operations to handle the case there is a query string
// parameter "asyncId" in the Location header of the response. This is used to poll the status of the operation.
func NewAPIManagementAPIPoller(cli *api.ApiClient, id api.ApiId, response *http.Response) *apiManagementAPIPoller {
	urlStr := response.Header.Get("location")
	var asyncId string
	if u, err := url.Parse(urlStr); err == nil {
		asyncId = u.Query().Get("asyncId")
	}

	// sometimes the poller is not required as the API directly return 200
	if asyncId == "" {
		return nil
	}

	return &apiManagementAPIPoller{
		client:  cli,
		id:      id,
		asyncID: asyncId,
	}
}

type options struct {
	asyncId string
}

func (p options) ToHeaders() *client.Headers {
	return &client.Headers{}
}

func (p options) ToOData() *odata.Query {
	return &odata.Query{}
}

func (p options) ToQuery() *client.QueryParams {
	q := client.QueryParams{}
	q.Append("asyncId", p.asyncId)
	return &q
}

func (p apiManagementAPIPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	if p.asyncID == "" {
		return &pollingSuccess, nil
	}

	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
			http.StatusAccepted,
			http.StatusCreated,
		},
		HttpMethod: http.MethodGet,
		Path:       p.id.ID(),
		OptionsObject: options{
			asyncId: p.asyncID,
		},
	}
	req, err := p.client.Client.NewRequest(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("building request: %+v", err)
	}

	resp, err := p.client.Client.Execute(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	// the response actually doesn't include a provisioningState property, so we only chech the http status code
	switch resp.StatusCode {
	case http.StatusOK:
		return &pollingSuccess, nil
	case http.StatusAccepted, http.StatusCreated:
		return &pollingInProgress, nil
	}

	return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
}
