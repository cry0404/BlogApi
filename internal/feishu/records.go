package feishu

//对于获取到的 records 进行处理
import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	//"log"
	"BlogApi/config"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
)

//飞书客户端的 sdk 实现的真是又臭又长
//列出记录， 可以考虑细化定义错误
// POST /open-apis/bitable/v1/apps/:app_token/tables/:table_id/records/search

//检索记录
// GET /open-apis/bitable/v1/apps/:app_token/tables/:table_id/records/:record_id

func getRecords(client *lark.Client, cfg *config.Config) ([]BookRecord, error) {
	req := larkbitable.
	NewSearchAppTableRecordReqBuilder().
	AppToken(cfg.FeiShu.FeiShuAppToken).
	TableId(cfg.FeiShu.BookTableID).
	PageSize(20).
	Body(larkbitable.NewSearchAppTableRecordReqBodyBuilder().Build()).
	Build()

	resp, err := client.Bitable.V1.
	AppTableRecord.Search(context.Background(), req)

	if err != nil {
		return nil, fmt.Errorf("请求时发生未知错误: %v", err)
	}

	if !resp.Success() {
		
		return nil, fmt.Errorf("API 错误")
	}

	var bookRecords []BookRecord
	// 文本字段一般都是 map[string]any
	for _, item := range resp.Data.Items {
		f := item.Fields
		
		title := parseTextField(f, "书名")
		if title == "" {
			continue  //跳过为空的部分
		}
		author := parseTextField(f, "作者")
		desc := parseTextField(f, "简介")
		comment := parseTextField(f, "书评")

		tag := parseStringField(f, "标签")
		recommend := parseStringField(f, "推荐状态")

		cover := parseFirstFileURL(f, "书籍封面")
		readDate := parseUnixFieldRFC3339(f, "完成阅读时期") // 或本地格式：UnixToLayout(..., "2006-01-02 15:04:05", true)

		bookRecords = append(bookRecords,  BookRecord{
			Title: title,
			Author: author,
			Description: desc,
			Comment: comment,
			Tag: tag,
			RecommendStatus: recommend,
			ConverImage: cover,
			ReadDate: readDate,
			RecordID: *item.RecordId,
		} )
	}



	return bookRecords, nil
}

func parseTextField(fields map[string]any, key string) string {
	v, ok := fields[key]
	if !ok || v == nil {
		return ""
	}
	switch arr := v.(type) {
	case []any:
		var b strings.Builder
		for _, it := range arr {
			if m, ok := it.(map[string]any); ok {
				if s, ok := m["text"].(string); ok {
					b.WriteString(s)
				}
			}
		}
		return b.String()
	case string:
		return arr
	}
	return ""
}

func parseStringField(fields map[string]any, key string) string {
	if v, ok := fields[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func parseFirstFileURL(fields map[string]any, key string) string {
	v, ok := fields[key]
	if !ok || v == nil {
		return ""
	}
	arr, ok := v.([]any)
	if !ok || len(arr) == 0 {
		return ""
	}
	m, ok := arr[0].(map[string]any)
	if !ok {
		return ""
	}
	if s, ok := m["url"].(string); ok && s != "" {
		return s
	}
	if s, ok := m["tmp_url"].(string); ok && s != "" {
		return s
	}
	return ""
}

func parseUnixFieldRFC3339(fields map[string]any, key string) string {
	v, ok := fields[key]
	if !ok || v == nil {
		return ""
	}
	var ts int64
	switch t := v.(type) {
	case float64:
		ts = int64(t)
	case int64:
		ts = t
	case int:
		ts = int64(t)
	case string:
		if p, err := strconv.ParseInt(t, 10, 64); err == nil {
			ts = p
		} else {
			return ""
		}
	default:
		return ""
	}
	sec, nsec := normalizeUnix(ts)
	return time.Unix(sec, nsec).UTC().Format(time.RFC3339)
}

func normalizeUnix(ts int64) (int64, int64) {
	switch {
	case ts >= 1e16: // 纳秒 ns（~1.7e18）
		return ts / 1e9, ts % 1e9
	case ts >= 1e13: // 微秒 μs（~1.7e15）
		return ts / 1e6, (ts % 1e6) * 1e3
	case ts >= 1e10: // 毫秒 ms（~1.7e12）
		return ts / 1e3, (ts % 1e3) * 1e6
	default: // 秒 s（~1.7e9）
		return ts, 0
	}
}
