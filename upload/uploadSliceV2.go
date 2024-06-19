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
    "io/ioutil"
)

// 解析初始化上传的响应
type UploadInitResponse struct {
    UploadID string `json:"uploadId"`
}

// 解析每个分片上传后的响应
type UploadPartResponse struct {
    Etag       string `json:"etag"`
    PartNumber int    `json:"partNumber"`
}

// V2分片上传
func UploadFileSliceV2(uploadToken, filePath, key string) (string, error) {
    // 分片大小
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
    upHost := "http://kodo-dev.up.jfcs-k8s-qa2.qiniu.io"
    initURL := fmt.Sprintf("%s/buckets/test-avif/objects/%s/uploads", upHost, base64.URLEncoding.EncodeToString([]byte(key)))
    
    // 创建一个新的HTTP请求
    req1, err := http.NewRequest("POST", initURL, bytes.NewBuffer([]byte{}))
    if err != nil {
        fmt.Println("创建请求失败:", err)
        return "", err
    }

    // 设置请求头
    req1.Header.Set("Content-Type", "application/json")
    req1.Header.Set("Authorization", "UpToken "+uploadToken)

    // 发送请求
    client1 := &http.Client{}
    resp, err := client1.Do(req1)
    if err != nil {
        fmt.Println("发送请求失败:", err)
        return "", err
    }
    defer resp.Body.Close()

    // 读取并解析响应体
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("读取响应体失败:", err)
        return "", err
    }

    var initResp UploadInitResponse
    if err := json.Unmarshal(body, &initResp); err != nil {
        fmt.Println("解析响应体失败:", err)
        return "", err
    }

    uploadID := initResp.UploadID
    // 分片上传
    partURL := fmt.Sprintf("/buckets/test-avif/objects/%s/uploads/%s/", base64.URLEncoding.EncodeToString([]byte(key)), uploadID)

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

    // 完整的请求URL
    completeURL := fmt.Sprintf("%s%s", upHost, partURL)

    // 构造请求体数据
    completeData := map[string]interface{}{
        "parts": uploadedParts,
    }

    completeBuffer := &bytes.Buffer{}
    json.NewEncoder(completeBuffer).Encode(completeData)
//    fmt.Println(completeBuffer.String()) // 打印编码后的JSON字符串
    // 创建请求
    req, err := http.NewRequest("POST", completeURL, completeBuffer)
    if err != nil {
        fmt.Println("创建请求失败:", err)
        return "", err
    }
    req.Header.Set("Authorization", "UpToken "+uploadToken)
    req.Header.Set("Content-Type", "application/json")

    // 发送请求
    client := &http.Client{}
    resp, err2 := client.Do(req)
    if err2 != nil {
        fmt.Println("发送请求失败:", err2)
        return "", err
    }
    defer resp.Body.Close()

    // 检查HTTP响应状态码
    if resp.StatusCode != http.StatusOK {
        fmt.Println("请求失败, 状态码:", resp.StatusCode)
        return "", err2
    }


    // 返回最终上传响应的ETag
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
