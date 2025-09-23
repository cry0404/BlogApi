package feishu

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// SaveAsJSON 通用的保存为 JSON 函数，支持所有类型的记录
func SaveAsJSON[T Record](records []T, recordType RecordType) error {
	var fileName string
	switch recordType {
	case BookRecordType:
		fileName = "./config/book.ndjson"
	case AnimeRecordType:
		fileName = "./config/anime.ndjson"
	case MovieRecordType:
		fileName = "./config/movie.ndjson"
	default:
		return fmt.Errorf("unsupported record type: %s", recordType)
	}

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(fileName), 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	bw := bufio.NewWriter(file)
	enc := json.NewEncoder(bw)
	enc.SetIndent("", " ")

	for _, record := range records {
		if record.GetHasDownload() {
			// 说明这次请求不是第一次下载了
			continue
		}

		info := record.ToOutputInfo()
		if err := enc.Encode(info); err != nil {
			return fmt.Errorf("编码 json 时发生错误: %w", err)
		}
	}

	if err := bw.Flush(); err != nil {
		return fmt.Errorf("刷新失败: %w", err)
	}

	return nil
}
