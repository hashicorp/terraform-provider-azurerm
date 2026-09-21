// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package preflightvalidation

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceValidationRequestResource struct {
	ApiVersion string `json:"apiVersion"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Properties any    `json:"properties"`
}
