// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logging

// Environment variables.
const (
	// EnvTfLogSdkMux is an environment variable that sets the logging level
	// of the mux logger. Infers root SDK logging level, if unset.
	EnvTfLogSdkMux = "TF_LOG_SDK_MUX"
)
