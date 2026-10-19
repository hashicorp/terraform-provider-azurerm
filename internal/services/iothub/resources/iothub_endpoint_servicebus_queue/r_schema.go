// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_endpoint_servicebus_queue

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	iothubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	devices "github.com/jackofallops/kermit/sdk/iothub/2022-04-30-preview/iothub" // azignore:AZG010 - package name does not match its path
)

func resourceIothubEndpointServicebusQueue() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: iothubValidate.IoTHubEndpointName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		// lintignore: S013
		"iothub_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: iothubValidate.IotHubID,
		},

		"authentication_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(devices.AuthenticationTypeKeyBased),
			ValidateFunc: validation.StringInEnumSlice(devices.PossibleAuthenticationTypeValues(), false),
		},

		"identity_id": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  commonids.ValidateUserAssignedIdentityID,
			ConflictsWith: []string{"connection_string"},
		},

		"endpoint_uri": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"entity_path"},
			ExactlyOneOf: []string{"endpoint_uri", "connection_string"},
		},

		"entity_path": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.QueueName(),
			RequiredWith: []string{"endpoint_uri"},
		},

		"connection_string": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
				sharedAccessKeyRegex := regexp.MustCompile("SharedAccessKey=[^;]+")
				sbProtocolRegex := regexp.MustCompile("sb://([^:]+)(:5671)?/;")

				maskedNew := sbProtocolRegex.ReplaceAllString(new, "sb://$1:5671/;")
				maskedNew = sharedAccessKeyRegex.ReplaceAllString(maskedNew, "SharedAccessKey=****")
				return (new == d.Get(k).(string)) && (maskedNew == old)
			},
			Sensitive:     true,
			ConflictsWith: []string{"identity_id"},
			ExactlyOneOf:  []string{"endpoint_uri", "connection_string"},
		},

		"subscription_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C : required since this property would always be set even if it isn't specified in the tf config, otherwise it would cause a diff and break existing users
			Computed:     true,
			ValidateFunc: validation.IsUUID,
		},
	}
}
