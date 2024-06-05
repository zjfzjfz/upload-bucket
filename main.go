package main

import (
	"fmt"

	"upload-bucket/conf"
	"upload-bucket/upload"
)

func main() {
	fmt.Println("Upload Token:", conf.UploadToken)

	//filePath := "D:/goproject/src/upload-bucket/2.txt"

	/*filePath1 := "/Users/junfengzhou/Desktop/1.jpg"
	key1:= "jpg1"

	response1, err := upload.UploadFileFormData(conf.UploadToken, filePath1, key1)
	if err != nil {
		fmt.Println("Error uploading file:", err)
	} else {
		fmt.Println("Upload response:", response1)
	}

	filePath2 := "/Users/junfengzhou/Desktop/5.pdf"
	key2:= "pdf5"

	response2, err := upload.UploadFileSliceV1(conf.UploadToken, filePath2, key2)
	if err != nil {
		fmt.Println("Error uploading file:", err)
	} else {
		fmt.Println("Upload response:", response2)
	}*/

	filePath3 := "/Users/junfengzhou/Desktop/6.pdf"
	key3:= "pdf6"

	response3, err := upload.UploadFileSliceV1(conf.UploadToken, filePath3, key3)
	if err != nil {
		fmt.Println("Error uploading file:", err)
	} else {
		fmt.Println("Upload response:", response3)
	}
}
