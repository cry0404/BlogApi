package feishu

import (
	"encoding/json"
	"fmt"
	"log"
)

// 使用示例
func ExampleUsage() {
	// 你的 JSON 数据
	jsonData := `{
		"code": 0,
		"data": {
			"has_more": false,
			"items": [
				{
					"fields": {
						"书名": [{"text": "昨日的世界", "type": "text"}],
						"书籍封面": [{
							"file_token": "IjFBb7huMoeddBx5CkacAwLZnKc",
							"name": "昨日的世界.jpg",
							"size": 79791,
							"tmp_url": "https://open.feishu.cn/open-apis/drive/v1/medias/batch_get_tmp_download_url?file_tokens=IjFBb7huMoeddBx5CkacAwLZnKc",
							"type": "image/jpeg",
							"url": "https://open.feishu.cn/open-apis/drive/v1/medias/IjFBb7huMoeddBx5CkacAwLZnKc/download"
						}],
						"书评": [{"text": "就像是身处第二次世界大战的人们并不知道自己已然陷入第二次世界大战...", "type": "text"}],
						"作者": [{"text": "茨威格", "type": "text"}],
						"完成阅读时期": 1754409600000,
						"推荐状态": "推荐",
						"标签": "文学",
						"简介": [{"text": "茨威格回忆录，记述一战前欧洲文化黄金时代与知识分子命运...", "type": "text"}]
					},
					"record_id": "recfnjcWkz"
				}
			],
			"total": 10
		},
		"msg": "success"
	}`

	// 解析书籍数据
	books, err := ParseBooksFromResponse([]byte(jsonData))
	if err != nil {
		log.Fatalf("解析失败: %v", err)
	}

	// 过滤有效书籍
	validBooks := FilterValidBooks(books)
	fmt.Printf("找到 %d 本有效书籍\n", len(validBooks))

	// 打印书籍信息
	for i, book := range validBooks {
		fmt.Printf("\n=== 书籍 %d ===\n", i+1)
		fmt.Printf("书名: %s\n", book.Title)
		fmt.Printf("作者: %s\n", book.Author)
		fmt.Printf("标签: %s\n", book.Tag)
		fmt.Printf("推荐状态: %s\n", book.RecommendStatus)
		fmt.Printf("简介: %s\n", book.Description)
		fmt.Printf("书评: %s\n", book.Comment)
		if book.CoverImage != nil {
			fmt.Printf("封面文件: %s (大小: %d bytes)\n", book.CoverImage.Name, book.CoverImage.Size)
		}
	}

	// 按推荐状态过滤
	recommendedBooks := FilterBooksByRecommendStatus(validBooks, "推荐")
	fmt.Printf("\n推荐书籍数量: %d\n", len(recommendedBooks))

	// 按标签过滤
	literatureBooks := FilterBooksByTag(validBooks, "文学")
	fmt.Printf("文学类书籍数量: %d\n", len(literatureBooks))

	// 输出为 JSON
	jsonOutput, _ := json.MarshalIndent(validBooks, "", "  ")
	fmt.Printf("\nJSON 输出:\n%s\n", string(jsonOutput))
}