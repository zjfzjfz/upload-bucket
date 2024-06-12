package manage

import (
    "fmt"
    "io/ioutil"
    "net/http"
	"upload-bucket/conf"
)

// 创建bucket
func CreateBucket(user conf.AccessInfo, bucketName string) {
    //uri := "http://uc.qiniuapi.com/mkbucketv3/" + bucketName + "/region/z0"
	uri := "http://kodo-dev.bucket.jfcs-k8s-qa2.qiniu.io/mkbucketv3/" + bucketName + "/region/z0"
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
    //uri := "http://uc.qiniuapi.com/buckets"
	uri := "http://kodo-dev.bucket.jfcs-k8s-qa2.qiniu.io/buckets"
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

// ListBuckets 列举所有的bucket
/*func ListBuckets(user conf.AccessInfo) {
    uri := "http://kodo-dev.bucket.jfcs-k8s-qa2.qiniu.io"
    method := "POST"

    // 创建一个multipart表单的body
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    // 假设我们需要传递一个参数（虽然在GET请求转POST的情况下可能并不需要）
    _ = writer.WriteField("paramName", "paramValue")
    // 关闭writer将会写入结尾的边界
    err := writer.Close()
    if err != nil {
        fmt.Println("构造表单数据失败:", err)
        return
    }

    // 使用SignQboxToken函数生成Authorization头部
    // 注意：在调用SignQboxToken时，不应该传递body，因为GET请求不包含body
    authorization := "QBox " + conf.SignQboxToken(user, uri, "")

    // 发送请求
    req, err := http.NewRequest(method, uri, body)
    if err != nil {
        fmt.Println("创建请求失败:", err)
        return
    }
    req.Header.Set("Authorization", authorization)
    req.Header.Set("Content-Type", writer.FormDataContentType())

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("列举bucket请求失败:", err)
        return
    }
    defer resp.Body.Close()

    responseBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("读取响应体失败:", err)
        return
    }
    fmt.Println("Bucket 列表:", string(responseBody))
}*/
