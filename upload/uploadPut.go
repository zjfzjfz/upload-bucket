package upload

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "mime"
    "net/http"
    "os"
    "path/filepath"
	"encoding/base64"
//	"net/url"
)

func UploadFilePut(upToken, filePath, key string) error {
    file, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    // 获取文件的 MIME 类型
    mimeType := mime.TypeByExtension(filepath.Ext(filePath))
    if mimeType == "" {
        mimeType = "application/octet-stream"
    }

    // 获取文件大小
    fileInfo, err := file.Stat()
    if err != nil {
        return err
    }
    fileSize := fileInfo.Size()

    // 重新读取文件内容
    file.Seek(0, 0)
    fileContent, err := ioutil.ReadAll(file)
    if err != nil {
        return err
    }

    // 构建请求 URL
    url := fmt.Sprintf("http://kodo-dev.up.jfcs-k8s-qa2.qiniu.io/put/%d/key/%s",
        fileSize, base64.URLEncoding.EncodeToString([]byte(key)))

    // 创建请求
    req, err := http.NewRequest("POST", url, bytes.NewReader(fileContent))
    if err != nil {
        return err
    }

    // 设置请求头
    req.Header.Set("Authorization", "UpToken "+upToken)
    req.Header.Set("Content-Type", mimeType)

    // 发送请求
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // 读取响应
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    if resp.StatusCode == http.StatusOK {
        fmt.Println("Upload successful.")
        fmt.Println(string(body))
    } else {
        return fmt.Errorf("upload failed: %s", string(body))
    }

    return nil
}