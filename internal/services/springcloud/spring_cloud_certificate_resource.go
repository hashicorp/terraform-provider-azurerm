// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package springcloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func resourceSpringCloudCertificate() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		DeprecationMessage: features.DeprecatedInFivePointOh("Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_certificate` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information."),

		Create: resourceSpringCloudCertificateCreate,
		Read:   resourceSpringCloudCertificateRead,
		Delete: resourceSpringCloudCertificateDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SpringCloudCertificateV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SpringCloudCertificateID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"service_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SpringCloudServiceName,
			},

			"certificate_content": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"key_vault_certificate_id", "certificate_content"},
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"exclude_private_key": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"key_vault_certificate_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"key_vault_certificate_id", "certificate_content"},
				ValidateFunc: keyVaultValidate.NestedItemId,
			},

			"thumbprint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSpringCloudCertificateCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.CertificatesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("service_name").(string)

	resourceId := parse.NewSpringCloudCertificateID(subscriptionId, resourceGroup, serviceName, name).ID()
	existing, err := client.Get(ctx, resourceGroup, serviceName, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing Spring Cloud Service Certificate %q (Spring Cloud Service %q / Resource Group %q): %+v", name, serviceName, resourceGroup, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_spring_cloud_certificate", resourceId)
	}

	cert := appplatform.CertificateResource{}
	if value, ok := d.GetOk("key_vault_certificate_id"); ok {
		keyVaultCertificateId, err := keyVaultParse.ParseNestedItemID(value.(string))
		if err != nil {
			return err
		}
		cert.Properties = &appplatform.KeyVaultCertificateProperties{
			VaultURI:          utils.String(strings.TrimSuffix(keyVaultCertificateId.KeyVaultBaseUrl, "/")),
			KeyVaultCertName:  &keyVaultCertificateId.Name,
			ExcludePrivateKey: utils.Bool(d.Get("exclude_private_key").(bool)),
		}
	}
	if value, ok := d.GetOk("certificate_content"); ok {
		cert.Properties = &appplatform.ContentCertificateProperties{
			Content: utils.String(value.(string)),
		}
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, name, cert)
	if err != nil {
		return fmt.Errorf("creating Spring Cloud Certificate %q (Spring Cloud Service %q / Resource Group %q): %+v", name, serviceName, resourceGroup, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %q(Spring Cloud Service %q / Resource Group %q): %+v", name, serviceName, resourceGroup, err)
	}

	d.SetId(resourceId)
	return resourceSpringCloudCertificateRead(d, meta)
}

func resourceSpringCloudCertificateRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.CertificatesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudCertificateID(d.Id())
	if err != nil {
		return err
	}
	resp, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.CertificateName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Spring Cloud Certificate %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Spring Cloud Certificate %q (Spring Cloud Service %q / Resource Group %q): %+v", id.CertificateName, id.SpringName, id.ResourceGroup, err)
	}

	d.Set("name", id.CertificateName)
	d.Set("service_name", id.SpringName)
	d.Set("resource_group_name", id.ResourceGroup)

	if props, ok := resp.Properties.AsKeyVaultCertificateProperties(); ok && props != nil {
		d.Set("thumbprint", props.Thumbprint)
		d.Set("exclude_private_key", props.ExcludePrivateKey)
	}
	if props, ok := resp.Properties.AsContentCertificateProperties(); ok && props != nil {
		d.Set("thumbprint", props.Thumbprint)
	}

	return nil
}

func resourceSpringCloudCertificateDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.CertificatesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudCertificateID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.SpringName, id.CertificateName)
	if err != nil {
		return fmt.Errorf("deleting Spring Cloud Certificate %q (Spring Cloud Service %q / Resource Group %q): %+v", id.CertificateName, id.SpringName, id.ResourceGroup, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %q: %+v", id, err)
	}

	return nil
}
