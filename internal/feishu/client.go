package feishu

import (
	config "BlogApi/config"
	"fmt"

	"github.com/larksuite/oapi-sdk-go/v3"
)



func UpdateBookCase(cfg *config.Config) error {
	client := lark.NewClient(cfg.FeiShu.FeiShuAppID, cfg.FeiShu.FeiShuAppSecret)

	bookRecords, err := getRecords(client, cfg)
	
	if err != nil {
		return fmt.Errorf("获取书籍记录失败: %v", err)
	}

	
	
	err = downloadImage(cfg, &bookRecords)

	if err != nil {
		return fmt.Errorf("下载图片失败: %v", err)
	}

	//打印测试 
	/*
	for _, bookRecord := range bookRecords {
		fmt.Println(bookRecord)
	}*/

	err = saveAsJson(bookRecords)

	if err != nil {
		return fmt.Errorf("保存为 json 失败: %v", err)
	}
	//下载图片，然后处理， 或者上传到 cdn cf 待实现

	
	fmt.Println("更新成功!开始尝试 webhook")
	return nil
}

