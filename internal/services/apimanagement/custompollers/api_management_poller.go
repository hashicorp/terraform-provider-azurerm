// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &apiManagementPoller{}

type apiManagementPoller struct {
	client     *apimanagementservice.ApiManagementServiceClient
	pollingUrl *url.URL
}

type operationResult struct {
	Status status `json:"status"`
}

type status string

var (
	pollingDeleteSuccess = pollers.PollResult{
		Status:       pollers.PollingStatusSucceeded,
		PollInterval: 20 * time.Second,
	}
	pollingDeleteInProgress = pollers.PollResult{
		Status:       pollers.PollingStatusInProgress,
		PollInterval: 20 * time.Second,
	}
)

// NewAPIManagementPoller - creates a new poller for API Management long-running delete operation to handle the case where the delete operation is asynchronous.
func NewAPIManagementPoller(cli *apimanagementservice.ApiManagementServiceClient, response *http.Response) (*apiManagementPoller, error) {
	pollingUrl := response.Header.Get("Azure-AsyncOperation")
	if pollingUrl == "" {
		pollingUrl = response.Header.Get("Location")
	}

	if pollingUrl == "" {
		return nil, fmt.Errorf("no polling URL found in response")
	}

	url, err := url.Parse(pollingUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid polling URL %q in response: %v", pollingUrl, err)
	}
	if !url.IsAbs() {
		return nil, fmt.Errorf("invalid polling URL %q in response: URL was not absolute", pollingUrl)
	}

	return &apiManagementPoller{
		client:     cli,
		pollingUrl: url,
	}, nil
}

func (p apiManagementPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	if p.pollingUrl == nil {
		return nil, fmt.Errorf("internal error: cannot poll without a pollingUrl")
	}

	reqOpts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
			http.StatusCreated,
			http.StatusAccepted,
			http.StatusNoContent,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: nil,
		Path:          p.pollingUrl.Path,
	}

	req, err := p.client.Client.NewRequest(ctx, reqOpts)
	if err != nil {
		return nil, fmt.Errorf("building request: %+v", err)
	}

	resp, err := p.client.Client.Execute(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", p.pollingUrl.String(), err)
	}

	if resp.Response != nil {
		var respBody []byte
		respBody, err = io.ReadAll(resp.Response.Body)
		if err != nil {
			return nil, fmt.Errorf("parsing response body: %+v", err)
		}
		resp.Response.Body.Close()

		resp.Response.Body = io.NopCloser(bytes.NewReader(respBody))

		if s, ok := resp.Response.Header["Retry-After"]; ok {
			if sleep, err := strconv.ParseInt(s[0], 10, 64); err == nil {
				pollingDeleteInProgress.PollInterval = time.Second * time.Duration(sleep)
			}
		}

		// 202's don't necessarily return a body, so there's nothing to deserialize
		if resp.StatusCode == http.StatusAccepted && resp.ContentLength == 0 {
			return &pollingDeleteInProgress, nil
		}

		// returns a 200 OK with no Body
		if resp.StatusCode == http.StatusOK && resp.ContentLength == 0 {
			return &pollingDeleteSuccess, nil
		}

		if resp.Response.StatusCode == http.StatusOK {
			contentType := resp.Response.Header.Get("Content-Type")
			var op operationResult
			if strings.Contains(strings.ToLower(contentType), "application/json") {
				if err = json.Unmarshal(respBody, &op); err != nil {
					return nil, fmt.Errorf("unmarshalling response body: %+v", err)
				}
			} else {
				return nil, fmt.Errorf("internal-error: polling support for the Content-Type %q was not implemented: %+v", contentType, err)
			}

			switch string(op.Status) {
			case string(pollers.PollingStatusInProgress):
				return &pollingDeleteInProgress, nil
			case string(pollers.PollingStatusSucceeded):
				return &pollingDeleteSuccess, nil
			}
		}
	}

	return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
}
