package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LocalDNSOverride struct {
	CacheDurationInSeconds      *int64                      `json:"cacheDurationInSeconds,omitempty"`
	ForwardDestination          *LocalDNSForwardDestination `json:"forwardDestination,omitempty"`
	ForwardPolicy               *LocalDNSForwardPolicy      `json:"forwardPolicy,omitempty"`
	MaxConcurrent               *int64                      `json:"maxConcurrent,omitempty"`
	Protocol                    *LocalDNSProtocol           `json:"protocol,omitempty"`
	QueryLogging                *LocalDNSQueryLogging       `json:"queryLogging,omitempty"`
	ServeStale                  *LocalDNSServeStale         `json:"serveStale,omitempty"`
	ServeStaleDurationInSeconds *int64                      `json:"serveStaleDurationInSeconds,omitempty"`
}
