package feishu

import (
	"BlogApi/config"
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	webpenc "github.com/chai2010/webp"
)

func DownloadImages[T Record](cfg *config.Config, records []T, recordType RecordType) error {
	token, err := getAccessToken(cfg)
	if err != nil {
		return err
	}

	// 根据记录类型确定下载目录
	var downloadDir string
	switch recordType {
	case BookRecordType:
		downloadDir = filepath.Join("public", "bookcase")
	case AnimeRecordType:
		downloadDir = filepath.Join("public", "anime")
	case MovieRecordType:
		downloadDir = filepath.Join("public", "movie")
	default:
		return fmt.Errorf("unsupported record type: %s", recordType)
	}

	indexPath := filepath.Join(downloadDir, "record.ndjson")
	downloaded, err := loadIndex(indexPath)
	if err != nil {
		return fmt.Errorf("load index failed: %w", err)
	}

	var newly []string
	for i := range records {
		record := records[i]
		if _, ok := downloaded[record.GetRecordID()]; ok {
			// 记录已存在：标记已下载，并尽力回填 ImageDir，避免后续写出空值
			record.SetHasDownload(true)
			if p, ok := findExistingImagePath(record.GetCoverImage(), downloadDir); ok {
				webpPath := imagePathFor(record.GetCoverImage(), ".webp", downloadDir)
				if filepath.Ext(p) != ".webp" {
					if _, err := os.Stat(webpPath); err != nil { // 不存在则尝试转换
						_ = convertToWebp(p, webpPath)
					}
					record.SetImageDir(webpPath)
				} else {
					record.SetImageDir(p)
				}
			} else {
				// 没找到现存文件时，按照约定的 token 计算 .webp 目标路径进行回填
				record.SetImageDir(imagePathFor(record.GetCoverImage(), ".webp", downloadDir))
			}
			continue
		}

		if p, ok := findExistingImagePath(record.GetCoverImage(), downloadDir); ok {
			webpPath := imagePathFor(record.GetCoverImage(), ".webp", downloadDir)
			record.SetHasDownload(true)
			if filepath.Ext(p) != ".webp" {
				if _, err := os.Stat(webpPath); err != nil {
					_ = convertToWebp(p, webpPath)
				}
				record.SetImageDir(webpPath)
			} else {
				record.SetImageDir(p)
			}
			downloaded[record.GetRecordID()] = struct{}{}
			newly = append(newly, record.GetRecordID())
			continue
		}

		if p, err := downloadAsWebp(record.GetCoverImage(), token, downloadDir); err != nil {
			continue
		} else {
			record.SetImageDir(p)
		}
		downloaded[record.GetRecordID()] = struct{}{}
		newly = append(newly, record.GetRecordID())
	}

	if len(newly) > 0 {
		if err := appendIndex(indexPath, newly); err != nil {
			return fmt.Errorf("append index failed: %w", err)
		}
	}

	return nil
}

func loadIndex(path string) (map[string]struct{}, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]struct{}{}, nil
		}
		return nil, err
	}
	defer f.Close()

	idx := make(map[string]struct{})
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		var r struct {
			RecordID string `json:"record_id"`
		}
		if err := json.Unmarshal(sc.Bytes(), &r); err == nil && r.RecordID != "" {
			idx[r.RecordID] = struct{}{}
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return idx, nil
}

func appendIndex(path string, ids []string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	for _, id := range ids {
		if err := enc.Encode(struct {
			RecordID string `json:"record_id"`
		}{RecordID: id}); err != nil {
			return err
		}
	}
	return nil
}

func findExistingImagePath(url string, downloadDir string) (string, bool) {
	// 尝试常见扩展名
	for _, ext := range []string{".jpg", ".png", ".webp", ".bin"} {
		p := imagePathFor(url, ext, downloadDir)
		if _, err := os.Stat(p); err == nil {
			return p, true
		}
	}
	return "", false
}

func fileTokenFromURL(u string) string {
	re := regexp.MustCompile(`/medias/([^/]+)/download`)
	if m := re.FindStringSubmatch(u); len(m) == 2 {
		return m[1]
	}
	sum := sha1.Sum([]byte(u))
	return hex.EncodeToString(sum[:])
}

func imagePathFor(u string, ext string, downloadDir string) string {
	token := fileTokenFromURL(u)
	if ext == "" {
		ext = ".bin"
	}
	return filepath.Join(downloadDir, token+ext)
}

func convertToWebp(srcPath string, dstPath string) error {
	in, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer in.Close()

	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("jpg", "jpg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("gif", "gif", gif.Decode, gif.DecodeConfig)

	img, _, err := image.Decode(in)
	if err != nil {
		return err
	}
	return encodeWebp(img, dstPath)
}

func downloadAsWebp(url, token string, downloadDir string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	// 注册常见解码器
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("jpg", "jpg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("gif", "gif", gif.Decode, gif.DecodeConfig)

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return "", fmt.Errorf("decode image failed: %w", err)
	}
	dst := imagePathFor(url, ".webp", downloadDir)
	if err := encodeWebp(img, dst); err != nil {
		return "", err
	}
	return dst, nil
}

func encodeWebp(img image.Image, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := webpenc.Encode(f, img, &webpenc.Options{
		Quality:  75,
		Lossless: false,
		Exact:    false,
	}); err != nil {
		return err
	}
	return nil
}
