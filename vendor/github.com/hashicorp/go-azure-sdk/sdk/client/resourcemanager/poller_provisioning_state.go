// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resourcemanager

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

const DefaultPollingInterval = 10 * time.Second

var _ pollers.PollerType = &provisioningStatePoller{}

type provisioningStatePoller struct {
	apiVersion           string
	client               *Client
	initialRetryDuration time.Duration
	originalUri          string
	resourcePath         string
}

func provisioningStatePollerFromResponse(response *client.Response, lroIsSelfReference bool, client *Client, pollingInterval time.Duration) (*provisioningStatePoller, error) {
	// if we've gotten to this point then we're polling against a Resource Manager resource/operation of some kind
	// we next need to determine if the current URI is a Resource Manager resource, or if we should be polling on the
	// resource (e.g. `/my/resource`) rather than an operation on the resource (e.g. `/my/resource/start`)
	if response.Request == nil {
		return nil, fmt.Errorf("request was nil")
	}
	if response.Request.URL == nil {
		return nil, fmt.Errorf("request url was nil")
	}
	originalUri := response.Request.URL.RequestURI()

	// all Resource Manager operations require the `api-version` querystring
	apiVersion := response.Request.URL.Query().Get("api-version")
	if apiVersion == "" {
		return nil, fmt.Errorf("unable to determine `api-version` from %q", originalUri)
	}

	resourcePath := originalUri
	if !lroIsSelfReference {
		// if it's a self-reference (such as API Management's API/API Schema)
		path, err := resourceManagerResourcePathFromUri(originalUri)
		if err != nil {
			return nil, fmt.Errorf("determining Resource Manager Resource Path from %q: %+v", originalUri, err)
		}
		resourcePath = *path
	}

	return &provisioningStatePoller{
		apiVersion:           apiVersion,
		client:               client,
		initialRetryDuration: pollingInterval,
		originalUri:          originalUri,
		resourcePath:         resourcePath,
	}, nil
}

func (p *provisioningStatePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		OptionsObject: provisioningStateOptions{
			apiVersion: p.apiVersion,
		},
		Path: p.resourcePath,
	}
	req, err := p.client.NewRequest(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("building request: %+v", err)
	}
	resp, err := p.client.Execute(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %+v", err)
	}
	if resp == nil {
		return nil, pollers.PollingDroppedConnectionError{}
	}

	var result provisioningStateResult
	if err := resp.Unmarshal(&result); err != nil {
		return nil, fmt.Errorf("unmarshaling result: %+v", err)
	}

	status := ""
	if string(result.Status) != "" {
		status = string(result.Status)
	}
	if string(result.Properties.ProvisioningState) != "" {
		status = string(result.Properties.ProvisioningState)
	}
	if status == "" {
		// Some Operations support both an LRO and immediate completion, but _don't_ return a provisioningState field
		// since we're checking for a 200 OK, if we didn't get a provisioningState field, for the moment we have to
		// assume that we're done.
		// Examples: `APIManagement` API Versions `2021-08-01` and `2022-08-01` - `Services.GlobalSchemaCreateOrUpdate`.
		// Examples: `Automation` API Versions `2020-01-13-preview` - `DscNodeConfiguration.CreateOrUpdate`.
		// https://github.com/hashicorp/go-azure-sdk/issues/542
		return &pollers.PollResult{
			PollInterval: p.initialRetryDuration,
			Status:       pollers.PollingStatusSucceeded,
		}, nil
	}

	if strings.EqualFold(status, string(statusCanceled)) || strings.EqualFold(status, string(statusCancelled)) {
		return nil, pollers.PollingCancelledError{
			HttpResponse: resp,
		}
	}

	if strings.EqualFold(status, string(statusFailed)) {
		return nil, pollers.PollingFailedError{
			HttpResponse: resp,
		}
	}

	if strings.EqualFold(status, string(statusSucceeded)) {
		return &pollers.PollResult{
			PollInterval: p.initialRetryDuration,
			Status:       pollers.PollingStatusSucceeded,
		}, nil
	}

	// some API's have unique provisioningStates (e.g. Storage Accounts has `ResolvingDns`)
	// if we don't recognise it, treat it as a polling status
	return &pollers.PollResult{
		PollInterval: p.initialRetryDuration,
		Status:       pollers.PollingStatusInProgress,
	}, nil
}

type provisioningStateResult struct {
	Properties provisioningStateResultProperties `json:"properties"`

	// others return Status, so we check that too
	Status status `json:"status"`
}

type provisioningStateResultProperties struct {
	// Some API's (such as Storage) return the Resource Representation from the LRO API, as such we need to check provisioningState
	ProvisioningState status `json:"provisioningState"`
}

func resourceManagerResourcePathFromUri(input string) (*string, error) {
	parsed, err := url.ParseRequestURI(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	segments := strings.Split(strings.TrimPrefix(parsed.Path, "/"), "/")
	if len(segments) == 0 {
		return nil, fmt.Errorf("polling uri was empty")
	}

	// Resources within Resource Manager are always a matching list of key-value pairs
	// however Operations against a given resource (e.g. `/start`) will be an uneven pair
	// as such, if we have matching pairs, this a Resource
	if len(segments)%2 != 0 {
		segments = segments[0 : len(segments)-1]
		parsed.Path = fmt.Sprintf("/%s", strings.Join(segments, "/"))
	}

	if parsed.Path == "/" {
		return nil, fmt.Errorf("expected a Resource Manager URI but got %q", parsed.Path)
	}

	return pointer.To(parsed.Path), nil
}

var _ client.Options = provisioningStateOptions{}

type provisioningStateOptions struct {
	apiVersion string
}

func (p provisioningStateOptions) ToHeaders() *client.Headers {
	return &client.Headers{}
}

func (p provisioningStateOptions) ToOData() *odata.Query {
	return &odata.Query{}
}

func (p provisioningStateOptions) ToQuery() *client.QueryParams {
	q := client.QueryParams{}
	q.Append("api-version", p.apiVersion)
	return &q
}
