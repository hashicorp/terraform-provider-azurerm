// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package armcore

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	sdkruntime "github.com/Azure/azure-sdk-for-go/sdk/internal/runtime"
)

const (
	// LogRPRegistration entries contain information specific to the automatic registration of an RP.
	// Entries of this classification are written IFF the policy needs to take any action.
	LogRPRegistration azcore.LogClassification = "RPRegistration"
)

// RegistrationOptions configures the registration policy's behavior.
// All zero-value fields will be initialized with their default values.
type RegistrationOptions struct {
	// MaxAttempts is the total number of times to attempt automatic registration
	// in the event that an attempt fails.
	// The default value is 3.
	// Set to a value less than zero to disable the policy.
	MaxAttempts int

	// PollingDelay is the amount of time to sleep between polling intervals.
	// The default value is 15 seconds.
	// A value less than zero means no delay between polling intervals (not recommended).
	PollingDelay time.Duration

	// PollingDuration is the amount of time to wait before abandoning polling.
	// The default valule is 5 minutes.
	// NOTE: Setting this to a small value might cause the policy to prematurely fail.
	PollingDuration time.Duration

	// HTTPClient sets the transport for making HTTP requests.
	HTTPClient azcore.Transport

	// Retry configures the built-in retry policy behavior.
	Retry azcore.RetryOptions

	// Telemetry configures the built-in telemetry policy behavior.
	Telemetry azcore.TelemetryOptions

	// Logging configures the built-in logging policy behavior.
	Logging azcore.LogOptions
}

// init sets any default values
func (r *RegistrationOptions) init() {
	if r.MaxAttempts == 0 {
		r.MaxAttempts = 3
	} else if r.MaxAttempts < 0 {
		r.MaxAttempts = 0
	}
	if r.PollingDelay == 0 {
		r.PollingDelay = 15 * time.Second
	} else if r.PollingDelay < 0 {
		r.PollingDelay = 0
	}
	if r.PollingDuration == 0 {
		r.PollingDuration = 5 * time.Minute
	}
}

// NewRPRegistrationPolicy creates a policy object configured using the specified endpoint,
// credentials and options.  The policy controls if an unregistered resource provider should
// automatically be registered. See https://aka.ms/rps-not-found for more information.
// Pass nil to accept the default options; this is the same as passing a zero-value options.
func NewRPRegistrationPolicy(endpoint string, cred azcore.Credential, o *RegistrationOptions) azcore.Policy {
	if o == nil {
		o = &RegistrationOptions{}
	}
	p := &rpRegistrationPolicy{
		endpoint: endpoint,
		pipeline: azcore.NewPipeline(o.HTTPClient,
			azcore.NewTelemetryPolicy(&o.Telemetry),
			azcore.NewRetryPolicy(&o.Retry),
			cred.AuthenticationPolicy(azcore.AuthenticationPolicyOptions{Options: azcore.TokenRequestOptions{Scopes: []string{endpointToScope(endpoint)}}}),
			azcore.NewLogPolicy(&o.Logging)),
		options: *o,
	}
	// init the copy
	p.options.init()
	return p
}

type rpRegistrationPolicy struct {
	endpoint string
	pipeline azcore.Pipeline
	options  RegistrationOptions
}

func (r *rpRegistrationPolicy) Do(req *azcore.Request) (*azcore.Response, error) {
	if r.options.MaxAttempts == 0 {
		// policy is disabled
		return req.Next()
	}
	const unregisteredRPCode = "MissingSubscriptionRegistration"
	const registeredState = "Registered"
	var rp string
	var resp *azcore.Response
	for attempts := 0; attempts < r.options.MaxAttempts; attempts++ {
		var err error
		// make the original request
		resp, err = req.Next()
		// getting a 409 is the first indication that the RP might need to be registered, check error response
		if err != nil || resp.StatusCode != http.StatusConflict {
			return resp, err
		}
		var reqErr requestError
		if err = resp.UnmarshalAsJSON(&reqErr); err != nil {
			return resp, newFrameError(err)
		}
		if reqErr.ServiceError == nil {
			return resp, newFrameError(errors.New("missing error information"))
		}
		if !strings.EqualFold(reqErr.ServiceError.Code, unregisteredRPCode) {
			// not a 409 due to unregistered RP
			return resp, err
		}
		// RP needs to be registered.  start by getting the subscription ID from the original request
		subID, err := getSubscription(req.URL.Path)
		if err != nil {
			return resp, newFrameError(err)
		}
		// now get the RP from the error
		rp, err = getProvider(reqErr)
		if err != nil {
			return resp, newFrameError(err)
		}
		logRegistrationExit := func(v interface{}) {
			azcore.Log().Writef(LogRPRegistration, "END registration for %s: %v", rp, v)
		}
		azcore.Log().Writef(LogRPRegistration, "BEGIN registration for %s", rp)
		// create client and make the registration request
		// we use the scheme and host from the original request
		rpOps := &providersOperations{
			p:     r.pipeline,
			u:     r.endpoint,
			subID: subID,
		}
		if _, err = rpOps.Register(req.Context(), rp); err != nil {
			logRegistrationExit(err)
			return resp, err
		}
		// RP was registered, however we need to wait for the registration to complete
		pollCtx, pollCancel := context.WithTimeout(req.Context(), r.options.PollingDuration)
		var lastRegState string
		for {
			// get the current registration state
			getResp, err := rpOps.Get(pollCtx, rp)
			if err != nil {
				pollCancel()
				logRegistrationExit(err)
				return resp, err
			}
			if getResp.Provider.RegistrationState != nil && !strings.EqualFold(*getResp.Provider.RegistrationState, lastRegState) {
				// registration state has changed, or was updated for the first time
				lastRegState = *getResp.Provider.RegistrationState
				azcore.Log().Writef(LogRPRegistration, "registration state is %s", lastRegState)
			}
			if strings.EqualFold(lastRegState, registeredState) {
				// registration complete
				pollCancel()
				logRegistrationExit(lastRegState)
				break
			}
			// wait before trying again
			select {
			case <-time.After(r.options.PollingDelay):
				// continue polling
			case <-pollCtx.Done():
				pollCancel()
				logRegistrationExit(pollCtx.Err())
				return resp, pollCtx.Err()
			}
		}
		// RP was successfully registered, retry the original request
		err = req.RewindBody()
		if err != nil {
			return resp, newFrameError(err)
		}
	}
	// if we get here it means we exceeded the number of attempts
	return resp, fmt.Errorf("exceeded attempts to register %s", rp)
}

func newFrameError(inner error) error {
	// skip ourselves
	return sdkruntime.NewFrameError(inner, false, 1, azcore.StackFrameCount)
}

func getSubscription(path string) (string, error) {
	parts := strings.Split(path, "/")
	for i, v := range parts {
		if v == "subscriptions" && (i+1) < len(parts) {
			return parts[i+1], nil
		}
	}
	return "", fmt.Errorf("failed to obtain subscription ID from %s", path)
}

func getProvider(re requestError) (string, error) {
	if len(re.ServiceError.Details) > 0 {
		return re.ServiceError.Details[0].Target, nil
	}
	return "", errors.New("unexpected empty Details")
}

// minimal error definitions to simplify detection
type requestError struct {
	ServiceError *serviceError `json:"error"`
}

type serviceError struct {
	Code    string                `json:"code"`
	Details []serviceErrorDetails `json:"details"`
}

type serviceErrorDetails struct {
	Code   string `json:"code"`
	Target string `json:"target"`
}

///////////////////////////////////////////////////////////////////////////////////////////////
// the following code was copied from module armresources, providers.go and models.go
// only the minimum amount of code was copied to get this working and some edits were made.
///////////////////////////////////////////////////////////////////////////////////////////////

type providersOperations struct {
	p     azcore.Pipeline
	u     string
	subID string
}

// Get - Gets the specified resource provider.
func (client *providersOperations) Get(ctx context.Context, resourceProviderNamespace string) (*ProviderResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceProviderNamespace)
	if err != nil {
		return nil, err
	}
	resp, err := client.p.Do(req)
	if err != nil {
		return nil, err
	}
	result, err := client.getHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// getCreateRequest creates the Get request.
func (client *providersOperations) getCreateRequest(ctx context.Context, resourceProviderNamespace string) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/{resourceProviderNamespace}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceProviderNamespace}", url.PathEscape(resourceProviderNamespace))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subID))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.u, urlPath))
	if err != nil {
		return nil, newFrameError(err)
	}
	query := req.URL.Query()
	query.Set("api-version", "2019-05-01")
	req.URL.RawQuery = query.Encode()
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *providersOperations) getHandleResponse(resp *azcore.Response) (*ProviderResponse, error) {
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.getHandleError(resp)
	}
	result := ProviderResponse{RawResponse: resp.Response}
	err := resp.UnmarshalAsJSON(&result.Provider)
	if err != nil {
		err = newFrameError(err)
	}
	return &result, err
}

// getHandleError handles the Get error response.
func (client *providersOperations) getHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return sdkruntime.NewResponseError(newFrameError(err), resp.Response)
	}
	if len(body) == 0 {
		return sdkruntime.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return sdkruntime.NewResponseError(errors.New(string(body)), resp.Response)
}

// Register - Registers a subscription with a resource provider.
func (client *providersOperations) Register(ctx context.Context, resourceProviderNamespace string) (*ProviderResponse, error) {
	req, err := client.registerCreateRequest(ctx, resourceProviderNamespace)
	if err != nil {
		return nil, err
	}
	resp, err := client.p.Do(req)
	if err != nil {
		return nil, err
	}
	result, err := client.registerHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// registerCreateRequest creates the Register request.
func (client *providersOperations) registerCreateRequest(ctx context.Context, resourceProviderNamespace string) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/{resourceProviderNamespace}/register"
	urlPath = strings.ReplaceAll(urlPath, "{resourceProviderNamespace}", url.PathEscape(resourceProviderNamespace))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subID))
	req, err := azcore.NewRequest(ctx, http.MethodPost, azcore.JoinPaths(client.u, urlPath))
	if err != nil {
		return nil, newFrameError(err)
	}
	query := req.URL.Query()
	query.Set("api-version", "2019-05-01")
	req.URL.RawQuery = query.Encode()
	return req, nil
}

// registerHandleResponse handles the Register response.
func (client *providersOperations) registerHandleResponse(resp *azcore.Response) (*ProviderResponse, error) {
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.registerHandleError(resp)
	}
	result := ProviderResponse{RawResponse: resp.Response}
	err := resp.UnmarshalAsJSON(&result.Provider)
	if err != nil {
		err = newFrameError(err)
	}
	return &result, err
}

// registerHandleError handles the Register error response.
func (client *providersOperations) registerHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return sdkruntime.NewResponseError(newFrameError(err), resp.Response)
	}
	if len(body) == 0 {
		return sdkruntime.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return sdkruntime.NewResponseError(errors.New(string(body)), resp.Response)
}

// ProviderResponse is the response envelope for operations that return a Provider type.
type ProviderResponse struct {
	// Resource provider information.
	Provider *Provider

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// Provider - Resource provider information.
type Provider struct {
	// The provider ID.
	ID *string `json:"id,omitempty"`

	// The namespace of the resource provider.
	Namespace *string `json:"namespace,omitempty"`

	// The registration policy of the resource provider.
	RegistrationPolicy *string `json:"registrationPolicy,omitempty"`

	// The registration state of the resource provider.
	RegistrationState *string `json:"registrationState,omitempty"`
}
