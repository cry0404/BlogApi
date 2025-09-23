package feishu

//对于获取到的 records 进行处理
import (
	"context"
	"fmt"
	

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

func getBookRecords(client *lark.Client, cfg *config.Config) ([]BookRecord, error) {
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

		cover := parseFirstFileURL(f, "封面")
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


