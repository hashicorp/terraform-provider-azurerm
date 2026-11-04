// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package playwright_workspace

type PlaywrightWorkspaceModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	DataplaneUri      string            `tfschema:"dataplane_uri"`
	Uuid              string            `tfschema:"uuid"`
	Tags              map[string]string `tfschema:"tags"`
}

func (Resource) ModelObject() any {
	return &PlaywrightWorkspaceModel{}
}
