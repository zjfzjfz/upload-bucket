package upload

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// 分片上传响应结构体
type BlockResponse struct {
	Ctx       string `json:"ctx"`
	Checksum  string `json:"checksum"`
	Offset    int    `json:"offset"`
	Host      string `json:"host"`
	Crc32     int    `json:"crc32"`
	Error     string `json:"error"`
	Code      int    `json:"code"`
	BlockSize int    `json:"blockSize"`
}

// 分片上传文件
func UploadFileSliceV1(uploadToken, filePath, key string) (string, error) {
	const chunkSize = 4 * 1024 * 1024 // 4MB
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return "", err
	}
	fileSize := fileInfo.Size()
	blockCount := (fileSize + chunkSize - 1) / chunkSize

	var blockCtxs []string
	for i := 0; i < int(blockCount); i++ {
		blockSize := chunkSize
		if i == int(blockCount)-1 && fileSize%chunkSize != 0 {
			blockSize = int(fileSize % chunkSize)
		}
		buf := make([]byte, blockSize)
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return "", err
		}
		if n != blockSize {
			return "", fmt.Errorf("read block %d: expected %d bytes, got %d bytes", i, blockSize, n)
		}

		blockCtx, err := uploadBlock(uploadToken, buf)
		if err != nil {
			return "", fmt.Errorf("upload block failed: %v", err)
		}
		blockCtxs = append(blockCtxs, blockCtx)
	}

	return makeFile(uploadToken, key, fileSize, blockCtxs)
}

// 上传块
func uploadBlock(uploadToken string, blockData []byte) (string, error) {
	url := "https://upload.qiniup.com/mkblk/" + strconv.Itoa(len(blockData))
	request, err := http.NewRequest("POST", url, bytes.NewReader(blockData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	request.Header.Add("Authorization", "UpToken "+uploadToken)
	request.Header.Add("Content-Type", "application/octet-stream")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer response.Body.Close()

	var blockResponse BlockResponse
	err = json.NewDecoder(response.Body).Decode(&blockResponse)
	if err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}
	if blockResponse.Code != 200 {
		return "", fmt.Errorf("upload block failed: %s (code: %d)", blockResponse.Error, blockResponse.Code)
	}

	return blockResponse.Ctx, nil
}

// 创建文件
func makeFile(uploadToken, key string, fileSize int64, blockCtxs []string) (string, error) {
	url := "https://upload.qiniup.com/mkfile/" + strconv.FormatInt(fileSize, 10) + "/key/" + base64.URLEncoding.EncodeToString([]byte(key))
	request, err := http.NewRequest("POST", url, bytes.NewReader([]byte(strings.Join(blockCtxs, ","))))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	request.Header.Add("Authorization", "UpToken "+uploadToken)
	request.Header.Add("Content-Type", "text/plain")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer response.Body.Close()

	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	return string(respBody), nil
}
