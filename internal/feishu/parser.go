package feishu

import (
	"encoding/json"
	"fmt"
	"log"
)

// 解析飞书响应并提取书籍信息
func ParseBooksFromResponse(jsonData []byte) ([]ProcessedBook, error) {
	var response FeishuResponse[TableListResponse]
	
	if err := json.Unmarshal(jsonData, &response); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("API 返回错误: %s", response.Msg)
	}

	var books []ProcessedBook
	
	for _, item := range response.Data.Items {
		// 跳过空记录
		if len(item.Fields.Title) == 0 {
			log.Printf("跳过空记录: %s", item.RecordID)
			continue
		}

		book := ProcessedBook{
			RecordID: item.RecordID,
		}

		// 提取书名
		if len(item.Fields.Title) > 0 {
			book.Title = item.Fields.Title[0].Text
		}

		// 提取作者
		if len(item.Fields.Author) > 0 {
			book.Author = item.Fields.Author[0].Text
		}

		// 提取书籍封面
		if len(item.Fields.CoverImage) > 0 {
			book.CoverImage = &item.Fields.CoverImage[0]
		}

		// 提取完成阅读时期
		book.ReadDate = item.Fields.ReadDate

		// 提取推荐状态
		book.RecommendStatus = item.Fields.RecommendStatus

		// 提取标签
		book.Tag = item.Fields.Tag

		// 提取简介
		if len(item.Fields.Description) > 0 {
			book.Description = item.Fields.Description[0].Text
		}

		// 提取书评
		if len(item.Fields.Comment) > 0 {
			book.Comment = item.Fields.Comment[0].Text
		}

		books = append(books, book)
	}

	return books, nil
}

// 过滤有效书籍（有书名的记录）
func FilterValidBooks(books []ProcessedBook) []ProcessedBook {
	var validBooks []ProcessedBook
	
	for _, book := range books {
		if book.Title != "" {
			validBooks = append(validBooks, book)
		}
	}
	
	return validBooks
}

// 按推荐状态过滤书籍
func FilterBooksByRecommendStatus(books []ProcessedBook, status string) []ProcessedBook {
	var filteredBooks []ProcessedBook
	
	for _, book := range books {
		if book.RecommendStatus == status {
			filteredBooks = append(filteredBooks, book)
		}
	}
	
	return filteredBooks
}

// 按标签过滤书籍
func FilterBooksByTag(books []ProcessedBook, tag string) []ProcessedBook {
	var filteredBooks []ProcessedBook
	
	for _, book := range books {
		if book.Tag == tag {
			filteredBooks = append(filteredBooks, book)
		}
	}
	
	return filteredBooks
}