package feishu

import (
	config "BlogApi/config"
	"fmt"

	lark "github.com/larksuite/oapi-sdk-go/v3"
)

func UpdateFeiShu(cfg *config.Config) error {
	client := lark.NewClient(cfg.FeiShu.FeiShuAppID, cfg.FeiShu.FeiShuAppSecret)

	// 处理书籍记录
	bookRecords, err := getBookRecords(client, cfg)
	if err != nil {
		return fmt.Errorf("获取书籍记录失败: %v", err)
	}

	bookRecordsInterface := make([]Record, len(bookRecords))
	for i := range bookRecords {
		bookRecordsInterface[i] = &bookRecords[i]
	}

	err = DownloadImages(cfg, bookRecordsInterface, BookRecordType)
	if err != nil {
		return fmt.Errorf("下载书架图片失败: %v", err)
	}

	err = SaveAsJSON(bookRecordsInterface, BookRecordType)
	if err != nil {
		return fmt.Errorf("保存书籍为 json 失败: %v", err)
	}

	// 处理动漫记录
	animeRecords, err := getAnimeRecords(client, cfg)
	if err != nil {
		return fmt.Errorf("获取动漫记录失败: %v", err)
	}

	animeRecordsInterface := make([]Record, len(animeRecords))
	for i := range animeRecords {
		animeRecordsInterface[i] = &animeRecords[i]
	}

	err = DownloadImages(cfg, animeRecordsInterface, AnimeRecordType)
	if err != nil {
		return fmt.Errorf("下载动漫图片失败: %v", err)
	}

	err = SaveAsJSON(animeRecordsInterface, AnimeRecordType)
	if err != nil {
		return fmt.Errorf("保存动漫为 json 失败: %v", err)
	}

	// 处理电影记录
	movieRecords, err := getMovieRecords(client, cfg)
	if err != nil {
		return fmt.Errorf("获取电影记录失败: %v", err)
	}

	movieRecordsInterface := make([]Record, len(movieRecords))
	for i := range movieRecords {
		movieRecordsInterface[i] = &movieRecords[i]
	}

	err = DownloadImages(cfg, movieRecordsInterface, MovieRecordType)
	if err != nil {
		return fmt.Errorf("下载电影图片失败: %v", err)
	}

	err = SaveAsJSON(movieRecordsInterface, MovieRecordType)
	if err != nil {
		return fmt.Errorf("保存电影为 json 失败: %v", err)
	}

	return nil
}
