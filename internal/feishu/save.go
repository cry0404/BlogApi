package feishu

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)


func saveAsJson( bookRecords []BookRecord) error{
	file, err := os.OpenFile("./feishu.ndjson", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	bw := bufio.NewWriter(file)
	enc := json.NewEncoder(bw)

	enc.SetIndent("", " ")

	for _, bookRecord := range bookRecords {

		if bookRecord.HasDownload {
			//说明这次请求不是第一次下载了
			continue
		}
		bookinfo := convertToInfo(bookRecord)
		if err := enc.Encode(bookinfo); err != nil {
			return fmt.Errorf("编码 json 时发生错误")
		}
	}
	if err := bw.Flush(); err != nil {
		return fmt.Errorf("刷新失败: %w", err)
	}
	
	return nil
}

func convertToInfo(bookRecord BookRecord) BookInfo {
	return BookInfo {
			Title: bookRecord.Title,
			Author: bookRecord.Author,
			ReadDate: bookRecord.ReadDate,
			RecommendStatus: bookRecord.RecommendStatus,
			Tag: bookRecord.Tag,
			Description: bookRecord.Description,
			Comment: bookRecord.Comment,
			ImageDir: bookRecord.ImageDir,
		}
	
}