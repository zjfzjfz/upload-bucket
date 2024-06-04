package upload

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"

)


// 上传文件
func UploadFileFormData(uploadToken string, filePath string, key string) (string, error) {
	url := "https://upload.qiniup.com"
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("token", uploadToken)
	writer.WriteField("key", key)
	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		return "", err
	}
	part.Write(fileData)
	writer.Close()

	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", err
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}