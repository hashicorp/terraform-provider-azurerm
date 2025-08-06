// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logging

// Environment variables.
const (
	// EnvTfLogSdkFramework is an environment variable that sets the logging
	// level of SDK framework loggers. Infers root SDK logging level, if
	// unset.
	EnvTfLogSdkFramework = "TF_LOG_SDK_FRAMEWORK"
)
