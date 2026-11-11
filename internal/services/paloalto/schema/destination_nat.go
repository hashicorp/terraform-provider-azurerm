// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"strconv"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/firewallresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type DestinationNAT struct {
	Name                  string                          `tfschema:"name"`
	Protocol              string                          `tfschema:"protocol"`
	FrontendConfiguration []FrontendEndpointConfiguration `tfschema:"frontend_config"`
	BackendConfiguration  []BackendEndpointConfiguration  `tfschema:"backend_config"`
}

type FrontendEndpointConfiguration struct {
	PublicIPID string `tfschema:"public_ip_address_id"`
	Port       int64  `tfschema:"port"`
}

type BackendEndpointConfiguration struct {
	PublicIP string `tfschema:"public_ip_address"`
	Port     int64  `tfschema:"port"`
}

// DestinationNATSchema returns the schema for a Palo Alto NGFW Front End Settings
func DestinationNATSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validate.DestinationNATName,
				},

				"protocol": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(firewallresources.PossibleValuesForProtocolType(), false),
				},

				"backend_config": BackendEndpointSchema(),

				"frontend_config": FrontendEndpointSchema(),
			},
		},
	}
}

func FrontendEndpointSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		MaxItems: 1,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"public_ip_address_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: commonids.ValidatePublicIPAddressID,
				},

				"port": {
					Type:         pluginsdk.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(1, 65535),
				},
			},
		},
	}
}

func BackendEndpointSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		MaxItems: 1,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"public_ip_address": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.IsIPv4Address,
				},

				"port": {
					Type:         pluginsdk.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(1, 65535),
				},
			},
		},
	}
}

func ExpandDestinationNAT(input []DestinationNAT) *[]firewallresources.FrontendSetting {
	fes := make([]firewallresources.FrontendSetting, 0)
	for _, v := range input {
		fe := firewallresources.FrontendSetting{
			Name:                  v.Name,
			Protocol:              firewallresources.ProtocolType(v.Protocol),
			BackendConfiguration:  firewallresources.EndpointConfiguration{},
			FrontendConfiguration: firewallresources.EndpointConfiguration{},
		}

		if len(v.FrontendConfiguration) > 0 {
			fec := v.FrontendConfiguration[0]
			fe.FrontendConfiguration = firewallresources.EndpointConfiguration{
				Address: firewallresources.IPAddress{
					ResourceId: pointer.To(fec.PublicIPID),
				},
				Port: strconv.FormatInt(fec.Port, 10),
			}
		}

		if len(v.BackendConfiguration) > 0 {
			bec := v.BackendConfiguration[0]
			fe.BackendConfiguration = firewallresources.EndpointConfiguration{
				Address: firewallresources.IPAddress{
					Address: pointer.To(bec.PublicIP),
				},
				Port: strconv.FormatInt(bec.Port, 10),
			}
		}

		fes = append(fes, fe)
	}

	return &fes
}

func FlattenDestinationNAT(input *[]firewallresources.FrontendSetting) []DestinationNAT {
	result := make([]DestinationNAT, 0)
	if feSettings := pointer.From(input); len(feSettings) > 0 {
		for _, v := range feSettings {
			bePort, _ := strconv.ParseInt(v.BackendConfiguration.Port, 10, 0)
			fePort, _ := strconv.ParseInt(v.FrontendConfiguration.Port, 10, 0)
			fe := DestinationNAT{
				Name:     v.Name,
				Protocol: string(v.Protocol),
				BackendConfiguration: []BackendEndpointConfiguration{{
					PublicIP: pointer.From(v.BackendConfiguration.Address.Address),
					Port:     bePort,
				}},
				FrontendConfiguration: []FrontendEndpointConfiguration{{
					PublicIPID: pointer.From(v.FrontendConfiguration.Address.ResourceId),
					Port:       fePort,
				}},
			}

			result = append(result, fe)
		}
	}
	return result
}
