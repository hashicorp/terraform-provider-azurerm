// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_virtual_machine_group

type MsSqlVirtualMachineGroupModel struct {
	Name          string `tfschema:"name"`
	ResourceGroup string `tfschema:"resource_group_name"`
	Location      string `tfschema:"location"`

	SqlImageOffer     string              `tfschema:"sql_image_offer"`
	SqlImageSku       string              `tfschema:"sql_image_sku"`
	WsfcDomainProfile []WsfcDomainProfile `tfschema:"wsfc_domain_profile"`
	Tags              map[string]string   `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &MsSqlVirtualMachineGroupModel{}
}
