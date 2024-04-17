// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/certificate"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementCertificate() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementCertificateCreateUpdate,
		Read:   resourceApiManagementCertificateRead,
		Update: resourceApiManagementCertificateCreateUpdate,
		Delete: resourceApiManagementCertificateDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := certificate.ParseCertificateID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": schemaz.SchemaApiManagementChildName(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"api_management_name": schemaz.SchemaApiManagementName(),

			"data": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				Sensitive:     true,
				ValidateFunc:  validation.StringIsBase64,
				AtLeastOneOf:  []string{"data", "key_vault_secret_id"},
				ConflictsWith: []string{"key_vault_secret_id", "key_vault_identity_client_id"},
			},

			"password": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				RequiredWith: []string{"data"},
			},

			"key_vault_secret_id": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ValidateFunc:  keyVaultValidate.NestedItemIdWithOptionalVersion,
				AtLeastOneOf:  []string{"data", "key_vault_secret_id"},
				ConflictsWith: []string{"data", "password"},
			},

			"key_vault_identity_client_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
				RequiredWith: []string{"key_vault_secret_id"},
			},

			"expiration": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"subject": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"thumbprint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceApiManagementCertificateCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.CertificatesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := certificate.NewCertificateID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("name").(string))

	data := d.Get("data").(string)
	password := d.Get("password").(string)
	keyVaultSecretId := d.Get("key_vault_secret_id").(string)
	keyVaultIdentity := d.Get("key_vault_identity_client_id").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_api_management_certificate", id.ID())
		}
	}

	parameters := certificate.CertificateCreateOrUpdateParameters{
		Properties: &certificate.CertificateCreateOrUpdateProperties{},
	}

	if keyVaultSecretId != "" {
		parsedSecretId, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(keyVaultSecretId)
		if err != nil {
			return err
		}

		parameters.Properties.KeyVault = &certificate.KeyVaultContractCreateProperties{
			SecretIdentifier: pointer.To(parsedSecretId.ID()),
		}

		if keyVaultIdentity != "" {
			parameters.Properties.KeyVault.IdentityClientId = pointer.To(keyVaultIdentity)
		}
	}

	if data != "" {
		parameters.Properties.Data = pointer.To(data)
		parameters.Properties.Password = pointer.To(password)
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters, certificate.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating or updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementCertificateRead(d, meta)
}

func resourceApiManagementCertificateRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.CertificatesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := certificate.ParseCertificateID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request for %s: %+v", *id, err)
	}

	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("api_management_name", id.ServiceName)

	if model := resp.Model; model != nil {
		d.Set("name", pointer.From(model.Name))
		if props := model.Properties; props != nil {
			d.Set("expiration", props.ExpirationDate)
			d.Set("subject", props.Subject)
			d.Set("thumbprint", props.Thumbprint)

			if keyvault := props.KeyVault; keyvault != nil {
				d.Set("key_vault_secret_id", pointer.From(keyvault.SecretIdentifier))
				d.Set("key_vault_identity_client_id", pointer.From(keyvault.IdentityClientId))
			}
		}
	}

	return nil
}

func resourceApiManagementCertificateDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.CertificatesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := certificate.ParseCertificateID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, *id, certificate.DeleteOperationOptions{}); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
