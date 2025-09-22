package steam

import (
	"bufio"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg" // 添加 JPEG 解码器
	_ "image/png"

	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"time"

	"github.com/chai2010/webp"
)
//也许可以设置一个 downloader 来控制并发量来并发下载 steam 图片
//依据 appid 判断是否存在
var SteamPath = "./public/steam"
func buildDownloadUrl(appid int, icon string) string {

	downloadImageURL := fmt.Sprintf("https://media.steampowered.com/steamcommunity/public/images/apps/%d/%s.jpg", appid, icon)

	return downloadImageURL
}

func DownloadImages(games []Game) error {
	//先 loadindex
	client := &http.Client {
		Timeout: 30 * time.Second,
	}

	if err := os.MkdirAll(SteamPath, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}
	configPath := filepath.Join(SteamPath, "record.ndjson")
	idx, err := loadIndex(configPath)
	if err != nil {
		return fmt.Errorf("载入 steam 的索引信息时失败")
	}

	appids := make([]int, 0, len(games))
	for _, game := range games {
		url := buildDownloadUrl(game.APPID, game.ImgIconURL)

		if _, ok := idx[game.APPID]; ok  {
			continue //如果存在直接跳过
		}
		err := downloadAndConvertSingleImage(client, url, game.APPID)
		if err != nil {
			return fmt.Errorf("下载图片时出现错误: %v", err)
			//log.Printf("下载时发生错误，发生错误的图片是: %v", game.APPID)
		}
		appids = append(appids, game.APPID)
		//将返回的 appid 写入 record.ndjson 文件

	}

	err = appendIndex(appids)

	if err != nil {
		return fmt.Errorf("写入索引时发生错误: %v", err)
	}
	//根据 appid 命名文件然后 json 对应渲染时长等内容， 如果没有找到再进行下载

	return nil
}

func downloadAndConvertSingleImage(client *http.Client, url string, appid int) error {
	resp, err := client.Get(url)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP 错误: %d", resp.StatusCode)
	}

	img, format, err := image.Decode(resp.Body)
	if err != nil {
		return fmt.Errorf("解码图片失败: %w, 格式: %s", err, format)
	}
	// 根据 appid 命名 .webp 图片
	imgWebp := filepath.Join(SteamPath, strconv.Itoa(appid))
	file, err := os.Create(imgWebp+".webp")
	if err != nil {
		return err
	}
	defer file.Close()

	err = webp.Encode(file, img, &webp.Options{
		Quality: 75,
		Lossless: false,
		Exact: false,
	})


	return err
}

//将对应的 map 先注入, 根据 appid 判断是否已经下载过
func loadIndex(path string) (map[int]struct{}, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[int]struct{}{}, nil
		}
		return nil, err
	}
	defer f.Close()

	idx := make(map[int]struct{})
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		var r struct {
			RecordID int `json:"record_id"`
		}
		if err := json.Unmarshal(sc.Bytes(), &r); err == nil  {
			idx[r.RecordID] = struct{}{}
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return idx, nil
}


func appendIndex(appids []int) error {
	filePath := filepath.Join(SteamPath, "record.ndjson")
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	for _, id := range appids {
		if err := enc.Encode(struct {
			RecordID int `json:"record_id"`
		}{RecordID: id}); err != nil {
			return err
		}
	}



	return nil
}