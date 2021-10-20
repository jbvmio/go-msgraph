package msgraph

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ListWin32LobApps .
func (g *GraphClient) ListWin32LobApps(opts ...ListQueryOption) (Win32LobApps, error) {
	resource := "/deviceAppManagement/mobileApps"
	var marsh struct {
		Apps Win32LobApps `json:"value"`
	}
	err := g.makeGETAPICall(resource, compileListQueryOptions(opts), &marsh)
	marsh.Apps.setGraphClient(g)
	return marsh.Apps, err
}

// ListWin32LobAppContentVersions .
func (g *GraphClient) ListWin32LobAppContentVersions(identifier string, opts ...ListQueryOption) ([]TypeAndID, error) {
	resource := fmt.Sprintf("/deviceAppManagement/mobileApps/%s/microsoft.graph.win32LobApp/contentVersions", identifier)
	var marsh struct {
		Versions []TypeAndID `json:"value"`
	}
	err := g.makeGETAPICall(resource, compileListQueryOptions(opts), &marsh)
	return marsh.Versions, err
}

// ListWin32LobAppContentFiles lists properties and relationships of the mobileAppContentFile objects.
func (g *GraphClient) ListWin32LobAppContentFiles(identifier, version string, opts ...ListQueryOption) (MobileAppContentFiles, error) {
	resource := fmt.Sprintf("/deviceAppManagement/mobileApps/%s/microsoft.graph.win32LobApp/contentVersions/%s/files", identifier, version)
	var marsh struct {
		Files MobileAppContentFiles `json:"value"`
	}
	err := g.makeGETAPICall(resource, compileListQueryOptions(opts), &marsh)
	marsh.Files.setGraphClient(g)
	return marsh.Files, err
}

// GetWin32LobApp .
func (g *GraphClient) GetWin32LobApp(identifier string, opts ...GetQueryOption) (Win32LobApp, error) {
	resource := fmt.Sprintf("/deviceAppManagement/mobileApps/%s", identifier)
	app := Win32LobApp{}
	err := g.makeGETAPICall(resource, compileGetQueryOptions(opts), &app)
	app.setGraphClient(g)
	return app, err
}

// GetWin32LobAppContentFile .
func (g *GraphClient) GetWin32LobAppContentFile(appID, contentVersionID, contentFileID string, opts ...GetQueryOption) (MobileAppContentFile, error) {
	resource := fmt.Sprintf("/deviceAppManagement/mobileApps/%s/microsoft.graph.win32LobApp/contentVersions/%s/files/%s", appID, contentVersionID, contentFileID)
	file := MobileAppContentFile{}
	err := g.makeGETAPICall(resource, compileGetQueryOptions(opts), &file)
	file.setGraphClient(g)
	return file, err
}

func (g *GraphClient) DeleteWin32LobApp(appID string) error {
	resource := fmt.Sprintf("/deviceAppManagement/mobileApps/%s", appID)
	err := g.makeDELETEAPICall(resource, compileGetQueryOptions([]GetQueryOption{}), nil)
	return err
}

func (g *GraphClient) DeleteWin32LobAppContentVersion(appID, contentVersionID string) error {
	resource := fmt.Sprintf("/deviceAppManagement/mobileApps/%s/microsoft.graph.win32LobApp/contentVersions/%s", appID, contentVersionID)
	err := g.makeDELETEAPICall(resource, compileGetQueryOptions([]GetQueryOption{}), nil)
	return err
}

// DeleteWin32LobAppContentFile deletes a mobileAppContentFile.
func (g *GraphClient) DeleteWin32LobAppContentFile(appID, contentVersionID, contentFileID string) error {
	resource := fmt.Sprintf("/deviceAppManagement/mobileApps/%s/microsoft.graph.win32LobApp/contentVersions/%s/files/%s", appID, contentVersionID, contentFileID)
	err := g.makeDELETEAPICall(resource, compileGetQueryOptions([]GetQueryOption{}), nil)
	return err
}

// CreateWin32LobApp Submits and Creates a New Win32LobApp in Intune.
func (g *GraphClient) CreateWin32LobApp(req Win32LobAppRequest, opts ...CreateQueryOption) (Win32LobApp, error) {
	resource := "/deviceAppManagement/mobileApps"
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return Win32LobApp(req), err
	}
	reader := bytes.NewReader(bodyBytes)
	var app Win32LobApp
	err = g.makePOSTAPICall(resource, compileCreateQueryOptions(opts), reader, &app)
	app.setGraphClient(g)
	return app, err
}

// CreateWin32LobAppContentVersion creates a new Content Version for the given Win32LobApp ID and returns
// the Newly created ContentVersion ID string and any errors.
func (g *GraphClient) CreateWin32LobAppContentVersion(appID string, opts ...CreateQueryOption) (TypeAndID, error) {
	resource := fmt.Sprintf("/deviceAppManagement/mobileApps/%s/microsoft.graph.win32LobApp/contentVersions", appID)
	var ID TypeAndID
	err := g.makePOSTAPICall(resource, compileCreateQueryOptions(opts), bytes.NewReader([]byte(`{}`)), &ID)
	return ID, err
}

// CreateWin32LobAppContentFile creates a new mobileAppContentFile object.
// https://docs.microsoft.com/en-us/graph/api/intune-apps-mobileappcontentfile-create?view=graph-rest-beta
func (g *GraphClient) CreateWin32LobAppContentFile(appID, contentVersionID string, appContentFileReq MobileAppContentFileRequest, opts ...CreateQueryOption) (MobileAppContentFile, error) {
	resource := fmt.Sprintf("/deviceAppManagement/mobileApps/%s/microsoft.graph.win32LobApp/contentVersions/%s/files", appID, contentVersionID)
	bodyBytes, err := json.Marshal(appContentFileReq)
	if err != nil {
		return MobileAppContentFile(appContentFileReq), err
	}
	reader := bytes.NewReader(bodyBytes)
	var appContentFileResp MobileAppContentFile
	err = g.makePOSTAPICall(resource, compileCreateQueryOptions(opts), reader, &appContentFileResp)
	appContentFileResp.setGraphClient(g)
	return appContentFileResp, err
}

// RenewWin32LobAppContentFileUpload uploads the given intunewin file to given MobileAppContentFile definition.
func (g *GraphClient) Win32LobAppContentFileUpload(intuneWinFile io.Reader, fileContent *MobileAppContentFile) error {
	const blocksize = 1024 * 1024 * 100
	err := fileContent.ContinueOnUploadState(UploadStateStorageReady, time.Minute*1)
	if err != nil {
		return fmt.Errorf("az storage wait error (storageReady): %w", err)
	}
	if fileContent.AzureStorageUri == "" {
		return fmt.Errorf("missing AzureStorageURI")
	}
	xmlMeta, err := GetIntuneWin32AppMetadata(intuneWinFile, true)
	if err != nil {
		return fmt.Errorf("xmlMeta Error: %w", err)
	}
	fileSize := len(xmlMeta.Data)
	if fileSize == 0 {
		return fmt.Errorf("file %q is invalid or missing data", intuneWinFile)
	}
	blockCount := int(math.Ceil(float64(fileSize) / blocksize))
	doneChan := make(chan error, blockCount)
	ticker := time.NewTicker(time.Minute * 12)
	blockIDs := make([]string, blockCount)
	var errs []error
	var done int
	for i := 0; i < blockCount; i++ {
		blockID := stdBase64Encode(fmt.Sprintf("%04d\n", i))
		blockIDs[i] = blockID
		start := i * blocksize
		stop := (i + 1) * blocksize
		var block []byte
		switch {
		case stop >= fileSize:
			block = xmlMeta.Data[start:]
		default:
			block = xmlMeta.Data[start:stop]
		}
		go win32LobAppUploadBlock(fileContent.AzureStorageUri, blockID, block, doneChan)
	}
uploadLoop:
	for {
		select {
		case <-ticker.C:
			err := g.Win32LobAppContentFileUploadRenew(fileContent.Context.AppID, fileContent.Context.ContentVersion, fileContent.ID)
			if err != nil {
				return fmt.Errorf("error renewing upload: %w", err)
			}
		case err := <-doneChan:
			done++
			if err != nil {
				errs = append(errs, err)
			}
			if done >= blockCount {
				break uploadLoop
			}
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("error uploading %d of %d: %v", len(errs), blockCount, errs)
	}
	fileContent.Refresh()
	err = win32LobAppUploadFinalize(fileContent.AzureStorageUri, blockIDs)
	if err != nil {
		return fmt.Errorf("error finalizing upload: %w", err)
	}
	err = g.Win32LobAppContentFileUploadCommit(xmlMeta.EncryptionInfo, fileContent.Context.AppID, fileContent.Context.ContentVersion, fileContent.ID)
	if err != nil {
		return fmt.Errorf("error committing upload: %w", err)
	}
	err = fileContent.ContinueOnUploadState(UploadStateFileCommitted, time.Minute*1)
	if err != nil {
		return fmt.Errorf("az storage wait error (fileCommitted): %w", err)
	}
	fileContent.Refresh()
	err = g.Win32LobAppContentFileVersionCommit(fileContent.Context.AppID, fileContent.Context.ContentVersion)
	if err != nil {
		return fmt.Errorf("error committing content version: %w", err)
	}
	return nil
}

// Win32LobAppContentFileUploadRenew renews the SAS URI for an application file upload.
func (g *GraphClient) Win32LobAppContentFileUploadRenew(appID, contentVersionID, contentFileID string) error {
	resource := fmt.Sprintf("/deviceAppManagement/mobileApps/%s/microsoft.graph.win32LobApp/contentVersions/%s/files/%s/renewUpload", appID, contentVersionID, contentFileID)
	err := g.makePOSTAPICall(resource, compileCreateQueryOptions([]CreateQueryOption{}), nil, nil)
	return err
}

func (g *GraphClient) Win32LobAppContentFileUploadCommit(encryptionInfo EncryptionInfo, appID, contentVersionID, contentFileID string) error {
	resource := fmt.Sprintf("/deviceAppManagement/mobileApps/%s/microsoft.graph.win32LobApp/contentVersions/%s/files/%s/commit", appID, contentVersionID, contentFileID)
	if encryptionInfo.ProfileIdentifier == "" {
		encryptionInfo.ProfileIdentifier = `ProfileVersion1`
	}
	encryption := FileEncryptionInfo{
		EncryptionInfo: encryptionInfo,
	}
	bodyBytes, err := json.Marshal(encryption)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(bodyBytes)
	err = g.makePOSTAPICall(resource, compileCreateQueryOptions([]CreateQueryOption{}), reader, nil)
	return err
}

func (g *GraphClient) Win32LobAppContentFileVersionCommit(appID, contentVersionID string) error {
	resource := fmt.Sprintf("/deviceAppManagement/mobileApps/%s", appID)
	tmp := struct {
		ODataType               string `json:"@odata.type"`
		CommittedContentVersion string `json:"committedContentVersion"`
	}{
		ODataType:               `#microsoft.graph.win32LobApp`,
		CommittedContentVersion: contentVersionID,
	}
	bodyBytes, err := json.Marshal(tmp)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(bodyBytes)
	err = g.makePATCHAPICall(resource, compileGetQueryOptions([]GetQueryOption{}), reader, nil)
	return err
}

func stdBase64Encode(v string) string {
	return base64.StdEncoding.EncodeToString([]byte(v))
}

func win32LobAppUploadBlock(storageURI, blockID string, data []byte, doneChan chan error) {
	params := url.Values{}
	params.Add(`comp`, `block`)
	params.Add(`blockid`, blockID)
	U := storageURI + `&` + params.Encode()
	payload := bytes.NewReader(data)
	req, err := http.NewRequest(`PUT`, U, payload)
	if err != nil {
		doneChan <- fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Add(`x-ms-blob-type`, `BlockBlob`)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		doneChan <- fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		doneChan <- fmt.Errorf("error reading response: %w", err)
	}
	//fmt.Println("Upload Status:", resp.Status, "Code:", resp.StatusCode)
	doneChan <- nil
}

func win32LobAppUploadFinalize(storageURI string, blockIDs []string) error {
	params := url.Values{}
	params.Add(`comp`, `blocklist`)
	U := storageURI + `&` + params.Encode()
	xml := `<?xml version="1.0" encoding="utf-8"?><BlockList>`
	for _, id := range blockIDs {
		xml += fmt.Sprintf("<Latest>%s</Latest>", id)
	}
	xml += `</BlockList>`
	payload := strings.NewReader(xml)
	req, err := http.NewRequest(`PUT`, U, payload)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}
	//fmt.Println("Finalize Status:", resp.Status, "Code:", resp.StatusCode)
	return nil
}

type FileEncryptionInfo struct {
	EncryptionInfo EncryptionInfo `json:"fileEncryptionInfo"`
}
