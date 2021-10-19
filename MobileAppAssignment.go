package msgraph

type MobileAppAssignment struct {
	ID       string                                 `json:"id"`
	Intent   string                                 `json:"intent,omitempty"`
	Target   DeviceAndAppManagementAssignmentTarget `json:"target,omitempty"`
	Settings interface{}                            `json:"settings,omitempty"`
	Source   string                                 `json:"source,omitempty"`
	SourceID string                                 `json:"sourceId,omitempty"`
}

type DeviceAndAppManagementAssignmentTarget struct {
	FilterID   string `json:"deviceAndAppManagementAssignmentFilterId"`
	FilterType string `json:"deviceAndAppManagementAssignmentFilterType"`
}
