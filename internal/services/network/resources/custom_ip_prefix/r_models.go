// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custom_ip_prefix

type CustomIpPrefixModel struct {
	CIDR                        string         `tfschema:"cidr"`
	CommissioningEnabled        bool           `tfschema:"commissioning_enabled"`
	InternetAdvertisingDisabled bool           `tfschema:"internet_advertising_disabled"`
	Location                    string         `tfschema:"location"`
	Name                        string         `tfschema:"name"`
	ParentCustomIPPrefixID      string         `tfschema:"parent_custom_ip_prefix_id"`
	ROAValidityEndDate          string         `tfschema:"roa_validity_end_date"`
	ResourceGroupName           string         `tfschema:"resource_group_name"`
	Tags                        map[string]any `tfschema:"tags"`
	WANValidationSignedMessage  string         `tfschema:"wan_validation_signed_message"`
	Zones                       []string       `tfschema:"zones"`
}

func (Resource) ModelObject() any {
	return &CustomIpPrefixModel{}
}
