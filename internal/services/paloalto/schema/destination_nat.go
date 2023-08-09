package schema

import (
	"strconv"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/firewalls"
	helpersValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
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
	Port       int    `tfschema:"port"`
}

type BackendEndpointConfiguration struct {
	PublicIP string `tfschema:"public_ip_address"`
	Port     int    `tfschema:"port"`
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
					ValidateFunc: validation.StringInSlice(firewalls.PossibleValuesForProtocolType(), false),
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
					ValidateFunc: networkValidate.PublicIpAddressID,
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
					ValidateFunc: helpersValidate.IPv4Address,
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

func ExpandDestinationNAT(input []DestinationNAT) *[]firewalls.FrontendSetting {
	fes := make([]firewalls.FrontendSetting, 0)
	for _, v := range input {
		fe := firewalls.FrontendSetting{
			Name:                  v.Name,
			Protocol:              firewalls.ProtocolType(v.Protocol),
			BackendConfiguration:  firewalls.EndpointConfiguration{},
			FrontendConfiguration: firewalls.EndpointConfiguration{},
		}

		if len(v.FrontendConfiguration) > 0 {
			fec := v.FrontendConfiguration[0]
			fe.FrontendConfiguration = firewalls.EndpointConfiguration{
				Address: firewalls.IPAddress{
					ResourceId: pointer.To(fec.PublicIPID),
				},
				Port: strconv.Itoa(fec.Port),
			}
		}

		if len(v.BackendConfiguration) > 0 {
			bec := v.BackendConfiguration[0]
			fe.BackendConfiguration = firewalls.EndpointConfiguration{
				Address: firewalls.IPAddress{
					Address: pointer.To(bec.PublicIP),
				},
				Port: strconv.Itoa(bec.Port),
			}
		}

		fes = append(fes, fe)
	}

	return &fes
}

func FlattenDestinationNAT(input *[]firewalls.FrontendSetting) []DestinationNAT {
	result := make([]DestinationNAT, 0)
	if feSettings := pointer.From(input); len(feSettings) > 0 {
		for _, v := range feSettings {
			bePort, _ := strconv.Atoi(v.BackendConfiguration.Port)
			fePort, _ := strconv.Atoi(v.FrontendConfiguration.Port)
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
