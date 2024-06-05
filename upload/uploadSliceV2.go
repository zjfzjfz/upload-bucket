package upload

/*import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	//"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	//"strconv"
)

const (
	uploadURL = "https://upload.qiniup.com"
	chunkSize = 4 * 1024 * 1024 // 4MB
)

func main() {
	accessKey := "<Your Access Key>"
	secretKey := "<Your Secret Key>"
	bucket := "<Your Bucket Name>"
	key := "<Your File Key>"

	filePath := "<Your File Path>"
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	fileName := fileInfo.Name()

	uploadID, err := initiateMultipartUpload(accessKey, secretKey, bucket, key)
	if err != nil {
		log.Fatal(err)
	}

	chunkCount := int(fileSize/chunkSize) + 1
	var etags []string

	for i := 0; i < chunkCount; i++ {
		partNumber := i + 1
		offset := int64(i * chunkSize)
		chunkSize := int64(chunkSize)
		chunkData := make([]byte, chunkSize)

		n, err := file.ReadAt(chunkData, offset)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		chunkData = chunkData[:n]
		etag, err := uploadChunk(accessKey, secretKey, bucket, key, uploadID, partNumber, chunkData)
		if err != nil {
			log.Fatal(err)
		}

		etags = append(etags, etag)
	}

	err = completeMultipartUpload(accessKey, secretKey, bucket, key, uploadID, etags)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("File uploaded successfully:", fileName)
}

func initiateMultipartUpload(accessKey, secretKey, bucket, key string) (string, error) {
	url := fmt.Sprintf("%s/v2/multipart/upload/%s", uploadURL, base64.URLEncoding.EncodeToString([]byte(bucket+"/"+key)))

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Qiniu "+signRequest(accessKey, secretKey, req))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		UploadID string `json:"uploadId"`
	}

	err = decodeJSON(resp.Body, &result)
	if err != nil {
		return "", err
	}

	return result.UploadID, nil
}

func uploadChunk(accessKey, secretKey, bucket, key, uploadID string, partNumber int, chunkData []byte) (string, error) {
	url := fmt.Sprintf("%s/v2/multipart/upload/%s/%d", uploadURL, uploadID, partNumber)

	req, err := http.NewRequest("PUT", url, bytes.NewReader(chunkData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Authorization", "Qiniu "+signRequest(accessKey, secretKey, req))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	etag := resp.Header.Get("Etag")
	return etag, nil
}

func completeMultipartUpload(accessKey, secretKey, bucket, key, uploadID string, etags []string) error {
	url := fmt.Sprintf("%s/v2/multipart/upload/%s/complete", uploadURL, uploadID)

	type part struct {
		PartNumber int    `json:"partNumber"`
		Etag       string `json:"etag"`
	}

	parts := make([]part, len(etags))
	for i, etag := range etags {
		parts[i] = part{
			PartNumber: i + 1,
			Etag:       etag,
		}
	}

	type completeRequest struct {
		Parts []part `json:"parts"`
	}

	reqData := completeRequest{
		Parts: parts,
	}

	req, err := http.NewRequest("POST", url, encodeJSON(reqData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Qiniu "+signRequest(accessKey, secretKey, req))

	client := &http.Client{}
	resp, err :=client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to complete multipart upload: %s", resp.Status)
	}

	return nil
}

func signRequest(accessKey, secretKey string, req *http.Request) string {
	signature := req.Method + " " + req.URL.Path + "\n"

	if req.URL.RawQuery != "" {
		signature += "?" + req.URL.RawQuery + "\n"
	}

	signature += fmt.Sprintf("Host: %s", req.Host)

	hmac := hmacSha1([]byte(secretKey), []byte(signature))
	sign := base64.URLEncoding.EncodeToString(hmac)

	return accessKey + ":" + sign
}

func hmacSha1(key, data []byte) []byte {
	//h := hmac.New(md5.New, key)
	h.Write(data)
	return h.Sum(nil)
}

func encodeJSON(data interface{}) io.Reader {
	//b, _ := json.Marshal(data)
	return bytes.NewReader(b)
}

func decodeJSON(r io.Reader, v interface{}) error {
//	return json.NewDecoder(r).Decode(v)
}*/