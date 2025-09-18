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
		if err := enc.Encode(bookRecord); err != nil {
			return fmt.Errorf("编码 json 时发生错误")
		}
	}
	if err := bw.Flush(); err != nil {
		return fmt.Errorf("刷新失败: %w", err)
	}
	
	return nil
}