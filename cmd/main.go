package main

import (
	"log"
	


	"BlogApi/utils"
	config "BlogApi/config"
	feishu "BlogApi/internal/feishu"


)

/*
SDK 使用 client. 业务域. 版本. 资源 .方法名称
来定位具体的 API 方法。以创建文档接口为例，HTTP URL
为 https://open.feishu.cn/open-apis/docx/v1/documents，其中
docx 为业务域，v1 为版本，documents 为资源，
相应的创建方法为 client.Docx.V1.Document.Create()。
*/

// baseURL = https://open.feishu.cn/open-apis/docs/v1/documents
// sdk 会自动管理 tenant_access_token 的生命周期
func main(){

	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("加载配置失败：%v, 请检查环境变量配置", err)
	}

	//feishuClient := lark.NewClient(cfg.FeiShu.FeiShuAppID, cfg.FeiShu.FeiShuAppSecret)
	err = utils.Update(cfg)


	if err != nil {
		log.Fatalf("更新失败: %v, 请检查错误", err)
	}
	
	feishu.ExampleUsage()

}

