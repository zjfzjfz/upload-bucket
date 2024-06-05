package upload

import (
    "bytes"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io"
	"strings"
    "net/http"
    "os"
)

const (
    UpHost        = "https://upload.qiniup.com"
    ChunkSize     = 4 * 1024 * 1024 // 4MB
    PartSize      = 1024 * 1024            // 每个片的大小
)

// mkblk 创建块
func mkblk(uploadToken string, blockSize int, firstChunk []byte) (string, error) {
    url := fmt.Sprintf("%s/mkblk/%d", UpHost, blockSize)
    req, err := http.NewRequest("POST", url, bytes.NewReader(firstChunk))
    if err != nil {
        return "", err
    }

    req.Header.Set("Authorization", "UpToken "+uploadToken)
    req.Header.Set("Content-Type", "application/octet-stream")
    req.Header.Set("Content-Length", fmt.Sprintf("%d", len(firstChunk)))

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var respData struct {
        Ctx string `json:"ctx"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
        return "", err
    }

    return respData.Ctx, nil
}

// bput 上传片
func bput(uploadToken string, ctx string, offset int, chunk []byte) (string, error) {
    url := fmt.Sprintf("%s/bput/%s/%d", UpHost, ctx, offset)
    req, err := http.NewRequest("POST", url, bytes.NewReader(chunk))
    if err != nil {
        return "", err
    }

    req.Header.Set("Authorization", "UpToken "+uploadToken)
    req.Header.Set("Content-Type", "application/octet-stream")
    req.Header.Set("Content-Length", fmt.Sprintf("%d", len(chunk)))

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var respData struct {
        Ctx string `json:"ctx"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
        return "", err
    }

    return respData.Ctx, nil
}

// mkfile 创建文件
func mkfile(uploadToken string, fileSize int64, key string, ctxs []string) error {
    url := fmt.Sprintf("%s/mkfile/%d/key/%s", UpHost, fileSize, base64.URLEncoding.EncodeToString([]byte(key)))
    req, err := http.NewRequest("POST", url, bytes.NewReader([]byte(strings.Join(ctxs, ","))))
    if err != nil {
        return err
    }

    req.Header.Set("Authorization", "UpToken "+uploadToken)
    req.Header.Set("Content-Type", "text/plain")
    req.Header.Set("Content-Length", fmt.Sprintf("%d", len(strings.Join(ctxs, ","))))

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return fmt.Errorf("mkfile failed with status code: %d", resp.StatusCode)
    }

    return nil
}

// UploadFileSliceV1 分片上传文件
func UploadFileSliceV1(uploadToken, filePath, key string) (string, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return "", err
    }
    defer file.Close()

    fileInfo, err := file.Stat()
    if err != nil {
        return "", err
    }
    fileSize := fileInfo.Size()

    var lastCtxs []string

    for {
        // 读取块
        chunk := make([]byte, ChunkSize)
        readSize, err := file.Read(chunk)
        if err != nil && err != io.EOF {
            return "", err
        }
        if readSize == 0 {
            break
        }
        chunk = chunk[:readSize]

        // 创建块
        ctx, err := mkblk(uploadToken, readSize, chunk[:PartSize])
        if err != nil {
            return "", err
        }

        // 上传剩余的片
        for offset := PartSize; offset < readSize; offset += PartSize {
            end := offset + PartSize
            if end > readSize {
                end = readSize
            }
            nextCtx, err := bput(uploadToken, ctx, offset, chunk[offset:end])
            if err != nil {
                return "", err
            }
            ctx = nextCtx
        }

        lastCtxs = append(lastCtxs, ctx)

        if readSize < ChunkSize {
            break
        }
    }

    // 创建文件
    if err := mkfile(uploadToken, fileSize, key, lastCtxs); err != nil {
        return "", err
    }

    return "上传成功", nil
}
