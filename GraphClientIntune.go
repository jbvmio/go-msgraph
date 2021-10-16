package msgraph

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
)

// ListWin32LobApps .
func (g *GraphClient) ListWin32LobApps(opts ...ListQueryOption) ([]Win32LobApp, error) {
	resource := "/deviceAppManagement/mobileApps"
	var marsh struct {
		Apps []Win32LobApp `json:"value"`
	}
	err := g.makeGETAPICall(resource, compileListQueryOptions(opts), &marsh)
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
func (g *GraphClient) ListWin32LobAppContentFiles(identifier, version string, opts ...ListQueryOption) ([]MobileAppContentFile, error) {
	resource := fmt.Sprintf("/deviceAppManagement/mobileApps/%s/microsoft.graph.win32LobApp/contentVersions/%s/files", identifier, version)
	var marsh struct {
		Files []MobileAppContentFile `json:"value"`
	}
	err := g.makeGETAPICall(resource, compileListQueryOptions(opts), &marsh)
	return marsh.Files, err
}

// GetWin32LobApp .
func (g *GraphClient) GetWin32LobApp(identifier string, opts ...GetQueryOption) (Win32LobApp, error) {
	resource := fmt.Sprintf("/deviceAppManagement/mobileApps/%s", identifier)
	app := Win32LobApp{}
	err := g.makeGETAPICall(resource, compileGetQueryOptions(opts), &app)
	return app, err
}

// GetWin32LobAppContentFile .
func (g *GraphClient) GetWin32LobAppContentFile(identifier, version, fileID string, opts ...GetQueryOption) (MobileAppContentFile, error) {
	resource := fmt.Sprintf("/deviceAppManagement/mobileApps/%s/microsoft.graph.win32LobApp/contentVersions/%s/files/%s", identifier, version, fileID)
	file := MobileAppContentFile{}
	err := g.makeGETAPICall(resource, compileGetQueryOptions(opts), &file)
	return file, err
}

// CreateWin32LobApp .
func (g *GraphClient) CreateWin32LobApp(app Win32LobApp, opts ...CreateQueryOption) (Win32LobApp, error) {
	resource := "/deviceAppManagement/mobileApps"
	bodyBytes, err := json.Marshal(app)
	if err != nil {
		return app, err
	}
	reader := bytes.NewReader(bodyBytes)
	err = g.makePOSTAPICall(resource, compileCreateQueryOptions(opts), reader, &app)
	return app, err
}

// CreateWin32LobAppContentVersion creates a new Content Version for the given Win32LobApp ID and returns
// the Newly created ContentVersion ID string and any errors.
func (g *GraphClient) CreateWin32LobAppContentVersion(identifier string, opts ...CreateQueryOption) (TypeAndID, error) {
	resource := fmt.Sprintf("/deviceAppManagement/mobileApps/%s/microsoft.graph.win32LobApp/contentVersions", identifier)
	var ID TypeAndID
	err := g.makePOSTAPICall(resource, compileCreateQueryOptions(opts), bytes.NewReader([]byte(`{}`)), &ID)
	return ID, err
}

func (g *GraphClient) DeleteWin32LobAppContentVersion(identifier, version string) error {
	resource := fmt.Sprintf("/deviceAppManagement/mobileApps/%s/microsoft.graph.win32LobApp/contentVersions/%s", identifier, version)
	err := g.makeDELETEAPICall(resource, compileGetQueryOptions([]GetQueryOption{}), nil)
	return err
}

// CreateWin32LobAppContentFile creates a new mobileAppContentFile object.
// https://docs.microsoft.com/en-us/graph/api/intune-apps-mobileappcontentfile-create?view=graph-rest-beta
func (g *GraphClient) CreateWin32LobAppContentFile(identifier, version string, appContentFileReq MobileAppContentFile, opts ...CreateQueryOption) (MobileAppContentFile, error) {
	resource := fmt.Sprintf("/deviceAppManagement/mobileApps/%s/microsoft.graph.win32LobApp/contentVersions/%s/files", identifier, version)
	bodyBytes, err := json.Marshal(appContentFileReq)
	if err != nil {
		return appContentFileReq, err
	}
	reader := bytes.NewReader(bodyBytes)
	var appContentFileResp MobileAppContentFile
	err = g.makePOSTAPICall(resource, compileCreateQueryOptions(opts), reader, &appContentFileResp)
	return appContentFileResp, err
}

// DeleteWin32LobAppContentFile deletes a mobileAppContentFile.
func (g *GraphClient) DeleteWin32LobAppContentFile(identifier, version, fileID string) error {
	resource := fmt.Sprintf("/deviceAppManagement/mobileApps/%s/microsoft.graph.win32LobApp/contentVersions/%s/files/%s", identifier, version, fileID)
	err := g.makeDELETEAPICall(resource, compileGetQueryOptions([]GetQueryOption{}), nil)
	return err
}

// RenewWin32LobAppContentFileUpload renews the SAS URI for an application file upload.
func (g *GraphClient) RenewWin32LobAppContentFileUpload(identifier, version, fileID string) error {
	resource := fmt.Sprintf("/deviceAppManagement/mobileApps/%s/microsoft.graph.win32LobApp/contentVersions/%s/files/%s/renewUpload", identifier, version, fileID)
	err := g.makePOSTAPICall(resource, compileCreateQueryOptions([]CreateQueryOption{}), nil, nil)
	return err
}

func (g *GraphClient) CommitWin32LobAppContentFileUpload(intuneWinFile, identifier, version, fileID string) error {
	resource := fmt.Sprintf("/deviceAppManagement/mobileApps/%s/microsoft.graph.win32LobApp/contentVersions/%s/files/%s/commit", identifier, version, fileID)
	xmlMeta, err := GetIntuneWin32AppMetadata(intuneWinFile, false)
	if err != nil {
		return fmt.Errorf("xmlMeta Error: %w", err)
	}
	encryptionInfo := xmlMeta.EncryptionInfo
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

// RenewWin32LobAppContentFileUpload renews the SAS URI for an application file upload.
func (g *GraphClient) Win32LobAppContentFileUpload(intuneWinFile string, fileContent MobileAppContentFile) error {
	const blocksize = 1024 * 1024 * 100
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

	for i := 0; i < blockCount; i++ {
		blockID := stdBase64Encode(fmt.Sprintf("%04d\n", i))
		start := i * blocksize
		stop := (i + 1) * blocksize
		var block []byte
		switch {
		case stop >= fileSize:
			block = xmlMeta.Data[start:]
		default:
			block = xmlMeta.Data[start:stop]
		}
		err := g.uploadBlock(fileContent.AzureStorageUri, blockID, block)
		if err != nil {
			return fmt.Errorf("error uploading block %d of %d: %v", i, blockCount, err)
		}
	}
	fmt.Println(fileSize, blocksize, blockCount)
	return nil
}

func stdBase64Encode(v string) string {
	return base64.StdEncoding.EncodeToString([]byte(v))
}

func (g *GraphClient) uploadBlock(storageURI, blockID string, data []byte) error {
	params := url.Values{}
	params.Add(`comp`, `block`)
	params.Add(`blockid`, blockID)
	U := storageURI + `&` + params.Encode()
	payload := bytes.NewReader(data)
	req, err := http.NewRequest(`PUT`, U, payload)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Add(`x-ms-blob-type`, `BlockBlob`)
	client := http.Client{}
	resp, err := client.Do(req)
	//resp, err := g.performRaw(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	fmt.Println("Status:", resp.Status, "Code:", resp.StatusCode)
	fmt.Printf("Body: %s\n", b)
	return nil
}

type FileEncryptionInfo struct {
	EncryptionInfo EncryptionInfo `json:"fileEncryptionInfo"`
}
