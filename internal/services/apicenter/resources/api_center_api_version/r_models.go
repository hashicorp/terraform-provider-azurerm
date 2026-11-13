// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_api_version

type ApiCenterApiVersionModel struct {
	ApiId          string `tfschema:"api_id"`
	Name           string `tfschema:"name"`
	LifecycleStage string `tfschema:"lifecycle_stage"`
	Title          string `tfschema:"display_name"`
}
