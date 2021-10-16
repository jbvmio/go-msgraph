package msgraph

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// MobileAppContentFiles is a collection of MobileAppContentFile.
type MobileAppContentFiles []MobileAppContentFile

func (F MobileAppContentFiles) setGraphClient(gC *GraphClient) {
	for i := 0; i < len(F); i++ {
		F[i].setGraphClient(gC)
	}
}

// MobileAppContentFile contains properties for a single installer file that is associated with a given mobileAppContent version.
// https://docs.microsoft.com/en-us/graph/api/resources/intune-apps-mobileappcontentfile?view=graph-rest-beta
type MobileAppContentFile struct {
	ODataContext                      string          `json:"@odata.context,omitempty"`
	Context                           MACFContext     `json:"-"`
	ID                                string          `json:"id,omitempty"`
	AzureStorageUri                   string          `json:"azureStorageUri,omitempty"`
	IsCommitted                       bool            `json:"isCommitted,omitempty"`
	CreatedDateTime                   *DateTimeOffset `json:"createdDateTime,omitempty"`
	Name                              string          `json:"name,omitempty"`
	Size                              int64           `json:"size,omitempty"`
	SizeEncrypted                     int64           `json:"sizeEncrypted,omitempty"`
	AzureStorageUriExpirationDateTime *DateTimeOffset `json:"azureStorageUriExpirationDateTime,omitempty"`
	Manifest                          []byte          `json:"manifest,omitempty"`
	UploadState                       string          `json:"uploadState,omitempty"`
	IsFrameworkFile                   bool            `json:"isFrameworkFile,omitempty"`
	IsDependency                      bool            `json:"isDependency"`

	graphClient *GraphClient // the graphClient that called the MobileAppContentFile
}

// https://graph.microsoft.com/beta/$metadata#deviceAppManagement/mobileApps('fb313955-3b52-4edf-9b0d-0222987084b7')/microsoft.graph.win32LobApp/contentVersions('1')/files/$entity
// https://graph.microsoft.com/beta/$metadata#deviceAppManagement/mobileApps('fb313955-3b52-4edf-9b0d-0222987084b7')/microsoft.graph.win32LobApp/contentVersions('1')/files/$entity

type MACFContext struct {
	AppID          string
	ContentVersion string
}

func (F *MobileAppContentFile) UploadIntuneWin(intuneWinFile string) error {
	return F.graphClient.Win32LobAppContentFileUpload(intuneWinFile, F)
}

func (F *MobileAppContentFile) ContinueWhenStorageReady(timeout time.Duration) error {
	F.Refresh()
	if F.StorageReady() {
		return nil
	}
	timer := time.NewTimer(timeout)
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-timer.C:
			return fmt.Errorf("timed out waiting for storage ready")
		case <-ticker.C:
			F.Refresh()
			if F.StorageReady() {
				return nil
			}
			switch {
			case strings.HasSuffix(F.UploadState, `Success`):
				return fmt.Errorf("unhandled success uploadState %q", F.UploadState)
			case strings.HasSuffix(F.UploadState, `Failed`):
				return fmt.Errorf("encountered failed uploadState %q", F.UploadState)
			case strings.HasSuffix(F.UploadState, `TimedOut`):
				return fmt.Errorf("encountered timed out uploadState %q", F.UploadState)
			}
		}
	}
}

func (F *MobileAppContentFile) Refresh() error {
	tmp, err := F.graphClient.GetWin32LobAppContentFile(F.Context.AppID, F.Context.ContentVersion, F.ID)
	if err != nil {
		return err
	}
	*F = tmp
	return nil
}

func (F *MobileAppContentFile) StorageReady() bool {
	return F.UploadState == `azureStorageUriRequestSuccess` && F.AzureStorageUri != ""
}

func (F *MobileAppContentFile) UnmarshalJSON(data []byte) error {
	type MACF MobileAppContentFile
	var tmp MACF
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	*F = MobileAppContentFile(tmp)
	F.makeContext()
	return nil
}

func (F *MobileAppContentFile) makeContext() {
	regex := regexp.MustCompile(`/mobileApps\('([a-z0-9-]+)'\)/.*/contentVersions\('([0-9]+)'\)/`)
	m := regex.FindStringSubmatch(F.ODataContext)
	if len(m) == 3 {
		F.Context.AppID = m[1]
		F.Context.ContentVersion = m[2]
	}
}

func (F *MobileAppContentFile) setGraphClient(gC *GraphClient) {
	F.graphClient = gC
}

type MobileAppContentFileRequest MobileAppContentFile
