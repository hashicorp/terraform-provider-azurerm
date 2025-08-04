// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package generators

import (
	"embed"
)

//go:embed templates/*
var Templatedir embed.FS
