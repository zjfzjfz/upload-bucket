package manage

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
	"upload-bucket/conf"

)

// ObjectInfo 用于解析七牛云返回的对象信息
type ObjectInfo struct {
    MimeType string `json:"mimeType"`
}

// 获取对象的MIME类型
func GetMimeType(user conf.AccessInfo, bucket, fileName string) {
    encodedEntry := base64.URLEncoding.EncodeToString([]byte(bucket + ":" + fileName))
    //uri := "http://rs.qiniu.com/stat/" + encodedEntry
	uri := "http://kodo-dev.rspub.jfcs-k8s-qa2.qiniu.io/stat/" + encodedEntry

    authorization := "QBox " + conf.SignQboxToken(user, uri, "")

    client := &http.Client{}
    req, _ := http.NewRequest("GET", uri, nil)
    req.Header.Set("Authorization", authorization)

    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("请求失败:", err)
        return
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)
    var info ObjectInfo
    if err := json.Unmarshal(body, &info); err != nil {
        fmt.Println("解析响应失败:", err)
        return
    }

    fmt.Println("MIME类型:", info.MimeType)
}