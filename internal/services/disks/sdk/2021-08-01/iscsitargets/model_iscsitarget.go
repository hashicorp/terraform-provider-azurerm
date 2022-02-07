package iscsitargets

type IscsiTarget struct {
	Id                *string               `json:"id,omitempty"`
	ManagedBy         *string               `json:"managedBy,omitempty"`
	ManagedByExtended *[]string             `json:"managedByExtended,omitempty"`
	Name              *string               `json:"name,omitempty"`
	Properties        IscsiTargetProperties `json:"properties"`
	SystemData        *SystemMetadata       `json:"systemData,omitempty"`
	Type              *string               `json:"type,omitempty"`
}
