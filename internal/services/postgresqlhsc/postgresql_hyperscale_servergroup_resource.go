package postgresqlhsc

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgresqlhsc/sdk/2020-10-05-privatepreview/servergroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type PostgreSQLHyperScaleServerGroupResource struct {
}

var _ sdk.Resource = PostgreSQLHyperScaleServerGroupResource{}
var _ sdk.ResourceWithUpdate = PostgreSQLHyperScaleServerGroupResource{}

// var _ sdk.ResourceWithCustomImporter = PostgreSQLHyperScaleServerGroupResource{}
// var _ sdk.ResourceWithCustomizeDiff = PostgreSQLHyperScaleServerGroupResource{}

type PostgreSQLHyperScaleServerGroupResourceModel struct {
	Name                    string                 `tfschema:"name"`
	ResourceGroup           string                 `tfschema:"resource_group_name"`
	Location                string                 `tfschema:"location"`
	Tags                    map[string]string      `tfschema:"tags"`
	CreateMode              string                 `tfschema:"create_mode"`
	AdministratorPassword   string                 `tfschema:"administrator_login_password"`
	BackupRetentionDays     int                    `tfschema:"backup_retention_days"`
	CitusVersion            string                 `tfschema:"citus_version"`
	EnableMX                bool                   `tfschema:"mx_enabled"`
	EnableZFS               bool                   `tfschema:"zfs_enabled"`
	PreviewFeatures         bool                   `tfschema:"preview_features"`
	PostgreSQLVersion       string                 `tfschema:"postgresql_version"`
	LogMinDurationStatement int                    `tfschema:"log_min_duration_statement"`
	TrackActivityQuerySize  int                    `tfschema:"track_activity_query_size"`
	ServerRoleGroup         []ServerRoleGroupModel `tfschema:"server_role_group"`
}

type ServerRoleGroupModel struct {
	Name             string `tfschema:"name"`
	Role             string `tfschema:"role"`
	ServerCount      int    `tfschema:"server_count"`
	ServerEdition    string `tfschema:"server_edition"`
	VCores           int    `tfschema:"vcores"`
	StorageQuotaInMB int    `tfschema:"storage_quota_in_mb"`
	EnableHA         bool   `tfschema:"ha_enabled"`
}

func (r PostgreSQLHyperScaleServerGroupResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": azure.SchemaResourceGroupName(),

		"location": location.Schema(),

		"create_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(servergroups.CreateModeDefault),
			ValidateFunc: validation.StringInSlice([]string{
				string(servergroups.CreateModeDefault),
				string(servergroups.CreateModePointInTimeRestore),
				string(servergroups.CreateModeReadReplica),
			}, false),
		},

		"administrator_login_password": {
			Type:      pluginsdk.TypeString,
			Optional:  true,
			Sensitive: true,
		},

		"backup_retention_days": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.IntBetween(7, 35),
		},

		"citus_version": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(servergroups.CitusVersionEightPointThree),
				string(servergroups.CitusVersionNinePointZero),
				string(servergroups.CitusVersionNinePointOne),
				string(servergroups.CitusVersionNinePointTwo),
				string(servergroups.CitusVersionNinePointThree),
				string(servergroups.CitusVersionNinePointFour),
				string(servergroups.CitusVersionNinePointFive),
				"10.2",
			}, false),
		},

		"postgresql_version": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(servergroups.PostgreSQLVersionOneOne),
				string(servergroups.PostgreSQLVersionOneTwo),
				"13",
			}, false),
		},

		"preview_features": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
		},

		"log_min_duration_statement": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
		},

		"mx_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
		},

		"zfs_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
		},

		"track_activity_query_size": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
		},

		"server_role_group": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "",
					},
					"role": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(servergroups.ServerRoleCoordinator),
							string(servergroups.ServerRoleWorker),
						}, false),
					},
					"server_count": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},
					"server_edition": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  servergroups.ServerEditionGeneralPurpose,
						ValidateFunc: validation.StringInSlice([]string{
							string(servergroups.ServerEditionGeneralPurpose),
							string(servergroups.ServerEditionMemoryOptimized),
						}, false),
					},
					"vcores": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},
					"storage_quota_in_mb": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},
					"ha_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"tags": tags.Schema(),
	}
}

func (r PostgreSQLHyperScaleServerGroupResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r PostgreSQLHyperScaleServerGroupResource) ModelObject() interface{} {
	return &PostgreSQLHyperScaleServerGroupResourceModel{}
}

func (r PostgreSQLHyperScaleServerGroupResource) ResourceType() string {
	return "azurerm_postgresql_hyperscale_servergroup"
}

func (r PostgreSQLHyperScaleServerGroupResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return servergroups.ValidateServerGroupsv2ID
}

// func (r PostgreSQLHyperScaleServerGroupResource) CustomImporter() sdk.ResourceRunFunc {
// 	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
// 		client := metadata.Client.PostgreSQLHSC.ServerGroupsClient
//
// 		id, err := servergroups.ParseServerGroupsv2ID(metadata.ResourceData.Id())
// 		if err != nil {
// 			return fmt.Errorf("while parsing resource ID: %+v", err)
// 		}
// 		_, err = client.Get(ctx, *id)
// 		if err != nil {
// 			return fmt.Errorf("while checking for Server Group's %q existence: %+v", id.ServerGroupName, err)
// 		}
//
// 		metadata.ResourceData.Set("create_mode", "Default")
// 		return nil
// 	}
// }
//
// func (k PostgreSQLHyperScaleServerGroupResource) CustomizeDiff() sdk.ResourceFunc {
// 	return sdk.ResourceFunc{
// 		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
// 			rd := metadata.ResourceDiff
// 			rd.SetNew("create_mode", "Default")
// 			return nil
// 		},
// 		Timeout: 30 * time.Minute,
// 	}
// }

func (r PostgreSQLHyperScaleServerGroupResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model PostgreSQLHyperScaleServerGroupResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			client := metadata.Client.PostgreSQLHSC.ServerGroupsClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := servergroups.NewServerGroupsv2ID(subscriptionId, model.ResourceGroup, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Server Group %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			// @favoretti: This can only be 'citus'
			administratorLogin := "citus"
			citusVersion := servergroups.CitusVersion(model.CitusVersion)
			createMode := servergroups.CreateMode(model.CreateMode)
			postgresqlVersion := servergroups.PostgreSQLVersion(model.PostgreSQLVersion)

			hscServerGroup := servergroups.ServerGroup{
				Name:     &model.Name,
				Location: model.Location,
				Tags:     &model.Tags,
				Properties: &servergroups.ServerGroupProperties{
					AdministratorLogin:         &administratorLogin,
					AdministratorLoginPassword: &model.AdministratorPassword,
					BackupRetentionDays:        utils.Int64(int64(model.BackupRetentionDays)),
					CitusVersion:               &citusVersion,
					CreateMode:                 &createMode,
					EnableMx:                   &model.EnableMX,
					EnableZfs:                  &model.EnableZFS,
					PostgresqlVersion:          &postgresqlVersion,
					ServerRoleGroups:           expandServerRoleGroupModel(model.ServerRoleGroup),
				},
			}

			_, err = client.CreateOrUpdate(ctx, id, hscServerGroup)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r PostgreSQLHyperScaleServerGroupResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {

			client := metadata.Client.PostgreSQLHSC.ServerGroupsClient
			id, err := servergroups.ParseServerGroupsv2ID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state PostgreSQLHyperScaleServerGroupResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Server Group %s: %v", id, err)
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Model.Tags = &state.Tags
			}

			_, err = client.CreateOrUpdate(ctx, *id, *existing.Model)
			if err != nil {
				return fmt.Errorf("updating Server Group %s: %+v", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r PostgreSQLHyperScaleServerGroupResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := servergroups.ParseServerGroupsv2ID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			client := metadata.Client.PostgreSQLHSC.ServerGroupsClient

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("while checking for Server Group's %q existence: %+v", id.ServerGroupName, err)
			}

			state := PostgreSQLHyperScaleServerGroupResourceModel{
				Name:                  id.ServerGroupName,
				Location:              location.NormalizeNilable(utils.String(resp.Model.Location)),
				ResourceGroup:         id.ResourceGroupName,
				CreateMode:            metadata.ResourceData.Get("create_mode").(string),
				AdministratorPassword: metadata.ResourceData.Get("administrator_login_password").(string),
			}

			if model := resp.Model; model != nil {
				if model.Tags != nil {
					state.Tags = *model.Tags
				}
				if props := model.Properties; props != nil {
					if props.BackupRetentionDays != nil {
						state.BackupRetentionDays = int(*props.BackupRetentionDays)
					}
					if props.CitusVersion != nil {
						state.CitusVersion = string(*props.CitusVersion)
					}
					if props.EnableMx != nil {
						state.EnableMX = *props.EnableMx
					}
					if props.EnableZfs != nil {
						state.EnableZFS = *props.EnableZfs
					}
					if props.PostgresqlVersion != nil {
						state.PostgreSQLVersion = string(*props.PostgresqlVersion)
					}
					if props.ServerRoleGroups != nil {
						state.ServerRoleGroup = flattenServerRoleGroupModel(props.ServerRoleGroups)
					}
				}
			}
			return metadata.Encode(&state)
		},
		Timeout: 5 * time.Minute,
	}
}

func (r PostgreSQLHyperScaleServerGroupResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := servergroups.ParseServerGroupsv2ID(metadata.ResourceData.Id())

			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			client := metadata.Client.PostgreSQLHSC.ServerGroupsClient

			_, err = client.Delete(ctx, *id)
			if err != nil {
				return fmt.Errorf("while removing Server Group %q: %+v", id.ServerGroupName, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func expandServerRoleGroupModel(input []ServerRoleGroupModel) *[]servergroups.ServerRoleGroup {
	var serverRoleGroups []servergroups.ServerRoleGroup
	if input == nil {
		return &serverRoleGroups
	}

	for _, v := range input {
		serverRole := servergroups.ServerRole(v.Role)
		serverEdition := servergroups.ServerEdition(v.ServerEdition)
		serverRoleGroups = append(serverRoleGroups, servergroups.ServerRoleGroup{
			Name:             utils.String(v.Name),
			Role:             &serverRole,
			ServerCount:      utils.Int64(int64(v.ServerCount)),
			ServerEdition:    &serverEdition,
			VCores:           utils.Int64(int64(v.VCores)),
			StorageQuotaInMb: utils.Int64(int64(v.StorageQuotaInMB)),
			EnableHa:         &v.EnableHA,
		})
	}

	return &serverRoleGroups
}

func flattenServerRoleGroupModel(input *[]servergroups.ServerRoleGroup) []ServerRoleGroupModel {
	var output []ServerRoleGroupModel
	if input == nil || len(*input) == 0 {
		return output
	}

	for _, v := range *input {
		if v.Name == nil {
			continue
		}

		output = append(output, ServerRoleGroupModel{
			Name:             *v.Name,
			Role:             string(*v.Role),
			ServerCount:      int(*v.ServerCount),
			ServerEdition:    string(*v.ServerEdition),
			VCores:           int(*v.VCores),
			StorageQuotaInMB: int(*v.StorageQuotaInMb),
			EnableHA:         *v.EnableHa,
		})
	}
	return output
}
