// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type RequestOptions struct {
	// ContentType is the content type of the request and should include the charset
	ContentType string

	// ExpectedStatusCodes is a slice of HTTP response codes considered valid for this request
	ExpectedStatusCodes []int

	// HttpMethod is the capitalized method verb for this request
	HttpMethod string

	// OptionsObject is used for dynamically modifying the request at runtime
	OptionsObject Options

	// Pager is an optional struct for handling custom pagination for this request. OData 4.0 compliant paging
	// is already handled implicitly and does not require a custom pager.
	Pager odata.CustomPager

	// Path is the absolute URI for this request, with a leading slash.
	Path string

	// RetryFunc is an optional function to determine whether a request should be automatically retried
	RetryFunc RequestRetryFunc
}

func (ro RequestOptions) Validate() error {
	if len(ro.ExpectedStatusCodes) == 0 {
		return fmt.Errorf("missing `ExpectedStatusCodes`")
	}
	if ro.HttpMethod == "" {
		return fmt.Errorf("missing `HttpMethod`")
	}
	if ro.Path == "" {
		return fmt.Errorf("missing `Path`")
	}
	return nil
}
