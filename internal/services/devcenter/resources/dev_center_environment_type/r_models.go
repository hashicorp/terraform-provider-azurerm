// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center_environment_type

func (r Resource) ModelObject() any {
	return &DevCenterEnvironmentTypeResourceModel{}
}

type DevCenterEnvironmentTypeResourceModel struct {
	Name        string            `tfschema:"name"`
	DevCenterId string            `tfschema:"dev_center_id"`
	Tags        map[string]string `tfschema:"tags"`
}
