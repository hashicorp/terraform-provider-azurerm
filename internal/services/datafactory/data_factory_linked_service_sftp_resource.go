// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/linkedservices"
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
		Create: resourceDataFactoryLinkedServiceSFTPCreateUpdate,
		Read:   resourceDataFactoryLinkedServiceSFTPRead,
		Update: resourceDataFactoryLinkedServiceSFTPCreateUpdate,
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
				ValidateFunc: validation.StringInSlice(linkedservices.PossibleValuesForSftpAuthenticationType(), false),
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
				Description:  "The user who has access to the SFTP server",
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"additional_properties": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"annotations": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"host_key_fingerprint": {
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

			"private_key_passphrase": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				Description:   "Specify the pass phrase or password to decrypt the private key if the key file or the key content is protected by a pass phrase",
				ConflictsWith: []string{"password"},
			},

			"password": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				Sensitive:     true,
				ValidateFunc:  validation.StringIsNotEmpty,
				ConflictsWith: []string{"private_key_content", "private_key_path", "private_key_passphrase"},
				AtLeastOneOf:  []string{"password", "private_key_content", "private_key_path"},
			},

			"private_key_content": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				Sensitive:     true,
				Description:   "SSH private key content in OpenSSH format",
				ValidateFunc:  validate.SSHPrivateKey,
				ConflictsWith: []string{"private_key_path", "password"},
				AtLeastOneOf:  []string{"password", "private_key_content", "private_key_path"},
			},

			"private_key_path": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				Description:   "Specify the absolute path to the private key file that the integration runtime can access. This applies only when using a self-hosted integration runtime as opposed to the default Azure-hosted runtime, as indicated by providing a value for `integration_runtime_name`.",
				ConflictsWith: []string{"private_key_content", "password"},
				AtLeastOneOf:  []string{"password", "private_key_content", "private_key_path"},
			},

			"skip_host_key_validation": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceDataFactoryLinkedServiceSFTPCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	subscriptionId := meta.(*clients.Client).DataFactory.LinkedServiceClient.SubscriptionID
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewLinkedServiceID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Data Factory SFTP %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_data_factory_linked_service_sftp", id.ID())
		}
	}

	authenticationType := d.Get("authentication_type").(string)

	host := d.Get("host").(string)
	port := d.Get("port").(int)
	username := d.Get("username").(string)

	//passphrase := datafactory.SecureString{
	//	Value: pointer.To(d.Get("private_key_passphrase").(string)),
	//	Type:  datafactory.TypeSecureString,
	//}

	sftpProperties := &datafactory.SftpServerLinkedServiceTypeProperties{
		Host:               pointer.To(host),
		Port:               port,
		AuthenticationType: datafactory.SftpAuthenticationType(authenticationType),
		UserName:           pointer.To(username),
		//Password:           &passwordSecureString,
		//PrivateKeyContent:  &privateKeyContent,
		//PrivateKeyPath:     d.Get("private_key_path").(string),
		//PassPhrase:         &passphrase,
	}

	if v, ok := d.GetOk("password"); ok {
		passwordSecureString := datafactory.SecureString{
			Value: pointer.To(v.(string)),
			Type:  datafactory.TypeSecureString,
		}
		sftpProperties.Password = &passwordSecureString
	}

	if v, ok := d.GetOk("private_key_content"); ok {
		privateKeyContent := datafactory.SecureString{
			Value: pointer.To(v.(string)),
			Type:  datafactory.TypeSecureString,
		}
		sftpProperties.PrivateKeyContent = &privateKeyContent
	}

	if v, ok := d.GetOk("private_key_passphrase"); ok {
		passphrase := datafactory.SecureString{
			Value: pointer.To(v.(string)),
			Type:  datafactory.TypeSecureString,
		}
		sftpProperties.PrivateKeyContent = &passphrase
	}

	if v, ok := d.GetOk("private_key_path"); ok {
		sftpProperties.PrivateKeyPath = pointer.To(v.(string))
	}

	sftpProperties.SkipHostKeyValidation = d.Get("skip_host_key_validation").(bool)
	sftpProperties.HostKeyFingerprint = d.Get("host_key_fingerprint").(string)
	description := d.Get("description").(string)

	sftpLinkedService := &datafactory.SftpServerLinkedService{
		Description:                           &description,
		SftpServerLinkedServiceTypeProperties: sftpProperties,
		Type:                                  datafactory.TypeBasicLinkedServiceTypeSftp,
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
		return fmt.Errorf("creating/updating Data Factory SFTP Anonymous %s: %+v", id, err)
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

	d.Set("name", resp.Name)
	d.Set("data_factory_id", dataFactoryId.ID())

	sftp, ok := resp.Properties.AsSftpServerLinkedService()
	if !ok {
		return fmt.Errorf("classifying Data Factory Linked Service SFTP %q (Data Factory %q / Resource Group %q): Expected: %q Received: %q", id.Name, id.FactoryName, id.ResourceGroup, datafactory.TypeBasicLinkedServiceTypeSftp, *resp.Type)
	}

	d.Set("authentication_type", sftp.AuthenticationType)
	d.Set("username", sftp.UserName)
	d.Set("port", sftp.Port)
	d.Set("host", sftp.Host)
	d.Set("private_key_path", sftp.PrivateKeyPath)

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

	return nil
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
