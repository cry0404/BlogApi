package feishu

//定义下载映射，将对应的 table 映射到对应的路径
import (
	"BlogApi/config"
	"context"
	"fmt"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
)

func getAnimeRecords(client *lark.Client, cfg *config.Config) ([]AnimeRecord, error) {
	req := larkbitable.
		NewSearchAppTableRecordReqBuilder().
		AppToken(cfg.FeiShu.FeiShuAppToken).
		TableId(cfg.FeiShu.AnimeTableID).
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

	var animeRecords []AnimeRecord
	// 文本字段一般都是 map[string]any
	for _, item := range resp.Data.Items {
		f := item.Fields

		title := parseTextField(f, "番名")
		if title == "" {
			continue //跳过为空的部分
		}

		desc := parseTextField(f, "简介")
		comment := parseTextField(f, "评价")
		grade := parseIntField(f, "评分")

		cover := parseFirstFileURL(f, "封面")

		animeRecords = append(animeRecords, AnimeRecord{
			Title:       title,
			Description: desc,
			Grade:       grade,
			Comment:     comment,
			ConverImage: cover,
			RecordID:    *item.RecordId,
		})
	}

	return animeRecords, nil
}
