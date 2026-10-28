// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dns_zone

func (Resource) ModelObject() any {
	return &DnsZoneResourceModel{}
}

type DnsZoneResourceModel struct {
	Name                  string                           `tfschema:"name"`
	ResourceGroupName     string                           `tfschema:"resource_group_name"`
	NumberOfRecordSets    int64                            `tfschema:"number_of_record_sets"`
	MaxNumberOfRecordSets int64                            `tfschema:"max_number_of_record_sets"`
	NameServers           []string                         `tfschema:"name_servers"`
	SoaRecord             []DnsZoneSoaRecordResourceRecord `tfschema:"soa_record"`
	Tags                  map[string]string                `tfschema:"tags"`
}
