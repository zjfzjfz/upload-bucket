package upload

import (
    "bytes"
    "crypto/md5"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "strconv"
)

// UploadInitResponse 结构体用来解析初始化上传的响应
type UploadInitResponse struct {
    UploadID string `json:"uploadId"`
}

// UploadPartResponse 结构体用来解析每个分片上传后的响应
type UploadPartResponse struct {
    Etag       string `json:"etag"`
    PartNumber int    `json:"partNumber"`
}

// UploadFileSliceV2 函数，使用七牛云的分片上传
func UploadFileSliceV2(uploadToken, filePath, key string) (string, error) {
    // 分片大小，例如5MB
    chunkSize := 5 * 1024 * 1024

    // 打开文件
    file, err := os.Open(filePath)
    if err != nil {
        return "", err
    }
    defer file.Close()

    // 获取文件的大小
    fileInfo, err := file.Stat()
    if err != nil {
        return "", err
    }
    fileSize := fileInfo.Size()

    // 初始化 multipart 上传
    initURL := fmt.Sprintf("/buckets/<BucketName>/objects/%s/uploads", base64.URLEncoding.EncodeToString([]byte(key)))
    upHost := "<UpHost>"
    resp, err := http.Post(initURL, "application/json", nil)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    var initResp UploadInitResponse
    if err := json.NewDecoder(resp.Body).Decode(&initResp); err != nil {
        return "", err
    }
    uploadID := initResp.UploadID

    // 分片上传
    partURL := fmt.Sprintf("/buckets/<BucketName>/objects/%s/uploads/%s/", base64.URLEncoding.EncodeToString([]byte(key)), uploadID)

    var uploadedParts []UploadPartResponse

    for partNumber := 1; fileSize > 0; partNumber++ {
        // 计算单个分片的大小（最后一个分片可能小于 chunkSize）
        partSize := int(min(fileSize, int64(chunkSize)))
        fileSize -= int64(partSize)

        // 读取分片大小的字节
        partBuffer := make([]byte, partSize)
        _, err := file.Read(partBuffer)
        if err != nil {
            return "", err
        }

        // 计算 MD5
        md5Hash := md5.Sum(partBuffer)
        md5Base64 := base64.StdEncoding.EncodeToString(md5Hash[:])

        // 创建一个请求并设置必要的头
        req, err := http.NewRequest("PUT", upHost+partURL+strconv.Itoa(partNumber), bytes.NewReader(partBuffer))
        if err != nil {
            return "", err
        }
        req.Header.Set("Authorization", "UpToken "+uploadToken)
        req.Header.Set("Content-Type", "application/octet-stream")
        req.Header.Set("Content-MD5", md5Base64)
        req.Header.Set("Content-Length", strconv.Itoa(partSize))

        // 发送请求
        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            return "", err
        }
        defer resp.Body.Close()

        // 检查 HTTP 响应
        if resp.StatusCode != http.StatusOK {
            return "", fmt.Errorf("bad status: %s", resp.Status)
        }

        // 解析响应数据
        var partResp UploadPartResponse
        if err := json.NewDecoder(resp.Body).Decode(&partResp); err != nil {
            return "", err
        }
        uploadedParts = append(uploadedParts, partResp)
    }

    // 合并分片
    completeURL := fmt.Sprintf("%s%s", upHost, partURL)
    completeData := map[string]interface{}{
        "parts":   uploadedParts,
        "fname":   fileInfo.Name(),
        "mimeType": "application/octet-stream",
        // 此处添加其他metadata和自定义变量
    }
    completeBuffer := &bytes.Buffer{}
    json.NewEncoder(completeBuffer).Encode(completeData)
    req, err := http.NewRequest("POST", completeURL, completeBuffer)
    if err != nil {
        return "", err
    }
    req.Header.Set("Authorization", "UpToken "+uploadToken)
    req.Header.Set("Content-Type", "application/json")

    // 发送合并分片的请求
    client := &http.Client{}
    resp, err = client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    // 检查 HTTP 响应
    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("bad status: %s", resp.Status)
    }

    // 返回最终上传响应的ETag（可以是其他有用的响应）
    var finalResponse map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&finalResponse); err != nil {
        return "", err
    }

    etag, ok := finalResponse["etag"].(string)
    if !ok {
        return "", fmt.Errorf("cannot get the final ETag")
    }

    return etag, nil
}

func min(a, b int64) int64 {
    if a < b {
        return a
	}
	return b
}