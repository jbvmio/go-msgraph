package msgraph

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strings"
)

// Win32LobApp contains properties and inherited properties for Win32 apps.
// https://docs.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobapp?view=graph-rest-beta
type Win32LobApp struct {
	ODataType                       string                       `json:"@odata.type"`
	ID                              string                       `json:"id,omitempty"`
	DisplayName                     string                       `json:"displayName,omitempty"`
	Description                     string                       `json:"description,omitempty"`
	Publisher                       string                       `json:"publisher,omitempty"`
	LargeIcon                       *MimeContent                 `json:"largeIcon,omitempty"`
	CreatedDateTime                 DateTimeOffset               `json:"createdDateTime,omitempty"`
	LastModifiedDateTime            DateTimeOffset               `json:"lastModifiedDateTime,omitempty"`
	IsFeatured                      bool                         `json:"isFeatured,omitempty"`
	PrivacyInformationUrl           string                       `json:"privacyInformationUrl,omitempty"`
	InformationUrl                  string                       `json:"informationUrl,omitempty"`
	Wwner                           string                       `json:"owner,omitempty"`
	Developer                       string                       `json:"developer,omitempty"`
	Notes                           string                       `json:"notes,omitempty"`
	UploadState                     int                          `json:"uploadState,omitempty"`
	PublishingState                 string                       `json:"publishingState,omitempty"`
	IsAssigned                      bool                         `json:"isAssigned,omitempty"`
	RoleScopeTagIds                 []string                     `json:"roleScopeTagIds,omitempty"`
	DependentAppCount               int                          `json:"dependentAppCount,omitempty"`
	SupersedingAppCount             int                          `json:"supersedingAppCount,omitempty"`
	SupersededAppCount              int                          `json:"supersededAppCount,omitempty"`
	CommittedContentVersion         string                       `json:"committedContentVersion,omitempty"`
	FileName                        string                       `json:"fileName,omitempty"`
	Size                            int64                        `json:"size,omitempty"`
	InstallCommandLine              string                       `json:"installCommandLine,omitempty"`
	UninstallCommandLine            string                       `json:"uninstallCommandLine,omitempty"`
	ApplicableArchitectures         string                       `json:"applicableArchitectures,omitempty"`
	MinimumSupportedOperatingSystem map[string]bool              `json:"minimumSupportedOperatingSystem,omitempty"`
	MinimumFreeDiskSpaceInMB        int                          `json:"minimumFreeDiskSpaceInMB,omitempty"`
	MinimumMemoryInMB               int                          `json:"minimumMemoryInMB,omitempty"`
	MinimumNumberOfProcessors       int                          `json:"minimumNumberOfProcessors,omitempty"`
	MinimumCpuSpeedInMHz            int                          `json:"minimumCpuSpeedInMHz,omitempty"`
	DetectionRules                  []Win32LobAppDetection       `json:"detectionRules,omitempty"`
	RequirementRules                []Win32LobAppRequirement     `json:"requirementRules,omitempty"`
	Rules                           []Win32LobAppRule            `json:"rules,omitempty"`
	InstallExperience               Win32LobAppInstallExperience `json:"installExperience,omitempty"`
	ReturnCodes                     []Win32LobAppReturnCode      `json:"returnCodes,omitempty"`
	MsiInformation                  *Win32LobAppMsiInformation   `json:"msiInformation,omitempty"`
	SetupFilePath                   string                       `json:"setupFilePath,omitempty"`
	MinimumSupportedWindowsRelease  string                       `json:"minimumSupportedWindowsRelease,omitempty"`
	DisplayVersion                  string                       `json:"displayVersion,omitempty"`
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

// NewWin32LobApp creates and returns a new Win32LobApp with defaults based on the given intunewin file.
func NewWin32LobApp(intuneWinFile string) (Win32LobApp, error) {
	meta, err := GetIntuneWin32AppMetadata(intuneWinFile, false)
	if err != nil {
		return Win32LobApp{}, err
	}

	win32LobApp := Win32LobApp{
		ODataType:     `#microsoft.graph.win32LobApp`,
		Developer:     `GoMSGraph`,
		Publisher:     `GoMSGraph`,
		Description:   meta.Name,
		DisplayName:   meta.Name,
		FileName:      meta.FileName,
		SetupFilePath: meta.SetupFile,
		InstallExperience: Win32LobAppInstallExperience{
			RunAsAccount:          `system`,
			DeviceRestartBehavior: `suppress`,
		},
		ApplicableArchitectures:        `x64`,
		MinimumSupportedWindowsRelease: `1607`,
		MinimumSupportedOperatingSystem: map[string]bool{
			"v10_1607": true,
		},
		ReturnCodes: DefaultWin32LobAppReturnCodes,
	}
	if meta.HasMsiInfo() {
		win32LobApp.MsiInformation = &Win32LobAppMsiInformation{}
		win32LobApp.Publisher = meta.MsiInfo.MsiPublisher
		win32LobApp.InstallCommandLine = fmt.Sprintf(`msiexec /%s "%s" /q`, "i", meta.SetupFile)
		win32LobApp.UninstallCommandLine = fmt.Sprintf(`msiexec /%s "%s" /q`, "x", meta.MsiInfo.MsiProductCode)
		win32LobApp.MsiInformation.ProductName = meta.Name
		win32LobApp.MsiInformation.Publisher = meta.MsiInfo.MsiPublisher
		win32LobApp.MsiInformation.ProductCode = meta.MsiInfo.MsiProductCode
		win32LobApp.MsiInformation.UpgradeCode = meta.MsiInfo.MsiUpgradeCode
		win32LobApp.MsiInformation.ProductVersion = meta.MsiInfo.MsiProductVersion
		win32LobApp.MsiInformation.RequiresReboot = meta.MsiInfo.MsiRequiresReboot
		win32LobApp.MsiInformation.PackageType = `dualPurpose`
		switch meta.MsiInfo.MsiExecutionContext {
		case `System`:
			win32LobApp.MsiInformation.PackageType = `perMachine`
		case `User`:
			win32LobApp.MsiInformation.PackageType = `perUser`
		}
		win32LobApp.DetectionRules = append(win32LobApp.DetectionRules, NewWin32LobAppProductCodeDetection(meta.MsiInfo.MsiProductCode))
	}
	return win32LobApp, nil
}

func GetIntuneWin32AppMetadata(intuneWinFile string, includeData bool) (*DetectionXML, error) {
	var detectionXML DetectionXML
	r, err := zip.OpenReader(intuneWinFile)
	if err != nil {
		return &detectionXML, err
	}
	defer r.Close()
	for _, f := range r.File {
		if strings.HasSuffix(f.Name, `etection.xml`) {
			rc, err := f.Open()
			if err != nil {
				return &detectionXML, fmt.Errorf("unable to open compressed XML: %w", err)
			}
			defer rc.Close()
			data, err := ioutil.ReadAll(rc)
			if err != nil {
				return &detectionXML, fmt.Errorf("unable to read compressed XML: %w", err)
			}
			err = xml.Unmarshal(data, &detectionXML)
			if err != nil {
				return &detectionXML, fmt.Errorf("unable to unmarshal XML: %w", err)
			}
			break
		}
	}
	for _, f := range r.File {
		if strings.HasSuffix(f.Name, detectionXML.FileName) {
			detectionXML.EncryptedContentSize = int64(f.CompressedSize64)
			if includeData {
				rc, err := f.Open()
				if err != nil {
					return &detectionXML, fmt.Errorf("unable to open compressed file: %w", err)
				}
				defer rc.Close()
				data, err := ioutil.ReadAll(rc)
				if err != nil {
					return &detectionXML, fmt.Errorf("unable to read compressed file: %w", err)
				}
				detectionXML.Data = data
			}
			break
		}
	}
	return &detectionXML, nil
}

type DetectionXML struct {
	XMLName                xml.Name `xml:"ApplicationInfo"`
	Name                   string
	EncryptedContentSize   int64
	UnencryptedContentSize int64
	FileName               string
	SetupFile              string
	EncryptionInfo         EncryptionInfo
	MsiInfo                MsiInfo
	Data                   []byte
}

func (d *DetectionXML) HasMsiInfo() bool {
	return d.MsiInfo != (MsiInfo{})
}

type EncryptionInfo struct {
	EncryptionKey        string `json:"encryptionKey"`
	MacKey               string `json:"macKey"`
	InitializationVector string `json:"initializationVector"`
	Mac                  string `json:"mac"`
	ProfileIdentifier    string `json:"profileIdentifier"`
	FileDigest           string `json:"fileDigest"`
	FileDigestAlgorithm  string `json:"fileDigestAlgorithm"`
}

type MsiInfo struct {
	MsiProductCode                string
	MsiProductVersion             string
	MsiPackageCode                string
	MsiUpgradeCode                string
	MsiExecutionContext           string
	MsiRequiresLogon              bool
	MsiRequiresReboot             bool
	MsiIsMachineInstall           bool
	MsiIsUserInstall              bool
	MsiIncludesServices           bool
	MsiIncludesODBCDataSource     bool
	MsiContainsSystemRegistryKeys bool
	MsiContainsSystemFolders      bool
	MsiPublisher                  string
}
