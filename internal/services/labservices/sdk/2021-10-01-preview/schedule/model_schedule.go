package schedule

type Schedule struct {
	Id         *string            `json:"id,omitempty"`
	Name       *string            `json:"name,omitempty"`
	Properties ScheduleProperties `json:"properties"`
	SystemData *SystemData        `json:"systemData,omitempty"`
	Type       *string            `json:"type,omitempty"`
}
