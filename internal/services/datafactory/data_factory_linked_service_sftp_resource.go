// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory" // nolint: staticcheck
)

func resourceDataFactoryLinkedServiceSFTP() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryLinkedServiceSFTPCreate,
		Read:   resourceDataFactoryLinkedServiceSFTPRead,
		Update: resourceDataFactoryLinkedServiceSFTPUpdate,
		Delete: resourceDataFactoryLinkedServiceSFTPDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.LinkedServiceID(id)
			return err
		}, importDataFactoryLinkedService(datafactory.TypeBasicLinkedServiceTypeSftp)),

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
				ValidateFunc: validate.LinkedServiceDatasetName,
			},

			"data_factory_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: factories.ValidateFactoryID,
			},

			"authentication_type": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"host": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"port": {
				Type:     pluginsdk.TypeInt,
				Required: true,
			},

			"username": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"password": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"integration_runtime_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"parameters": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"skip_host_key_validation": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"host_key_fingerprint": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"annotations": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"additional_properties": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}
}

func resourceDataFactoryLinkedServiceSFTPCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	subscriptionId := meta.(*clients.Client).DataFactory.LinkedServiceClient.SubscriptionID
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewLinkedServiceID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, d.Get("name").(string))

	existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing Data Factory SFTP %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_data_factory_linked_service_sftp", id.ID())
	}

	sftpProperties := &datafactory.SftpServerLinkedServiceTypeProperties{
		Host:               d.Get("host").(string),
		Port:               d.Get("port").(int),
		AuthenticationType: datafactory.SftpAuthenticationType(d.Get("authentication_type").(string)),
		UserName:           d.Get("username").(string),
		Password: pointer.To(datafactory.SecureString{
			Value: pointer.To(d.Get("password").(string)),
			Type:  datafactory.TypeSecureString,
		}),
	}

	if v, ok := d.GetOk("skip_host_key_validation"); ok {
		sftpProperties.SkipHostKeyValidation = pointer.To(v.(bool))
	}

	if v, ok := d.GetOk("host_key_fingerprint"); ok {
		sftpProperties.HostKeyFingerprint = pointer.To(v.(string))
	}

	sftpLinkedService := &datafactory.SftpServerLinkedService{
		SftpServerLinkedServiceTypeProperties: sftpProperties,
		Type:                                  datafactory.TypeBasicLinkedServiceTypeSftp,
	}

	if v, ok := d.GetOk("description"); ok {
		sftpLinkedService.Description = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("parameters"); ok {
		sftpLinkedService.Parameters = expandLinkedServiceParameters(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("integration_runtime_name"); ok {
		sftpLinkedService.ConnectVia = expandDataFactoryLinkedServiceIntegrationRuntime(v.(string))
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		sftpLinkedService.AdditionalProperties = v.(map[string]interface{})
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotations := v.([]interface{})
		sftpLinkedService.Annotations = &annotations
	}

	linkedService := datafactory.LinkedServiceResource{
		Properties: sftpLinkedService,
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, linkedService, ""); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDataFactoryLinkedServiceSFTPRead(d, meta)
}

func resourceDataFactoryLinkedServiceSFTPRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LinkedServiceID(d.Id())
	if err != nil {
		return err
	}

	dataFactoryId := factories.NewFactoryID(id.SubscriptionId, id.ResourceGroup, id.FactoryName)

	resp, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Data Factory SFTP %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("data_factory_id", dataFactoryId.ID())

	if resp.Properties != nil {
		sftp, ok := resp.Properties.AsSftpServerLinkedService()
		if !ok {
			return fmt.Errorf("classifying Data Factory Linked Service SFTP %s: Expected: %q Received: %q", id, datafactory.TypeBasicLinkedServiceTypeSftp, pointer.From(resp.Type))
		}

		d.Set("authentication_type", sftp.AuthenticationType)
		d.Set("username", sftp.UserName)
		d.Set("port", sftp.Port)
		d.Set("host", sftp.Host)

		d.Set("additional_properties", sftp.AdditionalProperties)
		d.Set("description", sftp.Description)

		annotations := flattenDataFactoryAnnotations(sftp.Annotations)
		if err := d.Set("annotations", annotations); err != nil {
			return fmt.Errorf("setting `annotations`: %+v", err)
		}

		parameters := flattenLinkedServiceParameters(sftp.Parameters)
		if err := d.Set("parameters", parameters); err != nil {
			return fmt.Errorf("setting `parameters`: %+v", err)
		}

		if connectVia := sftp.ConnectVia; connectVia != nil {
			if connectVia.ReferenceName != nil {
				d.Set("integration_runtime_name", connectVia.ReferenceName)
			}
		}

		if props := sftp.SftpServerLinkedServiceTypeProperties; props != nil {
			if skipHostKeyValidation := props.SkipHostKeyValidation; skipHostKeyValidation != nil {
				d.Set("skip_host_key_validation", skipHostKeyValidation.(bool))
			}

			if hostKeyFingerprint := props.HostKeyFingerprint; hostKeyFingerprint != nil {
				d.Set("host_key_fingerprint", hostKeyFingerprint)
			}
		}
	}

	return nil
}

func resourceDataFactoryLinkedServiceSFTPUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LinkedServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if resp.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` is nil", id)
	}

	sftp, ok := resp.Properties.AsSftpServerLinkedService()
	if !ok {
		return fmt.Errorf("classifying Data Factory Linked Service SFTP %s: Expected: %q Received: %q", id, datafactory.TypeBasicLinkedServiceTypeSftp, pointer.From(resp.Type))
	}

	if d.HasChange("authentication_type") {
		sftp.AuthenticationType = datafactory.SftpAuthenticationType(d.Get("authentication_type").(string))
	}

	if d.HasChange("host") {
		sftp.Host = pointer.To(d.Get("host").(string))
	}

	if d.HasChange("port") {
		sftp.Port = pointer.To(d.Get("port").(string))
	}

	if d.HasChange("username") {
		sftp.UserName = pointer.To(d.Get("username").(string))
	}

	if d.HasChange("password") {
		sftp.Password = datafactory.SecureString{
			Value: pointer.To(d.Get("password").(string)),
			Type:  datafactory.TypeSecureString,
		}
	}

	if d.HasChange("skip_host_key_validation") {
		sftp.SkipHostKeyValidation = pointer.To(d.Get("skip_host_key_validation").(bool))
	}

	if d.HasChange("host_key_fingerprint") {
		sftp.HostKeyFingerprint = pointer.To(d.Get("host_key_fingerprint").(string))
	}

	if d.HasChange("description") {
		sftp.Description = pointer.To(d.Get("description").(string))
	}

	if d.HasChange("parameters") {
		sftp.Parameters = expandLinkedServiceParameters(d.Get("parameters").(map[string]interface{}))
	}

	if d.HasChange("integration_runtime_name") {
		sftp.ConnectVia = expandDataFactoryLinkedServiceIntegrationRuntime(d.Get("integration_runtime_name").(string))
	}

	if d.HasChange("additional_properties") {
		sftp.AdditionalProperties = d.Get("additional_properties").(map[string]interface{})
	}

	if d.HasChange("annotations") {
		sftp.Annotations = pointer.To(d.Get("annotations").([]interface{}))
	}

	linkedService := datafactory.LinkedServiceResource{
		Properties: sftp,
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, linkedService, ""); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceDataFactoryLinkedServiceSFTPRead(d, meta)
}

func resourceDataFactoryLinkedServiceSFTPDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LinkedServiceID(d.Id())
	if err != nil {
		return err
	}

	response, err := client.Delete(ctx, id.ResourceGroup, id.FactoryName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(response) {
			return fmt.Errorf("deleting Data Factory SFTP %s: %+v", *id, err)
		}
	}

	return nil
}
