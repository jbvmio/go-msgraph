package msgraph

import "fmt"

// Win32LobApps is a collection of Win32LobApp.
type Win32LobApps []Win32LobApp

func (A Win32LobApps) setGraphClient(gC *GraphClient) {
	for i := 0; i < len(A); i++ {
		A[i].setGraphClient(gC)
	}
}

// Win32LobApp contains properties and inherited properties for Win32 apps.
// https://docs.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobapp?view=graph-rest-beta
type Win32LobApp struct {
	ODataType                       string                       `json:"@odata.type" yaml:"@odata.type"`
	ID                              string                       `json:"id,omitempty" yaml:"id,omitempty"`
	DisplayName                     string                       `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	Description                     string                       `json:"description,omitempty" yaml:"description,omitempty"`
	Publisher                       string                       `json:"publisher,omitempty" yaml:"publisher,omitempty"`
	LargeIcon                       *MimeContent                 `json:"largeIcon,omitempty" yaml:"largeIcon,omitempty"`
	CreatedDateTime                 DateTimeOffset               `json:"createdDateTime,omitempty" yaml:"createdDateTime,omitempty"`
	LastModifiedDateTime            DateTimeOffset               `json:"lastModifiedDateTime,omitempty" yaml:"lastModifiedDateTime,omitempty"`
	IsFeatured                      bool                         `json:"isFeatured,omitempty" yaml:"isFeatured,omitempty"`
	PrivacyInformationUrl           string                       `json:"privacyInformationUrl,omitempty" yaml:"privacyInformationUrl,omitempty"`
	InformationUrl                  string                       `json:"informationUrl,omitempty" yaml:"informationUrl,omitempty"`
	Owner                           string                       `json:"owner,omitempty" yaml:"owner,omitempty"`
	Developer                       string                       `json:"developer,omitempty" yaml:"developer,omitempty"`
	Notes                           string                       `json:"notes,omitempty" yaml:"notes,omitempty"`
	UploadState                     int                          `json:"uploadState,omitempty" yaml:"uploadState,omitempty"`
	PublishingState                 string                       `json:"publishingState,omitempty" yaml:"publishingState,omitempty"`
	IsAssigned                      bool                         `json:"isAssigned,omitempty" yaml:"isAssigned,omitempty"`
	RoleScopeTagIds                 []string                     `json:"roleScopeTagIds,omitempty" yaml:"roleScopeTagIds,omitempty"`
	DependentAppCount               int                          `json:"dependentAppCount,omitempty" yaml:"dependentAppCount,omitempty"`
	SupersedingAppCount             int                          `json:"supersedingAppCount,omitempty" yaml:"supersedingAppCount,omitempty"`
	SupersededAppCount              int                          `json:"supersededAppCount,omitempty" yaml:"supersededAppCount,omitempty"`
	CommittedContentVersion         string                       `json:"committedContentVersion,omitempty" yaml:"committedContentVersion,omitempty"`
	FileName                        string                       `json:"fileName,omitempty" yaml:"fileName,omitempty"`
	Size                            int64                        `json:"size,omitempty" yaml:"size,omitempty"`
	InstallCommandLine              string                       `json:"installCommandLine,omitempty" yaml:"installCommandLine,omitempty"`
	UninstallCommandLine            string                       `json:"uninstallCommandLine,omitempty" yaml:"uninstallCommandLine,omitempty"`
	ApplicableArchitectures         string                       `json:"applicableArchitectures,omitempty" yaml:"applicableArchitectures,omitempty"`
	MinimumSupportedOperatingSystem map[string]bool              `json:"minimumSupportedOperatingSystem,omitempty" yaml:"minimumSupportedOperatingSystem,omitempty"`
	MinimumFreeDiskSpaceInMB        int                          `json:"minimumFreeDiskSpaceInMB,omitempty" yaml:"minimumFreeDiskSpaceInMB,omitempty"`
	MinimumMemoryInMB               int                          `json:"minimumMemoryInMB,omitempty" yaml:"minimumMemoryInMB,omitempty"`
	MinimumNumberOfProcessors       int                          `json:"minimumNumberOfProcessors,omitempty" yaml:"minimumNumberOfProcessors,omitempty"`
	MinimumCpuSpeedInMHz            int                          `json:"minimumCpuSpeedInMHz,omitempty" yaml:"minimumCpuSpeedInMHz,omitempty"`
	DetectionRules                  []Win32LobAppDetection       `json:"detectionRules,omitempty" yaml:"detectionRules,omitempty"`
	RequirementRules                []Win32LobAppRequirement     `json:"requirementRules,omitempty" yaml:"requirementRules,omitempty"`
	Rules                           []Win32LobAppRule            `json:"rules,omitempty" yaml:"rules,omitempty"`
	InstallExperience               Win32LobAppInstallExperience `json:"installExperience,omitempty" yaml:"installExperience,omitempty"`
	ReturnCodes                     []Win32LobAppReturnCode      `json:"returnCodes,omitempty" yaml:"returnCodes,omitempty"`
	MsiInformation                  *Win32LobAppMsiInformation   `json:"msiInformation,omitempty" yaml:"msiInformation,omitempty"`
	SetupFilePath                   string                       `json:"setupFilePath,omitempty" yaml:"setupFilePath,omitempty"`
	MinimumSupportedWindowsRelease  string                       `json:"minimumSupportedWindowsRelease,omitempty" yaml:"minimumSupportedWindowsRelease,omitempty"`
	DisplayVersion                  string                       `json:"displayVersion,omitempty" yaml:"displayVersion,omitempty"`

	contentVersion TypeAndID
	graphClient    *GraphClient // the graphClient that called the Win32LobApp
}

func (A *Win32LobApp) CreateContentVersion() (TypeAndID, error) {
	tid, err := A.graphClient.CreateWin32LobAppContentVersion(A.ID)
	A.contentVersion = tid
	return tid, err
}

func (A *Win32LobApp) CreateContentFile(req MobileAppContentFileRequest) (MobileAppContentFile, error) {
	if A.contentVersion.ID == "" {
		tid, err := A.graphClient.CreateWin32LobAppContentVersion(A.ID)
		if err != nil {
			return MobileAppContentFile{}, fmt.Errorf("error generating new content version: %w", err)
		}
		A.contentVersion = tid
	}
	return A.graphClient.CreateWin32LobAppContentFile(A.ID, A.contentVersion.ID, req)
}

func (A *Win32LobApp) CreateContentFileWithVersion(version string, req MobileAppContentFileRequest) (MobileAppContentFile, error) {
	return A.graphClient.CreateWin32LobAppContentFile(A.ID, version, req)
}

func (A *Win32LobApp) setGraphClient(gC *GraphClient) {
	A.graphClient = gC
}

// MimeContent contains properties for a generic mime content.
// https://docs.microsoft.com/en-us/graph/api/resources/intune-shared-mimecontent?view=graph-rest-beta
type MimeContent struct {
	Type  string `json:"type"`
	Value []byte `json:"value"`
}

type Win32LobAppMsiInformation struct {
	ProductCode    string `json:"productCode"`
	ProductVersion string `json:"productVersion"`
	UpgradeCode    string `json:"upgradeCode"`
	RequiresReboot bool   `json:"requiresReboot"`
	PackageType    string `json:"packageType"`
	ProductName    string `json:"productName"`
	Publisher      string `json:"publisher"`
}

type Win32LobAppInstallExperience struct {
	RunAsAccount          string `json:"runAsAccount"`
	DeviceRestartBehavior string `json:"deviceRestartBehavior"`
}

type Win32LobAppReturnCode struct {
	ReturnCode int    `json:"returnCode"`
	Type       string `json:"type"`
}

type Win32LobAppRequirement interface {
}

type Win32LobAppRule interface {
}

var DefaultWin32LobAppReturnCodes = []Win32LobAppReturnCode{
	{
		ReturnCode: 0,
		Type:       `success`,
	},
	{
		ReturnCode: 1707,
		Type:       `success`,
	},
	{
		ReturnCode: 3010,
		Type:       `softReboot`,
	},
	{
		ReturnCode: 1641,
		Type:       `hardReboot`,
	},
	{
		ReturnCode: 1618,
		Type:       `retry`,
	},
}
