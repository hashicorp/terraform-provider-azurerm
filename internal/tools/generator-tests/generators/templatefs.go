// Copyright IBM Corp. 2023, 2025
// SPDX-License-Identifier: MPL-2.0

package generators

import (
	"embed"
)

//go:embed templates/*
var Templatedir embed.FS
