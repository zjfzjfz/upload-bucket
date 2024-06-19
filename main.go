package main

import (
	"fmt"
//	"time"
	"upload-bucket/conf"
	"upload-bucket/upload"
	"upload-bucket/manage"
)

//

func main() {
	// 创建bucket
    manage.CreateBucket(conf.User1, conf.Bucket)
    // 稍等一会儿
    //time.Sleep(2 * time.Second)
    // 列举所有bucket
    //manage.ListBuckets(conf.User1)


	//fmt.Println("Upload Token:", conf.UploadToken)

	// 定义一个包含6个文件路径的slice
    filePaths := []string{
        "/Users/junfengzhou/Desktop/1.avif",
        "/Users/junfengzhou/Desktop/2",
        "/Users/junfengzhou/Desktop/3.AvIf",
        "/Users/junfengzhou/Desktop/4.jpg",
        "/Users/junfengzhou/Desktop/5.avif",
//        "/Users/junfengzhou/Desktop/6.avif",
    }

    /*for i, filePath := range filePaths {
        key := fmt.Sprintf("avif%d", i+1) // 构造key名称，如avif1, avif2, ..., avif6

        response, err := upload.UploadFileFormData(conf.UploadToken, filePath, key)
        if err != nil {
            fmt.Printf("Error uploading file %s: %v\n", filePath, err)
        } else {
            fmt.Printf("Upload response for file %s: %v\n", filePath, response)
        }

        // 获取上传文件的MIME类型
        manage.GetMimeType(conf.User1, conf.Bucket, key)
      
    }*/

	for i, filePath := range filePaths {
        key := fmt.Sprintf("avif%d", i+1) // 构造key名称，如avif1, avif2, ..., avif6

        response, err := upload.UploadFileSliceV2(conf.UploadToken, filePath, key)
        if err != nil {
            fmt.Printf("Error uploading file %s: %v\n", filePath, err)
        } else {
            fmt.Printf("Upload response for file %s: %v\n", filePath, response)
        }

        // 获取上传文件的MIME类型
        manage.GetMimeType(conf.User1, conf.Bucket, key)
      
    }

	/*for i, filePath := range filePaths {
        key := fmt.Sprintf("avif%d", i+1) // 构造key名称，如avif1, avif2, ..., avif6

        err := upload.UploadFilePut(conf.UploadToken, filePath, key)
        if err != nil {
            fmt.Printf("Error uploading file %s: %v\n", filePath, err)
        }
        // 获取上传文件的MIME类型
        manage.GetMimeType(conf.User1, conf.Bucket, key)
      
    }*/

	
	
	
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
