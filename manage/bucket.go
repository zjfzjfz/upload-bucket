package manage

import (
    "fmt"
    "io/ioutil"
    "net/http"
	"upload-bucket/conf"
)

// 创建bucket
func CreateBucket(user conf.AccessInfo, bucketName string) {
    uri := "http://rs.qiniu.com/mkbucketv3/" + bucketName + "/region/z0"
    method := "POST"
    authorization := "QBox " + conf.SignQboxToken(user, uri, "")

    // 发送请求
    client := &http.Client{}
    req, _ := http.NewRequest(method, uri, nil)
    req.Header.Add("Authorization", authorization)
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("创建bucket请求失败:", err)
        return
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("创建bucket响应:", string(body))
}

// 列举所有的bucket
func ListBuckets(user conf.AccessInfo) {
    uri := "http://rs.qbox.me/buckets"
    method := "GET"
    authorization := "QBox " + conf.SignQboxToken(user, uri, "")

    // 发送请求
    req, _ := http.NewRequest(method, uri, nil)
    req.Header.Set("Authorization", authorization)
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("列举bucket请求失败:", err)
        return
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("Bucket 列表:", string(body))
}
