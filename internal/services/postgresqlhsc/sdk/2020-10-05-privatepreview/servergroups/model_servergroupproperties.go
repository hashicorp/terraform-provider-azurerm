package servergroups

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

type ServerGroupProperties struct {
	AdministratorLogin         *string                                        `json:"administratorLogin,omitempty"`
	AdministratorLoginPassword *string                                        `json:"administratorLoginPassword,omitempty"`
	AvailabilityZone           *string                                        `json:"availabilityZone,omitempty"`
	BackupRetentionDays        *int64                                         `json:"backupRetentionDays,omitempty"`
	CitusVersion               *CitusVersion                                  `json:"citusVersion,omitempty"`
	CreateMode                 *CreateMode                                    `json:"createMode,omitempty"`
	DelegatedSubnetArguments   *ServerGroupPropertiesDelegatedSubnetArguments `json:"delegatedSubnetArguments,omitempty"`
	EarliestRestoreTime        *string                                        `json:"earliestRestoreTime,omitempty"`
	EnableMx                   *bool                                          `json:"enableMx,omitempty"`
	EnableShardsOnCoordinator  *bool                                          `json:"enableShardsOnCoordinator,omitempty"`
	EnableZfs                  *bool                                          `json:"enableZfs,omitempty"`
	MaintenanceWindow          *MaintenanceWindow                             `json:"maintenanceWindow,omitempty"`
	PointInTimeUTC             *string                                        `json:"pointInTimeUTC,omitempty"`
	PostgresqlVersion          *PostgreSQLVersion                             `json:"postgresqlVersion,omitempty"`
	PrivateDnsZoneArguments    *ServerGroupPropertiesPrivateDnsZoneArguments  `json:"privateDnsZoneArguments,omitempty"`
	ReadReplicas               *[]string                                      `json:"readReplicas,omitempty"`
	ResourceProviderType       *ResourceProviderType                          `json:"resourceProviderType,omitempty"`
	ServerRoleGroups           *[]ServerRoleGroup                             `json:"serverRoleGroups,omitempty"`
	SourceLocation             *string                                        `json:"sourceLocation,omitempty"`
	SourceResourceGroupName    *string                                        `json:"sourceResourceGroupName,omitempty"`
	SourceServerGroup          *string                                        `json:"sourceServerGroup,omitempty"`
	SourceServerGroupName      *string                                        `json:"sourceServerGroupName,omitempty"`
	SourceSubscriptionId       *string                                        `json:"sourceSubscriptionId,omitempty"`
	StandbyAvailabilityZone    *string                                        `json:"standbyAvailabilityZone,omitempty"`
	State                      *ServerState                                   `json:"state,omitempty"`
}

func (o ServerGroupProperties) GetEarliestRestoreTimeAsTime() (*time.Time, error) {
	if o.EarliestRestoreTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EarliestRestoreTime, "2006-01-02T15:04:05Z07:00")
}

func (o ServerGroupProperties) SetEarliestRestoreTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EarliestRestoreTime = &formatted
}

func (o ServerGroupProperties) GetPointInTimeUTCAsTime() (*time.Time, error) {
	if o.PointInTimeUTC == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.PointInTimeUTC, "2006-01-02T15:04:05Z07:00")
}

func (o ServerGroupProperties) SetPointInTimeUTCAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.PointInTimeUTC = &formatted
}
