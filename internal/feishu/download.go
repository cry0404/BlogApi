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

//需要先获取 filetoken
// 下载素材
// GET https://open.feishu.cn/open-apis/drive/v1/medias/:file_token/download

// 获取 token
// POST https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal

// 获取临时下载链接
// GET https://open.feishu.cn/open-apis/drive/v1/medias/batch_get_tmp_download_url

func downloadImage(cfg *config.Config, bookRecords *[]BookRecord) error {
	token, err := getAccessToken(cfg)
	if err != nil {
		return err
	}

	indexPath := filepath.Join(cfg.FeiShu.DownLoadDir, "record.ndjson")
	downloaded, err := loadIndex(indexPath)
	if err != nil {
		return fmt.Errorf("load index failed: %w", err)
	}

	var newly []string
	s := *bookRecords
	for i := range s {
		br := &s[i]
		if _, ok := downloaded[br.RecordID]; ok {
			br.HasDownload = true //这里作为 true
			continue
		}
		if p, ok := findExistingImagePath(br.ConverImage); ok {
			webpPath := imagePathFor(br.ConverImage, ".webp")
			br.HasDownload = true
			if filepath.Ext(p) != ".webp" {
				if _, err := os.Stat(webpPath); err != nil {
					_ = convertToWebp(p, webpPath)
				}
				br.ImageDir = webpPath
				//fmt.Printf("webpath(existing->webp): %s\n", br.ImageDir)
			} else {
				br.ImageDir = p
				//fmt.Printf("webpath(existing): %s\n", br.ImageDir)
			}
			downloaded[br.RecordID] = struct{}{}
			newly = append(newly, br.RecordID)
			continue
		} //说明已经存在了，这里说明以前都下载过

		if p, err := downloadAsWebp(br.ConverImage, token); err != nil {
			continue
		} else {
			br.ImageDir = p
			//fmt.Printf("webpath(downloaded): %s\n", br.ImageDir)
		}
		downloaded[br.RecordID] = struct{}{}
		newly = append(newly, br.RecordID)
	}

	*bookRecords = s

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

func findExistingImagePath(url string) (string, bool) {
	// 尝试常见扩展名
	for _, ext := range []string{".jpg", ".png", ".webp", ".bin"} {
		p := imagePathFor(url, ext)
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

func imagePathFor(u string, ext string) string {
	token := fileTokenFromURL(u)
	if ext == "" {
		ext = ".bin"
	}
	return filepath.Join("public", "bookcase", token+ext)
}

func convertToWebp(srcPath string, dstPath string) error {
	in, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer in.Close()

	// 注册常见解码器，应该是和注册数据库驱动一样的道理？
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

func downloadAsWebp(url, token string) (string, error) {
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
	dst := imagePathFor(url, ".webp")
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
	if err := webpenc.Encode(f, img, &webpenc.Options{Quality: 75}); err != nil {
		return err
	}
	return nil
}
