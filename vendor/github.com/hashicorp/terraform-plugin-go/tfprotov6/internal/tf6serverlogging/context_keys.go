// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf6serverlogging

// Context key types.
// Reference: https://staticcheck.io/docs/checks/#SA1029

// ContextKeyDownstreamRequestStartTime is a context.Context key to store the
// time.Time when the server began a downstream request.
type ContextKeyDownstreamRequestStartTime struct{}
