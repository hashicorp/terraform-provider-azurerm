package azurerm

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/arm/containerinstance"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmContainerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmContainerGroupCreate,
		Read:   resourceArmContainerGroupRead,
		Update: resourceArmContainerGroupCreate,
		Delete: resourceArmContainerGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"ip_address_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"os_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"tags": tagsSchema(),

			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
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
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},

						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func resourceArmContainerGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	containerGroupsClient := client.containerGroupsClient

	// container group properties
	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	location := d.Get("location").(string)
	OSType := d.Get("os_type").(string)
	if OSType == "" {
		OSType = "linux"
	}
	IPAddressType := d.Get("ip_address_type").(string)
	if IPAddressType == "" {
		IPAddressType = "public"
	}
	tags := d.Get("tags").(map[string]interface{})

	containersConfig := d.Get("container").([]interface{})
	containers := make([]containerinstance.Container, len(containersConfig))
	containerGroupPorts := make([]containerinstance.Port, len(containersConfig))

	for index, containerConfig := range containersConfig {
		data := containerConfig.(map[string]interface{})

		// required
		name := data["name"].(string)
		image := data["image"].(string)

		// optional
		cpu := data["cpu"].(float64)
		memory := data["memory"].(float64)
		port := int32(data["port"].(int))
		protocol := data["protocol"].(string)

		containerPort := containerinstance.ContainerPort{
			Port: &port,
		}

		// container group port (port number + protocol)
		containerGroupPort := containerinstance.Port{
			Port: &port,
		}

		if strings.ToUpper(protocol) == "TCP" || strings.ToUpper(protocol) == "UDP" {
			containerGroupPort.Protocol = containerinstance.ContainerGroupNetworkProtocol(strings.ToUpper(protocol))
		}
		containerGroupPorts[index] = containerGroupPort

		container := containerinstance.Container{
			Name: &name,
			ContainerProperties: &containerinstance.ContainerProperties{
				Image: &image,
				Ports: &[]containerinstance.ContainerPort{containerPort},
				Resources: &containerinstance.ResourceRequirements{
					Requests: &containerinstance.ResourceRequests{
						MemoryInGB: &memory,
						CPU:        &cpu,
					},
				},
			},
		}

		containers[index] = container
	}

	// type ContainerGroupProperties struct {
	// 	ProvisioningState        *string                    `json:"provisioningState,omitempty"`
	// 	Containers               *[]Container               `json:"containers,omitempty"`
	// 	ImageRegistryCredentials *[]ImageRegistryCredential `json:"imageRegistryCredentials,omitempty"`
	// 	RestartPolicy            ContainerRestartPolicy     `json:"restartPolicy,omitempty"`
	// 	IPAddress                *IPAddress                 `json:"ipAddress,omitempty"`
	// 	OsType                   OperatingSystemTypes       `json:"osType,omitempty"`
	// 	State                    *string                    `json:"state,omitempty"`
	// 	Volumes                  *[]Volume                  `json:"volumes,omitempty"`
	// }

	// type ContainerProperties struct {
	// 	Image                *string                          `json:"image,omitempty"`
	// 	Command              *[]string                        `json:"command,omitempty"`
	// 	Ports                *[]ContainerPort                 `json:"ports,omitempty"`
	// 	EnvironmentVariables *[]EnvironmentVariable           `json:"environmentVariables,omitempty"`
	// 	InstanceView         *ContainerPropertiesInstanceView `json:"instanceView,omitempty"`
	// 	Resources            *ResourceRequirements            `json:"resources,omitempty"`
	// 	VolumeMounts         *[]VolumeMount                   `json:"volumeMounts,omitempty"`
	// }

	containerGroup := containerinstance.ContainerGroup{
		Name:     &name,
		Location: &location,
		Tags:     expandTags(tags),
		ContainerGroupProperties: &containerinstance.ContainerGroupProperties{
			Containers: &containers,
			IPAddress: &containerinstance.IPAddress{
				Type:  &IPAddressType,
				Ports: &containerGroupPorts,
			},
			OsType: containerinstance.OperatingSystemTypes(OSType),
		},
	}

	_, error := containerGroupsClient.CreateOrUpdate(resGroup, name, containerGroup)
	if error != nil {
		return error
	}

	read, readErr := containerGroupsClient.Get(resGroup, name)
	if readErr != nil {
		return readErr
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read container group %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return nil
}
func resourceArmContainerGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	containterGroupsClient := client.containerGroupsClient

	id, err := parseAzureResourceID(d.Id())

	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["containerGroups"]

	resp, error := containterGroupsClient.Get(resGroup, name)

	if error != nil {
		return error
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("location", azureRMNormalizeLocation(*resp.Location))
	flattenAndSetTags(d, resp.Tags)

	d.Set("os_type", string(resp.OsType))
	d.Set("ip_address_type", *resp.IPAddress.Type)
	d.Set("ip_address", *resp.IPAddress.IP)

	// ports := *resp.IPAddress.Ports
	// containers := *resp.Containers

	return nil
}
func resourceArmContainerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	containterGroupsClient := client.containerGroupsClient

	// container group properties
	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	_, error := containterGroupsClient.Delete(resGroup, name)

	if error != nil {
		return error
	}

	return nil
}
