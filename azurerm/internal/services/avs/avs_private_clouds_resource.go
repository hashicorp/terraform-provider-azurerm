package avs

import (
    "fmt"
    "github.com/Azure/azure-sdk-for-go/services/avs/mgmt/2020-03-20/avs"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
    "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
    "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
    "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
    "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
    "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/avs/parse"
    "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
    azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
    "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
    "github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
    "log"
    "time"
)
func resourceArmAvsPrivateCloud() *schema.Resource {
    return &schema.Resource{
        Create: resourceArmAvsPrivateCloudCreate,
        Read:   resourceArmAvsPrivateCloudRead,
        Update: resourceArmAvsPrivateCloudUpdate,
        Delete: resourceArmAvsPrivateCloudDelete,

        Timeouts: &schema.ResourceTimeout{
            Create: schema.DefaultTimeout(30 * time.Minute),
            Read: schema.DefaultTimeout(5 * time.Minute),
            Update: schema.DefaultTimeout(30 * time.Minute),
            Delete: schema.DefaultTimeout(30 * time.Minute),
        },

        Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
            _, err := parse.AvsPrivateCloudID(id)
            return err
        }),


        Schema: map[string]*schema.Schema{
            "name": {
                Type: schema.TypeString,
                Required: true,
                ForceNew: true,
            },

            "resource_group_name": azure.SchemaResourceGroupName(),

            "location": azure.SchemaLocation(),

            "sku": {
                Type: schema.TypeList,
                Required: true,
                ForceNew: true,
                MaxItems: 1,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "name": {
                            Type: schema.TypeString,
                            Required: true,
                            ForceNew: true,
                        },
                    },
                },
            },

            "management_cluster": {
                Type: schema.TypeList,
                Required: true,
                MaxItems: 1,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "cluster_size": {
                            Type: schema.TypeInt,
                            Required: true,
                        },

                        "cluster_id": {
                            Type: schema.TypeInt,
                            Computed: true,
                        },

                        "hosts": {
                            Type: schema.TypeSet,
                            Computed: true,
                            Elem: &schema.Schema{
                                Type: schema.TypeString,
                            },
                        },
                    },
                },
            },

            "network_block": {
                Type: schema.TypeString,
                Required: true,
            },

            "circuit": {
                Type: schema.TypeList,
                Computed: true,
                MaxItems: 1,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "express_route_id": {
                            Type: schema.TypeString,
                            Computed: true,
                        },

                        "express_route_private_peering_id": {
                            Type: schema.TypeString,
                            Computed: true,
                        },

                        "primary_subnet": {
                            Type: schema.TypeString,
                            Computed: true,
                        },

                        "secondary_subnet": {
                            Type: schema.TypeString,
                            Computed: true,
                        },
                    },
                },
            },

            "identity_source": {
                Type: schema.TypeSet,
                Optional: true,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "name": {
                            Type: schema.TypeString,
                            Optional: true,
                        },

                        "alias": {
                            Type: schema.TypeString,
                            Optional: true,
                        },

                        "base_group_dn": {
                            Type: schema.TypeString,
                            Optional: true,
                        },

                        "base_user_dn": {
                            Type: schema.TypeString,
                            Optional: true,
                        },

                        "domain": {
                            Type: schema.TypeString,
                            Optional: true,
                        },

                        "password": {
                            Type: schema.TypeString,
                            Optional: true,
                        },

                        "primary_server": {
                            Type: schema.TypeString,
                            Optional: true,
                        },

                        "secondary_server": {
                            Type: schema.TypeString,
                            Optional: true,
                        },

                        "ssl": {
                            Type: schema.TypeBool,
                            Optional: true,
                        },

                        "username": {
                            Type: schema.TypeString,
                            Optional: true,
                        },
                    },
                },
            },

            "internet": {
                Type: schema.TypeBool,
                Optional: true,
            },

            "nsxt_password": {
                Type: schema.TypeString,
                Optional: true,
            },

            "vcenter_password": {
                Type: schema.TypeString,
                Optional: true,
            },

            "endpoint": {
                Type: schema.TypeList,
                Computed: true,
                MaxItems: 1,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "hcx_cloud_manager": {
                            Type: schema.TypeString,
                            Computed: true,
                        },

                        "nsxt_manager": {
                            Type: schema.TypeString,
                            Computed: true,
                        },

                        "vcsa": {
                            Type: schema.TypeString,
                            Computed: true,
                        },
                    },
                },
            },

            "management_network": {
                Type: schema.TypeString,
                Computed: true,
            },

            "nsxt_certificate_thumbprint": {
                Type: schema.TypeString,
                Computed: true,
            },

            "provisioning_network": {
                Type: schema.TypeString,
                Computed: true,
            },

            "type": {
                Type: schema.TypeString,
                Computed: true,
            },

            "vcenter_certificate_thumbprint": {
                Type: schema.TypeString,
                Computed: true,
            },

            "vmotion_network": {
                Type: schema.TypeString,
                Computed: true,
            },

            "tags": tags.Schema(),
        },
    }
}
func resourceArmAvsPrivateCloudCreate(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*clients.Client).Avs.PrivateCloudClient
    ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
    defer cancel()

    name := d.Get("name").(string)
    resourceGroup := d.Get("resource_group_name").(string)

    existing, err :=client.Get(ctx, resourceGroup, name)
    if err != nil {
        if !utils.ResponseWasNotFound(existing.Response) {
            return fmt.Errorf("checking for present of existing Avs PrivateCloud %q (Resource Group %q): %+v", name, resourceGroup, err)
        }
    }
    if existing.ID != nil && *existing.ID != "" {
        return tf.ImportAsExistsError("azurerm_avs_private_cloud", *existing.ID)
    }

    internet := avs.Disabled
    if d.Get("internet").(bool) {
        internet = avs.Enabled
    }

    privateCloud := avs.PrivateCloud{
        Location: utils.String(location.Normalize(d.Get("location").(string))),
        Sku: expandArmPrivateCloudSku(d.Get("sku").([]interface{})),
        PrivateCloudProperties: &avs.PrivateCloudProperties{
            ManagementCluster: expandArmPrivateCloudManagementCluster(d.Get("management_cluster").([]interface{})),
            NetworkBlock: utils.String(d.Get("network_block").(string)),
            IdentitySources: expandArmPrivateCloudIdentitySourceArray(d.Get("identity_source").(*schema.Set).List()),
            Internet: internet,
            NsxtPassword: utils.String(d.Get("nsxt_password").(string)),
            VcenterPassword: utils.String(d.Get("vcenter_password").(string)),
        },
        Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
    }
    future, err :=client.CreateOrUpdate(ctx, resourceGroup, name, privateCloud)
    if err != nil {
        return fmt.Errorf("creating Avs PrivateCloud %q (Resource Group %q): %+v", name, resourceGroup, err)
    }

    if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
        return fmt.Errorf("waiting on creating future for Avs PrivateCloud %q (Resource Group %q): %+v", name, resourceGroup, err)
    }

    resp, err :=client.Get(ctx, resourceGroup, name)
    if err != nil {
        return fmt.Errorf("retrieving Avs PrivateCloud %q (Resource Group %q): %+v", name, resourceGroup, err)
    }

    if resp.ID == nil || *resp.ID== "" {
        return fmt.Errorf("empty or nil ID returned for Avs PrivateCloud %q (Resource Group %q) ID", name, resourceGroup)
    }

    d.SetId(*resp.ID)
    return resourceArmAvsPrivateCloudRead(d, meta)
}

func resourceArmAvsPrivateCloudRead(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*clients.Client).Avs.PrivateCloudClient
    ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
    defer cancel()

    id, err := parse.AvsPrivateCloudID(d.Id())
    if err != nil {
        return err
    }

    resp, err :=client.Get(ctx, id.ResourceGroup, id.Name)
    if err != nil {
        if utils.ResponseWasNotFound(resp.Response) {
            log.Printf("[INFO] avs %q does not exist - removing from state", d.Id())
            d.SetId("")
            return nil
        }
        return fmt.Errorf("retrieving Avs PrivateCloud %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
    }
    d.Set("name", id.Name)
    d.Set("resource_group_name", id.ResourceGroup)
    d.Set("location", location.NormalizeNilable(resp.Location))
    if props := resp.PrivateCloudProperties; props != nil {
        if err := d.Set("management_cluster", flattenArmPrivateCloudManagementCluster(props.ManagementCluster)); err != nil {
            return fmt.Errorf("setting `management_cluster`: %+v", err)
        }
        d.Set("network_block", props.NetworkBlock)
        if err := d.Set("circuit", flattenArmPrivateCloudCircuit(props.Circuit)); err != nil {
            return fmt.Errorf("setting `circuit`: %+v", err)
        }
        if err := d.Set("identity_source", flattenArmPrivateCloudIdentitySourceArray(props.IdentitySources)); err != nil {
            return fmt.Errorf("setting `identity_source`: %+v", err)
        }
        d.Set("internet", props.Internet==avs.Enabled)
        d.Set("nsxt_password", props.NsxtPassword)
        d.Set("vcenter_password", props.VcenterPassword)
        if err := d.Set("endpoint", flattenArmPrivateCloudEndpoints(props.Endpoints)); err != nil {
            return fmt.Errorf("setting `endpoint`: %+v", err)
        }
        d.Set("management_network", props.ManagementNetwork)
        d.Set("nsxt_certificate_thumbprint", props.NsxtCertificateThumbprint)
        d.Set("provisioning_network", props.ProvisioningNetwork)
        d.Set("vcenter_certificate_thumbprint", props.VcenterCertificateThumbprint)
        d.Set("vmotion_network", props.VmotionNetwork)
    }
    if err := d.Set("sku", flattenArmPrivateCloudSku(resp.Sku)); err != nil {
        return fmt.Errorf("setting `sku`: %+v", err)
    }
    d.Set("type", resp.Type)
    return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmAvsPrivateCloudUpdate(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*clients.Client).Avs.PrivateCloudClient
    ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
    defer cancel()

    id, err := parse.AvsPrivateCloudID(d.Id())
    if err != nil {
        return err
    }

    privateCloudUpdate := avs.PrivateCloudUpdate{
        PrivateCloudUpdateProperties: &avs.PrivateCloudUpdateProperties{
        },
    }
    if d.HasChange("management_cluster") {
        privateCloudUpdate.PrivateCloudUpdateProperties.ManagementCluster = expandArmPrivateCloudManagementCluster(d.Get("management_cluster").([]interface{}))
    }
    if d.HasChange("internet") {
        internet := avs.Disabled
        if d.Get("internet").(bool) {
            internet = avs.Enabled
        }
        privateCloudUpdate.PrivateCloudUpdateProperties.Internet = internet
    }
    if d.HasChange("identity_source") {
        privateCloudUpdate.PrivateCloudUpdateProperties.IdentitySources = expandArmPrivateCloudIdentitySourceArray(d.Get("identity_source").(*schema.Set).List())
    }
    if d.HasChange("tags") {
        privateCloudUpdate.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
    }

    future, err :=client.Update(ctx, id.ResourceGroup, id.Name, privateCloudUpdate)
    if err != nil {
        return fmt.Errorf("updating Avs PrivateCloud %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
    }

    if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
        return fmt.Errorf("waiting on updating future for Avs PrivateCloud %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
    }
    return resourceArmAvsPrivateCloudRead(d, meta)
}

func resourceArmAvsPrivateCloudDelete(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*clients.Client).Avs.PrivateCloudClient
    ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
    defer cancel()

    id, err := parse.AvsPrivateCloudID(d.Id())
    if err != nil {
        return err
    }

    future, err :=client.Delete(ctx, id.ResourceGroup, id.Name)
    if err != nil {
        return fmt.Errorf("deleting Avs PrivateCloud %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
    }

    if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
        return fmt.Errorf("waiting on deleting future for Avs PrivateCloud %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
    }
    return nil
}

func expandArmPrivateCloudSku(input []interface{}) *avs.Sku {
    if len(input) == 0 {
        return nil
    }
    v := input[0].(map[string]interface{})
    return &avs.Sku{
        Name: utils.String(v["name"].(string)),
    }
}

func expandArmPrivateCloudManagementCluster(input []interface{}) *avs.ManagementCluster {
    if len(input) == 0 {
        return nil
    }
    v := input[0].(map[string]interface{})
    return &avs.ManagementCluster{
        ClusterSize: utils.Int32(int32(v["cluster_size"].(int))),
    }
}

func expandArmPrivateCloudIdentitySourceArray(input []interface{}) *[]avs.IdentitySource {
    results := make([]avs.IdentitySource, 0)
    for _, item := range input {
        v := item.(map[string]interface{})
        sSL := avs.SslEnumDisabled
        if v["ssl"].(bool) {
            sSL = avs.SslEnumEnabled
        }
        result := avs.IdentitySource{
            Name: utils.String(v["name"].(string)),
            Alias: utils.String(v["alias"].(string)),
            Domain: utils.String(v["domain"].(string)),
            BaseUserDN: utils.String(v["base_user_dn"].(string)),
            BaseGroupDN: utils.String(v["base_group_dn"].(string)),
            PrimaryServer: utils.String(v["primary_server"].(string)),
            SecondaryServer: utils.String(v["secondary_server"].(string)),
            Ssl: sSL,
            Username: utils.String(v["username"].(string)),
            Password: utils.String(v["password"].(string)),
        }
        results = append(results, result)
    }
    return &results
}

func flattenArmPrivateCloudManagementCluster(input *avs.ManagementCluster) []interface{} {
    if input == nil {
        return make([]interface{}, 0)
    }

    var clusterSize int32
    if input.ClusterSize != nil {
        clusterSize = *input.ClusterSize
    }
    var clusterId int32
    if input.ClusterID != nil {
        clusterId = *input.ClusterID
    }
    return []interface{}{
        map[string]interface{}{
            "cluster_size": clusterSize,
            "cluster_id": clusterId,
            "hosts": utils.FlattenStringSlice(input.Hosts),
        },
    }
}

func flattenArmPrivateCloudCircuit(input *avs.Circuit) []interface{} {
    if input == nil {
        return make([]interface{}, 0)
    }

    var expressRouteId string
    if input.ExpressRouteID != nil {
        expressRouteId = *input.ExpressRouteID
    }
    var expressRoutePrivatePeeringId string
    if input.ExpressRoutePrivatePeeringID != nil {
        expressRoutePrivatePeeringId = *input.ExpressRoutePrivatePeeringID
    }
    var primarySubnet string
    if input.PrimarySubnet != nil {
        primarySubnet = *input.PrimarySubnet
    }
    var secondarySubnet string
    if input.SecondarySubnet != nil {
        secondarySubnet = *input.SecondarySubnet
    }
    return []interface{}{
        map[string]interface{}{
            "express_route_id": expressRouteId,
            "express_route_private_peering_id": expressRoutePrivatePeeringId,
            "primary_subnet": primarySubnet,
            "secondary_subnet": secondarySubnet,
        },
    }
}

func flattenArmPrivateCloudIdentitySourceArray(input *[]avs.IdentitySource) []interface{} {
    results := make([]interface{}, 0)
    if input == nil {
        return results
    }

    for _, item := range *input {
        var name string
        if item.Name != nil {
            name = *item.Name
        }
        var alias string
        if item.Alias != nil {
            alias = *item.Alias
        }
        var baseGroupDn string
        if item.BaseGroupDN != nil {
            baseGroupDn = *item.BaseGroupDN
        }
        var baseUserDn string
        if item.BaseUserDN != nil {
            baseUserDn = *item.BaseUserDN
        }
        var domain string
        if item.Domain != nil {
            domain = *item.Domain
        }
        var password string
        if item.Password != nil {
            password = *item.Password
        }
        var primaryServer string
        if item.PrimaryServer != nil {
            primaryServer = *item.PrimaryServer
        }
        var secondaryServer string
        if item.SecondaryServer != nil {
            secondaryServer = *item.SecondaryServer
        }
        sSL := item.Ssl==avs.SslEnumEnabled
        var username string
        if item.Username != nil {
            username = *item.Username
        }
        v := map[string]interface{}{
            "name": name,
            "alias": alias,
            "base_group_dn": baseGroupDn,
            "base_user_dn": baseUserDn,
            "domain": domain,
            "password": password,
            "primary_server": primaryServer,
            "secondary_server": secondaryServer,
            "ssl": sSL,
            "username": username,
        }
        results = append(results, v)
    }
    return results
}

func flattenArmPrivateCloudEndpoints(input *avs.Endpoints) []interface{} {
    if input == nil {
        return make([]interface{}, 0)
    }

    var hcxCloudManager string
    if input.HcxCloudManager != nil {
        hcxCloudManager = *input.HcxCloudManager
    }
    var nsxtManager string
    if input.NsxtManager != nil {
        nsxtManager = *input.NsxtManager
    }
    var vcsa string
    if input.Vcsa != nil {
        vcsa = *input.Vcsa
    }
    return []interface{}{
        map[string]interface{}{
            "hcx_cloud_manager": hcxCloudManager,
            "nsxt_manager": nsxtManager,
            "vcsa": vcsa,
        },
    }
}

func flattenArmPrivateCloudSku(input *avs.Sku) []interface{} {
    if input == nil {
        return make([]interface{}, 0)
    }

    var name string
    if input.Name != nil {
        name = *input.Name
    }
    return []interface{}{
        map[string]interface{}{
            "name": name,
        },
    }
}
