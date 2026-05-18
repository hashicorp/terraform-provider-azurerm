// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/connection"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationConnectionServicePrincipal() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAutomationConnectionServicePrincipalCreate,
		Read:   resourceAutomationConnectionServicePrincipalRead,
		Update: resourceAutomationConnectionServicePrincipalUpdate,
		Delete: resourceAutomationConnectionServicePrincipalDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := connection.ParseConnectionID(id)
			return err
		}, importAutomationConnection("AzureServicePrincipal")),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ConnectionName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"automation_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AutomationAccount(),
			},

			"application_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},

			"certificate_thumbprint": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"subscription_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},

			"tenant_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAutomationConnectionServicePrincipalCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.Connection
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Automation Connection creation.")

	id := connection.NewConnectionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("automation_account_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_automation_connection_service_principal", id.ID())
	}

	parameters := connection.ConnectionCreateOrUpdateParameters{
		Name: id.ConnectionName,
		Properties: connection.ConnectionCreateOrUpdateProperties{
			Description: pointer.To(d.Get("description").(string)),
			ConnectionType: connection.ConnectionTypeAssociationProperty{
				Name: pointer.To("AzureServicePrincipal"),
			},
			FieldDefinitionValues: &map[string]string{
				"ApplicationId":         d.Get("application_id").(string),
				"CertificateThumbprint": d.Get("certificate_thumbprint").(string),
				"SubscriptionId":        d.Get("subscription_id").(string),
				"TenantId":              d.Get("tenant_id").(string),
			},
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceAutomationConnectionServicePrincipalRead(d, meta)
}

func resourceAutomationConnectionServicePrincipalUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.Connection
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Automation Connection update.")

	id, err := connection.ParseConnectionID(d.Id())
	if err != nil {
		return err
	}

	// The API represents application_id, certificate_thumbprint, subscription_id, and tenant_id
	// as a single FieldDefinitionValues map — the map is replaced as a whole on update. We GET the
	// existing values first, then overlay only the changed fields so that ignore_changes on
	// individual properties is respected.
	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving existing %s: %+v", *id, err)
	}

	fieldDefinitionValues := make(map[string]string)
	if existing.Model != nil && existing.Model.Properties != nil && existing.Model.Properties.FieldDefinitionValues != nil {
		fieldDefinitionValues = *existing.Model.Properties.FieldDefinitionValues
	}

	if d.HasChange("application_id") {
		fieldDefinitionValues["ApplicationId"] = d.Get("application_id").(string)
	}

	if d.HasChange("certificate_thumbprint") {
		fieldDefinitionValues["CertificateThumbprint"] = d.Get("certificate_thumbprint").(string)
	}

	if d.HasChange("subscription_id") {
		fieldDefinitionValues["SubscriptionId"] = d.Get("subscription_id").(string)
	}

	if d.HasChange("tenant_id") {
		fieldDefinitionValues["TenantId"] = d.Get("tenant_id").(string)
	}

	parameters := connection.ConnectionUpdateParameters{
		Name: &id.ConnectionName,
		Properties: &connection.ConnectionUpdateProperties{
			FieldDefinitionValues: &fieldDefinitionValues,
		},
	}

	if d.HasChange("description") {
		parameters.Properties.Description = pointer.To(d.Get("description").(string))
	}

	if _, err := client.Update(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceAutomationConnectionServicePrincipalRead(d, meta)
}

func resourceAutomationConnectionServicePrincipalRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.Connection
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := connection.ParseConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("read request on %s: %+v", *id, err)
	}

	d.Set("name", id.ConnectionName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("automation_account_name", id.AutomationAccountName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			description := ""
			if props.Description != nil {
				description = *props.Description
			}
			d.Set("description", description)

			if props.FieldDefinitionValues != nil {
				fieldDefinitionValues := *props.FieldDefinitionValues
				if v, ok := fieldDefinitionValues["ApplicationId"]; ok {
					d.Set("application_id", v)
				}
				if v, ok := fieldDefinitionValues["CertificateThumbprint"]; ok {
					d.Set("certificate_thumbprint", v)
				}
				if v, ok := fieldDefinitionValues["SubscriptionId"]; ok {
					d.Set("subscription_id", v)
				}
				if v, ok := fieldDefinitionValues["TenantId"]; ok {
					d.Set("tenant_id", v)
				}
			}
		}
	}

	return nil
}

func resourceAutomationConnectionServicePrincipalDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.Connection
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := connection.ParseConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}

		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
