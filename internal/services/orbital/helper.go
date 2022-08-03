package orbital

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/orbital/2022-03-01/spacecraft"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SpacecraftLinkModel struct {
	BandwidthMhz       float64 `tfschema:"bandwidth_mhz"`
	CenterFrequencyMhz float64 `tfschema:"center_frequency_mhz"`
	Direction          string  `tfschema:"direction"`
	Polarization       string  `tfschema:"polarization"`
	Name               string  `tfschema:"name"`
}

type EndPointModel struct {
	EndPointName string `tfschema:"end_point_name"`
	IpAddress    string `tfschema:"ip_address"`
	Port         string `tfschema:"port"`
	Protocol     string `tfschema:"protocol"`
}

func SpacecraftLinkSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		ForceNew: true,
		MinItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*schema.Schema{
				"bandwidth_mhz": {
					Type:         pluginsdk.TypeFloat,
					Required:     true,
					ValidateFunc: validation.FloatAtLeast(0),
				},

				"center_frequency_mhz": {
					Type:         pluginsdk.TypeFloat,
					Required:     true,
					ValidateFunc: validation.FloatAtLeast(0),
				},

				"direction": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(spacecraft.DirectionUplink),
						string(spacecraft.DirectionDownlink),
					}, true),
				},

				"polarization": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(spacecraft.PolarizationLHCP),
						string(spacecraft.PolarizationLinearHorizontal),
						string(spacecraft.PolarizationRHCP),
						string(spacecraft.PolarizationLinearVertical),
					}, false),
				},

				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func expandSpacecraftLinks(input []SpacecraftLinkModel) ([]spacecraft.SpacecraftLink, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("links should be defined")
	}
	var spacecraftLink []spacecraft.SpacecraftLink
	for _, v := range input {
		link := spacecraft.SpacecraftLink{
			BandwidthMHz:       v.BandwidthMhz,
			CenterFrequencyMHz: v.CenterFrequencyMhz,
			Direction:          spacecraft.Direction(v.Direction),
			Polarization:       spacecraft.Polarization(v.Polarization),
			Name:               v.Name,
		}
		spacecraftLink = append(spacecraftLink, link)
	}

	return spacecraftLink, nil
}

func flattenSpacecraftLinks(input []spacecraft.SpacecraftLink) ([]SpacecraftLinkModel, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("links are missing")
	}

	var spacecraftLinkModel []SpacecraftLinkModel
	for _, v := range input {
		linkModel := SpacecraftLinkModel{
			BandwidthMhz:       v.BandwidthMHz,
			CenterFrequencyMhz: v.CenterFrequencyMHz,
			Direction:          string(v.Direction),
			Polarization:       string(v.Polarization),
			Name:               v.Name,
		}
		spacecraftLinkModel = append(spacecraftLinkModel, linkModel)
	}

	return spacecraftLinkModel, nil
}
