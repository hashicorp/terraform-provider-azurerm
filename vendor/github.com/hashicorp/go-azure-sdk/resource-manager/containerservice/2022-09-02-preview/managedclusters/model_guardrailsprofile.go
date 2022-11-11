package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GuardrailsProfile struct {
	ExcludedNamespaces       *[]string `json:"excludedNamespaces,omitempty"`
	Level                    Level     `json:"level"`
	SystemExcludedNamespaces *[]string `json:"systemExcludedNamespaces,omitempty"`
	Version                  string    `json:"version"`
}
