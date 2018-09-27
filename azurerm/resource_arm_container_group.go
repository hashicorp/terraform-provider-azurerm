package azurerm

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/2018-04-01/containerinstance"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmContainerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmContainerGroupCreate,
		Read:   resourceArmContainerGroupRead,
		Delete: resourceArmContainerGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"ip_address_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "Public",
				ForceNew:         true,
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
				ValidateFunc: validation.StringInSlice([]string{
					"Public",
				}, true),
			},

			"os_type": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
				ValidateFunc: validation.StringInSlice([]string{
					"windows",
					"linux",
				}, true),
			},

			"image_registry_credential": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
							ForceNew:     true,
						},

						"username": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
							ForceNew:     true,
						},

						"password": {
							Type:         schema.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validation.NoZeroValues,
							ForceNew:     true,
						},
					},
				},
			},

			"tags": tagsForceNewSchema(),

			"restart_policy": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          string(containerinstance.Always),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
				ValidateFunc: validation.StringInSlice([]string{
					string(containerinstance.Always),
					string(containerinstance.Never),
					string(containerinstance.OnFailure),
				}, true),
			},

			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"dns_name_label": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"volume": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{

					Schema: map[string]*schema.Schema{
						"empty_dir": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem:     &schema.Resource{},
						},

						"secret": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},

									"data": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},

						"git_repo": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"repository": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"directory": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},

						"azure_share": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"share_name": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},

									"storage_account_name": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},

									"storage_account_key": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},

						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			"container": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"image": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"cpu": {
							Type:     schema.TypeFloat,
							Required: true,
							ForceNew: true,
						},

						"memory": {
							Type:     schema.TypeFloat,
							Required: true,
							ForceNew: true,
						},

						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},

						"protocol": {
							Type:             schema.TypeString,
							Optional:         true,
							ForceNew:         true,
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
							ValidateFunc: validation.StringInSlice([]string{
								"tcp",
								"udp",
							}, true),
						},

						"environment_variables": {
							Type:     schema.TypeMap,
							Optional: true,
							ForceNew: true,
						},

						"command": {
							Type:       schema.TypeString,
							Optional:   true,
							ForceNew:   true,
							Computed:   true,
							Deprecated: "Use `commands` instead.",
						},

						"commands": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},

						"volume_mount": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"volume_name": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},

									"mount_path": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},

									"read_only": {
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: true,
										Default:  false,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceArmContainerGroupCreate(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext
	containerGroupsClient := meta.(*ArmClient).containerGroupsClient

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	OSType := d.Get("os_type").(string)
	IPAddressType := d.Get("ip_address_type").(string)
	tags := d.Get("tags").(map[string]interface{})
	restartPolicy := d.Get("restart_policy").(string)

	containers, containerGroupPorts := expandContainerGroupContainers(d)
	containerGroupVolumes, err := expandContainerVolumes(d)
	if err != nil {
		return err
	}

	containerGroup := containerinstance.ContainerGroup{
		Name:     &name,
		Location: &location,
		Tags:     expandTags(tags),
		ContainerGroupProperties: &containerinstance.ContainerGroupProperties{
			Containers:    containers,
			RestartPolicy: containerinstance.ContainerGroupRestartPolicy(restartPolicy),
			IPAddress: &containerinstance.IPAddress{
				Type:  &IPAddressType,
				Ports: containerGroupPorts,
			},
			OsType:                   containerinstance.OperatingSystemTypes(OSType),
			Volumes:                  containerGroupVolumes,
			ImageRegistryCredentials: expandContainerImageRegistryCredentials(d),
		},
	}

	if dnsNameLabel := d.Get("dns_name_label").(string); dnsNameLabel != "" {
		containerGroup.ContainerGroupProperties.IPAddress.DNSNameLabel = &dnsNameLabel
	}

	_, err = containerGroupsClient.CreateOrUpdate(ctx, resGroup, name, containerGroup)
	if err != nil {
		return err
	}

	read, err := containerGroupsClient.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read container group %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmContainerGroupRead(d, meta)
}

func resourceArmContainerGroupRead(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext
	client := meta.(*ArmClient).containerGroupsClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["containerGroups"]

	resp, err := client.Get(ctx, resourceGroup, name)

	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Container Group %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if props := resp.ContainerGroupProperties; props != nil {
		containerConfigs := flattenContainerGroupContainers(d, resp.Containers, props.IPAddress.Ports, props.Volumes)
		if err := d.Set("container", containerConfigs); err != nil {
			return fmt.Errorf("Error setting `container`: %+v", err)
		}

		if err := d.Set("image_registry_credential", flattenContainerImageRegistryCredentials(d, props.ImageRegistryCredentials)); err != nil {
			return fmt.Errorf("Error setting `capabilities`: %+v", err)
		}

		if address := props.IPAddress; address != nil {
			d.Set("ip_address_type", address.Type)
			d.Set("ip_address", address.IP)
			d.Set("dns_name_label", address.DNSNameLabel)
			d.Set("fqdn", address.Fqdn)
		}

		d.Set("restart_policy", string(props.RestartPolicy))
		d.Set("os_type", string(props.OsType))

		if err := d.Set("volume", flattenContainerVolumes(d, props.Volumes)); err != nil {
			return fmt.Errorf("Error setting `volume`: %+v", err)
		}
	}
	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmContainerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext
	client := meta.(*ArmClient).containerGroupsClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["containerGroups"]

	resp, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error deleting Container Group %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func flattenContainerGroupContainers(d *schema.ResourceData, containers *[]containerinstance.Container, containerGroupPorts *[]containerinstance.Port, containerGroupVolumes *[]containerinstance.Volume) []interface{} {

	containerConfigs := make([]interface{}, 0, len(*containers))
	for _, container := range *containers {
		containerConfig := make(map[string]interface{})
		containerConfig["name"] = *container.Name
		containerConfig["image"] = *container.Image

		if resources := container.Resources; resources != nil {
			if resourceRequests := resources.Requests; resourceRequests != nil {
				containerConfig["cpu"] = *resourceRequests.CPU
				containerConfig["memory"] = *resourceRequests.MemoryInGB
			}
		}

		if len(*container.Ports) > 0 {
			containerPort := *(*container.Ports)[0].Port
			containerConfig["port"] = containerPort
			// protocol isn't returned in container config, have to search in container group ports
			protocol := ""
			if containerGroupPorts != nil {
				for _, cgPort := range *containerGroupPorts {
					if *cgPort.Port == containerPort {
						protocol = string(cgPort.Protocol)
					}
				}
			}
			if protocol != "" {
				containerConfig["protocol"] = protocol
			}
		}

		if container.EnvironmentVariables != nil {
			if len(*container.EnvironmentVariables) > 0 {
				containerConfig["environment_variables"] = flattenContainerEnvironmentVariables(container.EnvironmentVariables)
			}
		}

		commands := make([]string, 0)
		if command := container.Command; command != nil {
			containerConfig["command"] = strings.Join(*command, " ")

			for _, v := range *command {
				commands = append(commands, v)
			}
		}
		containerConfig["commands"] = commands

		if container.VolumeMounts != nil {
			containerConfig["volume_mount"] = flattenContainerVolumeMounts(container.VolumeMounts)
		}

		containerConfigs = append(containerConfigs, containerConfig)
	}

	return containerConfigs
}

func flattenContainerEnvironmentVariables(input *[]containerinstance.EnvironmentVariable) map[string]interface{} {
	output := make(map[string]interface{})
	if input == nil {
		return output
	}

	for _, envVar := range *input {
		if envVar.Name != nil && envVar.Value != nil {
			output[*envVar.Name] = *envVar.Value
		}
	}

	return output
}

func flattenContainerVolumeMounts(volumeMounts *[]containerinstance.VolumeMount) []interface{} {
	volumeConfigs := make([]interface{}, 0)

	if volumeMounts == nil {
		return volumeConfigs
	}

	for _, vm := range *volumeMounts {
		volumeConfig := make(map[string]interface{})
		if vm.Name != nil {
			volumeConfig["volume_name"] = *vm.Name
		}
		if vm.MountPath != nil {
			volumeConfig["mount_path"] = *vm.MountPath
		}
		if vm.ReadOnly != nil {
			volumeConfig["read_only"] = *vm.ReadOnly
		}
		volumeConfigs = append(volumeConfigs, volumeConfig)
	}

	return volumeConfigs
}

func flattenContainerVolumes(d *schema.ResourceData, containerGroupVolumes *[]containerinstance.Volume) []interface{} {
	volumeConfigs := make([]interface{}, 0)

	if containerGroupVolumes == nil {
		return volumeConfigs
	}

	oldVolumeConfigs := d.Get("volume").([]interface{})

	for i, volume := range *containerGroupVolumes {

		volumeConfig := make(map[string]interface{})
		volumeConfig["name"] = *volume.Name

		if volume.EmptyDir != nil {
			volumeConfig["empty_dir"] = [1]interface{}{}
			volumeConfigs = append(volumeConfigs, volumeConfig)
			continue
		}

		if volume.GitRepo != nil {
			volumeConfig["git_repo"] = [1]interface{}{
				map[string]string{
					"repository": *volume.GitRepo.Repository,
					"directory":  *volume.GitRepo.Directory,
				},
			}
			volumeConfigs = append(volumeConfigs, volumeConfig)
			continue
		}

		var oldVolumeConfig map[string]interface{}
		// Secrets aren't returned so check the old config for them.
		if len(oldVolumeConfigs) > i {
			if oldVolumeConfigCheck, exists := oldVolumeConfigs[i].(map[string]interface{}); exists {
				if oldVolumeConfigName, exists := oldVolumeConfigCheck["name"]; exists && oldVolumeConfigName == *volume.Name {
					oldVolumeConfig = oldVolumeConfigCheck
				}
			}
		}

		if volume.AzureFile != nil {
			azureShare := make(map[string]interface{})
			azureShare["share_name"] = *volume.AzureFile.ShareName
			azureShare["storage_account_name"] = *volume.AzureFile.StorageAccountName
			if oldVolumeConfig != nil {
				azureShare["storage_account_key"] = oldVolumeConfig["azure_share"].([]interface{})[0].(map[string]interface{})["storage_account_key"]
			}
			volumeConfig["azure_share"] = []interface{}{
				azureShare,
			}
			volumeConfigs = append(volumeConfigs, volumeConfig)
			continue
		}

		if oldVolumeConfig != nil {
			if _, exists := validateArmContainerGroupVolumeVolumeExists(oldVolumeConfig, "secret"); exists {
				volumeConfig["secret"] = oldVolumeConfig["secret"]
			}
		}

		volumeConfigs = append(volumeConfigs, volumeConfig)
	}

	return volumeConfigs
}

func expandContainerGroupContainers(d *schema.ResourceData) (*[]containerinstance.Container, *[]containerinstance.Port) {
	containersConfig := d.Get("container").([]interface{})
	containers := make([]containerinstance.Container, 0)
	containerGroupPorts := make([]containerinstance.Port, 0)

	for _, containerConfig := range containersConfig {
		data := containerConfig.(map[string]interface{})

		name := data["name"].(string)
		image := data["image"].(string)
		cpu := data["cpu"].(float64)
		memory := data["memory"].(float64)

		container := containerinstance.Container{
			Name: utils.String(name),
			ContainerProperties: &containerinstance.ContainerProperties{
				Image: utils.String(image),
				Resources: &containerinstance.ResourceRequirements{
					Requests: &containerinstance.ResourceRequests{
						MemoryInGB: utils.Float(memory),
						CPU:        utils.Float(cpu),
					},
				},
			},
		}

		if v, _ := data["port"]; v != 0 {
			port := int32(v.(int))

			// container port (port number)
			container.Ports = &[]containerinstance.ContainerPort{
				{
					Port: &port,
				},
			}

			// container group port (port number + protocol)
			containerGroupPort := containerinstance.Port{
				Port: &port,
			}

			if v, ok := data["protocol"]; ok {
				protocol := v.(string)
				containerGroupPort.Protocol = containerinstance.ContainerGroupNetworkProtocol(strings.ToUpper(protocol))
			}

			containerGroupPorts = append(containerGroupPorts, containerGroupPort)
		}

		if v, ok := data["environment_variables"]; ok {
			container.EnvironmentVariables = expandContainerEnvironmentVariables(v)
		}

		if v, ok := data["commands"]; ok {
			c := v.([]interface{})
			command := make([]string, 0)
			for _, v := range c {
				command = append(command, v.(string))
			}

			container.Command = &command
		}

		if container.Command == nil {
			if v, _ := data["command"]; v != "" {
				command := strings.Split(v.(string), " ")
				container.Command = &command
			}
		}

		if v, ok := data["volume_mount"]; ok {
			volumeMounts := expandContainerVolumeMounts(v)
			container.VolumeMounts = volumeMounts
		}

		containers = append(containers, container)
	}

	return &containers, &containerGroupPorts
}

func expandContainerEnvironmentVariables(input interface{}) *[]containerinstance.EnvironmentVariable {
	envVars := input.(map[string]interface{})
	output := make([]containerinstance.EnvironmentVariable, 0)

	for k, v := range envVars {
		ev := containerinstance.EnvironmentVariable{
			Name:  utils.String(k),
			Value: utils.String(v.(string)),
		}
		output = append(output, ev)
	}

	return &output
}

func expandContainerImageRegistryCredentials(d *schema.ResourceData) *[]containerinstance.ImageRegistryCredential {
	credsRaw := d.Get("image_registry_credential").([]interface{})
	if len(credsRaw) == 0 {
		return nil
	}

	output := make([]containerinstance.ImageRegistryCredential, 0, len(credsRaw))

	for _, c := range credsRaw {
		credConfig := c.(map[string]interface{})

		output = append(output, containerinstance.ImageRegistryCredential{
			Server:   utils.String(credConfig["server"].(string)),
			Password: utils.String(credConfig["password"].(string)),
			Username: utils.String(credConfig["username"].(string)),
		})
	}

	return &output
}

func flattenContainerImageRegistryCredentials(d *schema.ResourceData, input *[]containerinstance.ImageRegistryCredential) []interface{} {
	if input == nil {
		return nil
	}
	configsOld := d.Get("image_registry_credential").([]interface{})

	output := make([]interface{}, 0)
	for i, cred := range *input {
		credConfig := make(map[string]interface{})
		if cred.Server != nil {
			credConfig["server"] = *cred.Server
		}
		if cred.Username != nil {
			credConfig["username"] = *cred.Username
		}

		if len(configsOld) > i {
			data := configsOld[i].(map[string]interface{})
			oldServer := data["server"].(string)
			if cred.Server != nil && *cred.Server == oldServer {
				if v, ok := d.GetOk(fmt.Sprintf("image_registry_credential.%d.password", i)); ok {
					credConfig["password"] = v.(string)
				}
			}
		}

		output = append(output, credConfig)
	}
	return output
}

func validateArmContainerGroupVolumeVolumeExists(volume map[string]interface{}, name string) (map[string]interface{}, bool) {
	if volumeList := volume[name]; len(volumeList.([]interface{})) > 0 {
		if volumeItem := volumeList.([]interface{})[0]; volumeItem != nil {
			return volumeItem.(map[string]interface{}), true
		}
	}
	return nil, false
}

func validateArmContainerGroupVolume(v interface{}) (errors []error) {
	volumesRaw := v.([]interface{})
	for _, volumeRaw := range volumesRaw {
		volumeConfig := volumeRaw.(map[string]interface{})

		volumeTypeCount := 0
		var azureShareExists, gitRepoExists, secretExists, emptyDirExists bool

		if _, exists := validateArmContainerGroupVolumeVolumeExists(volumeConfig, "azure_share"); exists {
			volumeTypeCount++
			azureShareExists = true
		} else if _, exists := validateArmContainerGroupVolumeVolumeExists(volumeConfig, "git_repo"); exists {
			volumeTypeCount++
			gitRepoExists = true
		} else if _, exists := validateArmContainerGroupVolumeVolumeExists(volumeConfig, "secret"); exists {
			volumeTypeCount++
			secretExists = true
		} else {
			//default to using empty_dir
			volumeTypeCount++
			emptyDirExists = true

		}

		if volumeTypeCount != 1 {
			errors = append(errors,
				fmt.Errorf(
					"validation failed for volume: '%v'. Only one volume type may be defined. It has azure_share: %t, git_repo: %t, secret: %t, empty_dir: %t",
					volumeConfig["name"], azureShareExists, gitRepoExists, secretExists, emptyDirExists))
		}
	}

	return
}

func expandContainerVolumes(d *schema.ResourceData) (*[]containerinstance.Volume, error) {
	volumesRaw := d.Get("volume").([]interface{})

	if len(volumesRaw) == 0 {
		return nil, nil
	}

	errors := validateArmContainerGroupVolume(volumesRaw)
	if len(errors) > 0 {
		return nil, fmt.Errorf("invalid volume configuration: %+v", errors)
	}

	containerGroupVolumes := make([]containerinstance.Volume, 0)

	for _, volumeRaw := range volumesRaw {
		volumeConfig := volumeRaw.(map[string]interface{})

		name := utils.String(volumeConfig["name"].(string))

		if azureShare, exists := validateArmContainerGroupVolumeVolumeExists(volumeConfig, "azure_share"); exists {
			shareName := azureShare["share_name"].(string)
			storageAccountName := azureShare["storage_account_name"].(string)
			storageAccountKey := azureShare["storage_account_key"].(string)

			cv := containerinstance.Volume{
				Name: name,
				AzureFile: &containerinstance.AzureFileVolume{
					ShareName:          utils.String(shareName),
					StorageAccountName: utils.String(storageAccountName),
					StorageAccountKey:  utils.String(storageAccountKey),
				},
			}

			containerGroupVolumes = append(containerGroupVolumes, cv)
			continue
		}

		if gitRepo, exists := validateArmContainerGroupVolumeVolumeExists(volumeConfig, "git_repo"); exists {
			repository := gitRepo["repository"].(string)
			directory := gitRepo["directory"].(string)

			cv := containerinstance.Volume{
				Name: name,
				GitRepo: &containerinstance.GitRepoVolume{
					Repository: utils.String(repository),
					Directory:  utils.String(directory),
				},
			}

			containerGroupVolumes = append(containerGroupVolumes, cv)
			continue
		}

		if _, exists := validateArmContainerGroupVolumeVolumeExists(volumeConfig, "secret"); exists {
			secretsConverted := map[string]*string{}
			secrets := volumeConfig["secret"].([]interface{})

			for _, v := range secrets {

				secret := v.(map[string]interface{})
				secretName := secret["name"].(string)
				secretData := secret["data"].(string)
				secretsConverted[secretName] = &secretData
			}

			cv := containerinstance.Volume{
				Name:   name,
				Secret: secretsConverted,
			}
			containerGroupVolumes = append(containerGroupVolumes, cv)
			continue
		}

		//default to empty_dir
		containerGroupVolumes = append(containerGroupVolumes, containerinstance.Volume{
			Name: name,
			// Note: object needs to be initialized for call to succeed but doesn't need to contain anything.
			EmptyDir: map[string]string{},
		})
	}

	return &containerGroupVolumes, nil
}

func expandContainerVolumeMounts(input interface{}) *[]containerinstance.VolumeMount {
	volumeMountsRaw := input.([]interface{})

	if len(volumeMountsRaw) == 0 {
		return nil
	}

	volumeMounts := make([]containerinstance.VolumeMount, 0)

	for _, mountRaw := range volumeMountsRaw {
		volumeMountConfig := mountRaw.(map[string]interface{})
		name := volumeMountConfig["volume_name"].(string)
		mountPath := volumeMountConfig["mount_path"].(string)
		readOnly := volumeMountConfig["read_only"].(bool)

		vm := containerinstance.VolumeMount{
			Name:      utils.String(name),
			MountPath: utils.String(mountPath),
			ReadOnly:  utils.Bool(readOnly),
		}

		volumeMounts = append(volumeMounts, vm)
	}

	return &volumeMounts
}
