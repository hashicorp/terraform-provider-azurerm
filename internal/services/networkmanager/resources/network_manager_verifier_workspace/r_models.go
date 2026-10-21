// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_verifier_workspace

func (Resource) ModelObject() any {
	return &ManagerVerifierWorkspaceResourceModel{}
}

type ManagerVerifierWorkspaceResourceModel struct {
	Description      string            `tfschema:"description"`
	Location         string            `tfschema:"location"`
	Name             string            `tfschema:"name"`
	NetworkManagerId string            `tfschema:"network_manager_id"`
	Tags             map[string]string `tfschema:"tags"`
}
