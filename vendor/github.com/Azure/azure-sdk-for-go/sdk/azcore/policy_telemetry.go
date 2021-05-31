// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
)

// TelemetryOptions configures the telemetry policy's behavior.
type TelemetryOptions struct {
	// Value is a string prepended to each request's User-Agent and sent to the service.
	// The service records the user-agent in logs for diagnostics and tracking of client requests.
	Value string

	// ApplicationID is an application-specific identification string used in telemetry.
	// It has a maximum length of 24 characters and must not contain any spaces.
	ApplicationID string

	// Disabled will prevent the addition of any telemetry data to the User-Agent.
	Disabled bool
}

type telemetryPolicy struct {
	telemetryValue string
}

// NewTelemetryPolicy creates a telemetry policy object that adds telemetry information to outgoing HTTP requests.
// The format is [<application_id> ]azsdk-<sdk_language>-<package_name>/<package_version> <platform_info> [<custom>].
// Pass nil to accept the default values; this is the same as passing a zero-value options.
func NewTelemetryPolicy(o *TelemetryOptions) Policy {
	if o == nil {
		o = &TelemetryOptions{}
	}
	tp := telemetryPolicy{}
	if o.Disabled {
		return &tp
	}
	b := &bytes.Buffer{}
	// normalize ApplicationID
	if o.ApplicationID != "" {
		o.ApplicationID = strings.ReplaceAll(o.ApplicationID, " ", "/")
		if len(o.ApplicationID) > 24 {
			o.ApplicationID = o.ApplicationID[:24]
		}
		b.WriteString(o.ApplicationID)
		b.WriteRune(' ')
	}
	// write out telemetry string
	if o.Value != "" {
		b.WriteString(o.Value)
		b.WriteRune(' ')
	}
	b.WriteString(UserAgent)
	b.WriteRune(' ')
	b.WriteString(platformInfo)
	tp.telemetryValue = b.String()
	return &tp
}

func (p telemetryPolicy) Do(req *Request) (*Response, error) {
	if p.telemetryValue == "" {
		return req.Next()
	}
	// preserve the existing User-Agent string
	if ua := req.Request.Header.Get(HeaderUserAgent); ua != "" {
		p.telemetryValue = fmt.Sprintf("%s %s", p.telemetryValue, ua)
	}
	var rt requestTelemetry
	if req.OperationValue(&rt) {
		p.telemetryValue = fmt.Sprintf("%s %s", string(rt), p.telemetryValue)
	}
	req.Request.Header.Set(HeaderUserAgent, p.telemetryValue)
	return req.Next()
}

// NOTE: the ONLY function that should write to this variable is this func
var platformInfo = func() string {
	operatingSystem := runtime.GOOS // Default OS string
	switch operatingSystem {
	case "windows":
		operatingSystem = os.Getenv("OS") // Get more specific OS information
	case "linux": // accept default OS info
	case "freebsd": //  accept default OS info
	}
	return fmt.Sprintf("(%s; %s)", runtime.Version(), operatingSystem)
}()

// used for adding per-request telemetry
type requestTelemetry string
