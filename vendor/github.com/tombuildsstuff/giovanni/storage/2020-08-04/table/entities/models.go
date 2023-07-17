// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package entities

type MetaDataLevel string

var (
	NoMetaData      MetaDataLevel = "nometadata"
	MinimalMetaData MetaDataLevel = "minimalmetadata"
	FullMetaData    MetaDataLevel = "fullmetadata"
)
