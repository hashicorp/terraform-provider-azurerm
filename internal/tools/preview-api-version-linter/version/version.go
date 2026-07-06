// Copyright IBM Corp. 2023, 2026
// SPDX-License-Identifier: MPL-2.0

package version

type Version struct {
	Module  string // example: github.com/hashicorp/go-azure-sdk/resource-manager
	Service string // example: compute
	Version string // example: 2020-02-02-preview
}
