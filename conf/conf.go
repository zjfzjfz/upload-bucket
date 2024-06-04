package conf

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"time"
)


// 七牛云配置信息
const (
	AccessKey = "_k6nwy23UG9PHBabmkfVn47wROJNry-cRx-pelKr"
	SecretKey = "Bb3665vtPrCYh7NvA7aPYyGbdYfsXZAkJ6oyNgMz"
	Bucket    = "zjf-db1"
)

var(
	UploadToken = generateUploadToken()
)

// 上传凭证结构体
type PutPolicy struct {
	Scope      string `json:"scope"`
	Deadline   int64  `json:"deadline"`
	ReturnBody string `json:"returnBody"`
}

// 生成上传凭证
func generateUploadToken() string {
	putPolicy := PutPolicy{
		Scope:    Bucket,
		Deadline: time.Now().Unix() + 3600,
	}
	putPolicyJson, _ := json.Marshal(putPolicy)
	encodedPutPolicy := base64.URLEncoding.EncodeToString(putPolicyJson)

	h := hmac.New(sha1.New, []byte(SecretKey))
	h.Write([]byte(encodedPutPolicy))
	sign := base64.URLEncoding.EncodeToString(h.Sum(nil))

	return AccessKey + ":" + sign + ":" + encodedPutPolicy
}