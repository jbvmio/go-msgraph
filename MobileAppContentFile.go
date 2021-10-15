package msgraph

// MobileAppContentFile contains properties for a single installer file that is associated with a given mobileAppContent version.
// https://docs.microsoft.com/en-us/graph/api/resources/intune-apps-mobileappcontentfile?view=graph-rest-beta
type MobileAppContentFile struct {
	ID                                string         `json:"id,omitempty"`
	AzureStorageUri                   string         `json:"azureStorageUri,omitempty"`
	IsCommitted                       bool           `json:"isCommitted,omitempty"`
	CreatedDateTime                   DateTimeOffset `json:"createdDateTime,omitempty"`
	Name                              string         `json:"name,omitempty"`
	Size                              int64          `json:"size,omitempty"`
	SizeEncrypted                     int64          `json:"sizeEncrypted,omitempty"`
	AzureStorageUriExpirationDateTime DateTimeOffset `json:"azureStorageUriExpirationDateTime,omitempty"`
	Manifest                          []byte         `json:"manifest,omitempty"`
	UploadState                       string         `json:"uploadState,omitempty"`
	IsFrameworkFile                   bool           `json:"isFrameworkFile,omitempty"`
	IsDependency                      bool           `json:"isDependency,omitempty"`
}
