package feishu

import (
	"BlogApi/config"
	"context"
	"fmt"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
)

func getMovieRecords(client *lark.Client, cfg *config.Config) ([]MovieRecord, error) {
	req := larkbitable.
		NewSearchAppTableRecordReqBuilder().
		AppToken(cfg.FeiShu.FeiShuAppToken).
		TableId(cfg.FeiShu.MovieTableID).
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

	var movieRecords []MovieRecord
	// 文本字段一般都是 map[string]any
	for _, item := range resp.Data.Items {
		f := item.Fields

		title := parseTextField(f, "影名")
		if title == "" {
			continue //跳过为空的部分
		}
		author := parseTextField(f, "作者")
		desc := parseTextField(f, "简介")
		comment := parseTextField(f, "影评")
		grade := parseIntField(f, "评价")

		cover := parseFirstFileURL(f, "封面")
		readDate := parseUnixFieldRFC3339(f, "日期")

		movieRecords = append(movieRecords, MovieRecord{
			Title:       title,
			Author:      author,
			Description: desc,
			Comment:     comment,
			Grade:       grade,
			Date:        readDate,
			ConverImage: cover,
			RecordID:    *item.RecordId,
		})
	}

	return movieRecords, nil
}
