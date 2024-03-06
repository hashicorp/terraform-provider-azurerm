// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package orbital

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/orbital/2022-11-01/contactprofile"
	"github.com/hashicorp/go-azure-sdk/resource-manager/orbital/2022-11-01/spacecraft"
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

type ContactProfileLinkModel struct {
	Polarization string                       `tfschema:"polarization"`
	Direction    string                       `tfschema:"direction"`
	Name         string                       `tfschema:"name"`
	Channels     []ContactProfileChannelModel `tfschema:"channels"`
}

type ContactProfileChannelModel struct {
	BandwidthMhz              float64         `tfschema:"bandwidth_mhz"`
	CenterFrequencyMhz        float64         `tfschema:"center_frequency_mhz"`
	EndPoint                  []EndPointModel `tfschema:"end_point"`
	Name                      string          `tfschema:"name"`
	ModulationConfiguration   string          `tfschema:"modulation_configuration"`
	DemodulationConfiguration string          `tfschema:"demodulation_configuration"`
}

type EndPointModel struct {
	EndPointName string `tfschema:"end_point_name"`
	IpAddress    string `tfschema:"ip_address"`
	Port         string `tfschema:"port"`
	Protocol     string `tfschema:"protocol"`
}

func ContactProfileLinkSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		ForceNew: true,
		MinItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*schema.Schema{
				"channels": ChannelSchema(),

				"direction": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(contactprofile.DirectionDownlink),
						string(contactprofile.DirectionUplink),
					}, true),
				},

				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"polarization": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(contactprofile.PolarizationLHCP),
						string(contactprofile.PolarizationRHCP),
						string(contactprofile.PolarizationLinearVertical),
						string(contactprofile.PolarizationLinearHorizontal),
					}, false),
				},
			},
		},
	}
}

func ChannelSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MinItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

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

				"end_point": EndPointSchema(),

				"modulation_configuration": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"demodulation_configuration": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func EndPointSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Required: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*schema.Schema{
				"end_point_name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"ip_address": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsIPAddress,
				},

				"port": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},

				"protocol": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(contactprofile.ProtocolTCP),
						string(contactprofile.ProtocolUDP),
					}, false),
				},
			},
		},
	}
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

func flattenContactProfileLinks(input []contactprofile.ContactProfileLink) ([]ContactProfileLinkModel, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("links are missing")
	}

	var contactProfileLinkModel []ContactProfileLinkModel
	for _, v := range input {
		linkModel := ContactProfileLinkModel{
			Polarization: string(v.Polarization),
			Direction:    string(v.Direction),
			Name:         v.Name,
		}
		contactProfileChannel, err := flattenContactProfileChannel(v.Channels)
		if err != nil {
			return nil, err
		}
		linkModel.Channels = contactProfileChannel
		contactProfileLinkModel = append(contactProfileLinkModel, linkModel)
	}

	return contactProfileLinkModel, nil
}

func flattenContactProfileChannel(input []contactprofile.ContactProfileLinkChannel) ([]ContactProfileChannelModel, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("channels are missing")
	}

	var contactProfileChannelModel []ContactProfileChannelModel
	for _, v := range input {
		channelModel := ContactProfileChannelModel{
			BandwidthMhz:              v.BandwidthMHz,
			CenterFrequencyMhz:        v.CenterFrequencyMHz,
			Name:                      v.Name,
			ModulationConfiguration:   pointer.From(v.ModulationConfiguration),
			DemodulationConfiguration: pointer.From(v.DemodulationConfiguration),
		}
		endPoint := v.EndPoint
		endPointModel := EndPointModel{
			EndPointName: endPoint.EndPointName,
			IpAddress:    endPoint.IPAddress,
			Port:         endPoint.Port,
			Protocol:     string(endPoint.Protocol),
		}
		channelModel.EndPoint = []EndPointModel{endPointModel}
		contactProfileChannelModel = append(contactProfileChannelModel, channelModel)
	}

	return contactProfileChannelModel, nil
}

func expandContactProfileLinks(input []ContactProfileLinkModel) ([]contactprofile.ContactProfileLink, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("links should be defined")
	}

	var contactProfileLink []contactprofile.ContactProfileLink
	for _, v := range input {
		link := contactprofile.ContactProfileLink{
			Direction:    contactprofile.Direction(v.Direction),
			Name:         v.Name,
			Polarization: contactprofile.Polarization(v.Polarization),
		}
		channel, err := expandContactProfileChannel(v.Channels)
		if err != nil {
			return nil, err
		}
		link.Channels = channel
		contactProfileLink = append(contactProfileLink, link)
	}

	return contactProfileLink, nil
}

func expandContactProfileChannel(input []ContactProfileChannelModel) ([]contactprofile.ContactProfileLinkChannel, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("contact profile channel should be defined")
	}

	var contactProfileChannel []contactprofile.ContactProfileLinkChannel
	for _, v := range input {
		channel := contactprofile.ContactProfileLinkChannel{
			BandwidthMHz:              v.BandwidthMhz,
			CenterFrequencyMHz:        v.CenterFrequencyMhz,
			DemodulationConfiguration: pointer.To(v.DemodulationConfiguration),
			EndPoint:                  contactprofile.EndPoint{},
			ModulationConfiguration:   pointer.To(v.ModulationConfiguration),
			Name:                      v.Name,
		}
		endPoint, err := expandEndPoint(v.EndPoint)
		if err != nil {
			return nil, err
		}
		channel.EndPoint = endPoint[0]
		contactProfileChannel = append(contactProfileChannel, channel)
	}

	return contactProfileChannel, nil
}

func expandEndPoint(input []EndPointModel) ([]contactprofile.EndPoint, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("end point should be defined")
	}

	v := input[0]
	endPoint := contactprofile.EndPoint{
		EndPointName: v.EndPointName,
		IPAddress:    v.IpAddress,
		Port:         v.Port,
		Protocol:     contactprofile.Protocol(v.Protocol),
	}

	return []contactprofile.EndPoint{endPoint}, nil
}
