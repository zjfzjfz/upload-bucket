package conf

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"time"
	"net/url"
)


// 七牛云配置信息
const (
	//Bucket    = "zjf-db1"
	Bucket    = "test-avif"
)

type AccessInfo struct {
	Key      string `json:"key"`
	Secret   string `json:"secret"`
	UserName string
	Password string
	Uid      uint32
}



var(
	UploadToken = generateUploadToken()
//	Authorization = SignQboxToken(User1, "http://rs.qiniu.com/buckets", "")

	User1 = AccessInfo{
		Key:      "4u26TrA3ZdoHwqjX21uXpfgO638K3w7xJ2o5pllE",
		Secret:   "zc3tcNlN4ErDfiU3Aqh9B4z2ndHD95QCXDVFiS8J",
		//Key:      "_k6nwy23UG9PHBabmkfVn47wROJNry-cRx-pelKr",
		//Secret:	  "Bb3665vtPrCYh7NvA7aPYyGbdYfsXZAkJ6oyNgMz",
		UserName: "general_storage_001@test.qiniu.io",
		Password: "Test@123456",
		Uid:      1380469261,
	}
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

	h := hmac.New(sha1.New, []byte(User1.Secret))
	h.Write([]byte(encodedPutPolicy))
	sign := base64.URLEncoding.EncodeToString(h.Sum(nil))

	return User1.Key + ":" + sign + ":" + encodedPutPolicy
}

// QBox Authorization
// APIDOC: https://github.com/qbox/product/blob/master/kodo/auths/Qbox.md
// 只有在 <ContentType> 为 application/x-www-form-urlencoded 时才签进去。
func SignQboxToken(user AccessInfo, uri, body string) string {
	u, err := url.Parse(uri)
	if err != nil {
		println("Parse url failed, url = %d", uri)
	}

	data := u.Path

	if u.RawQuery != "" {
		data += "?" + u.RawQuery
	}
	data += "\n"

	if body != "" {
		data += body
	}

	h := hmac.New(sha1.New, []byte(user.Secret))
	h.Write([]byte(data))
	sign := base64.URLEncoding.EncodeToString(h.Sum(nil))

	return user.Key + ":" + sign
}
