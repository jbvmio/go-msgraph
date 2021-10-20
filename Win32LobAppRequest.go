package msgraph

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

// Win32LobAppRequest is used to contruct a request to submit a new Win32LobApp.
type Win32LobAppRequest Win32LobApp

// NewWin32LobAppRequest creates and returns a new NewWin32LobAppRequest with defaults based on the given intunewin file.
func NewWin32LobAppRequest(xmlMeta *DetectionXML) Win32LobAppRequest {
	win32LobApp := Win32LobAppRequest{
		ODataType:     `#microsoft.graph.win32LobApp`,
		Developer:     `GoMSGraph`,
		Publisher:     `GoMSGraph`,
		Description:   xmlMeta.Name,
		DisplayName:   xmlMeta.Name,
		FileName:      xmlMeta.FileName,
		SetupFilePath: xmlMeta.SetupFile,
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
	if xmlMeta.HasMsiInfo() {
		win32LobApp.MsiInformation = &Win32LobAppMsiInformation{}
		win32LobApp.Publisher = xmlMeta.MsiInfo.MsiPublisher
		win32LobApp.DisplayVersion = xmlMeta.MsiInfo.MsiProductVersion
		win32LobApp.InstallCommandLine = fmt.Sprintf(`msiexec /%s "%s" /q`, "i", xmlMeta.SetupFile)
		win32LobApp.UninstallCommandLine = fmt.Sprintf(`msiexec /%s "%s" /q`, "x", xmlMeta.MsiInfo.MsiProductCode)
		win32LobApp.MsiInformation.ProductName = xmlMeta.Name
		win32LobApp.MsiInformation.Publisher = xmlMeta.MsiInfo.MsiPublisher
		win32LobApp.MsiInformation.ProductCode = xmlMeta.MsiInfo.MsiProductCode
		win32LobApp.MsiInformation.UpgradeCode = xmlMeta.MsiInfo.MsiUpgradeCode
		win32LobApp.MsiInformation.ProductVersion = xmlMeta.MsiInfo.MsiProductVersion
		win32LobApp.MsiInformation.RequiresReboot = xmlMeta.MsiInfo.MsiRequiresReboot
		win32LobApp.MsiInformation.PackageType = `dualPurpose`
		switch xmlMeta.MsiInfo.MsiExecutionContext {
		case `System`:
			win32LobApp.MsiInformation.PackageType = `perMachine`
		case `User`:
			win32LobApp.MsiInformation.PackageType = `perUser`
		}
		win32LobApp.DetectionRules = append(win32LobApp.DetectionRules, NewWin32LobAppProductCodeDetection(xmlMeta.MsiInfo.MsiProductCode))
	}
	return win32LobApp
}

func GetIntuneWin32AppMetadataFromFile(intuneWinFile string, includeData bool) (*DetectionXML, error) {
	file, err := ioutil.ReadFile(intuneWinFile)
	if err != nil {
		return &DetectionXML{}, err
	}
	return GetIntuneWin32AppMetadata(bytes.NewReader(file), includeData)
}

func GetIntuneWin32AppMetadata(intuneWinFile io.Reader, includeData bool) (*DetectionXML, error) {
	var detectionXML DetectionXML
	//r, err := zip.OpenReader(intuneWinFile)
	data, err := ioutil.ReadAll(intuneWinFile)
	if err != nil {
		return &detectionXML, err
	}
	l := len(data)
	r, err := zip.NewReader(bytes.NewReader(data), int64(l))
	if err != nil {
		return &detectionXML, err
	}
	//defer r.Close()
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
