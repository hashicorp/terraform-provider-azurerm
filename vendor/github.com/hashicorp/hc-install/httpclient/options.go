// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package httpclient

import (
	"github.com/hashicorp/go-retryablehttp"
)

type Option func(c *retryablehttp.Client)
