package main

import (
	"fmt"

	"upload-bucket/conf"
	"upload-bucket/upload"
	"upload-bucket/manage"
	"time"
)

func main() {
	// 创建bucket
    manage.CreateBucket(conf.User1, conf.Bucket)
    // 稍等一会儿
    time.Sleep(2 * time.Second)
    // 列举所有bucket
    manage.ListBuckets(conf.User1)


	fmt.Println("Upload Token:", conf.UploadToken)

	filePath1 := "/Users/junfengzhou/Desktop/2.avif"
	key1:= "avif2"
	response1, err := upload.UploadFileFormData(conf.UploadToken, filePath1, key1)
	if err != nil {
		fmt.Println("Error uploading file:", err)
	} else {
		fmt.Println("Upload response:", response1)
	}

	manage.GetMimeType(conf.User1, conf.Bucket, key1)

	/*filePath2 := "/Users/junfengzhou/Desktop/5.pdf"
	key2:= "pdf5"
	response2, err := upload.UploadFileSliceV1(conf.UploadToken, filePath2, key2)
	if err != nil {
		fmt.Println("Error uploading file:", err)
	} else {
		fmt.Println("Upload response:", response2)
	}

	filePath3 := "/Users/junfengzhou/Desktop/6.pdf"
	key3:= "pdf6"
	response3, err := upload.UploadFileSliceV1(conf.UploadToken, filePath3, key3)
	if err != nil {
		fmt.Println("Error uploading file:", err)
	} else {
		fmt.Println("Upload response:", response3)
	}*/
}
