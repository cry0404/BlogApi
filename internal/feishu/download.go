package feishu

import (
	"BlogApi/config"
	"fmt"
	//"log"
	"os"

	lark "github.com/larksuite/oapi-sdk-go/v3"
)


//需要先获取 filetoken
// 下载素材
// GET https://open.feishu.cn/open-apis/drive/v1/medias/:file_token/download


// 获取临时下载链接
// GET https://open.feishu.cn/open-apis/drive/v1/medias/batch_get_tmp_download_url
func DownloadBookCovers(cfg *config.Config, 
						books []BookRecord, 
						downloadDir string, 
						client lark.Client) error {
					
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return fmt.Errorf("创建下载目录失败")
	}
/*
	for _, book := range books {
		if book.CoverImage == nil {
			log.Printf("书籍《%s》没有封面图片")
			continue
		}


	}*/
	return nil
}

func gettempDownloadURL(client *lark.Client, fileToken string) {

}
//压缩成 webp 结构， 也可以考虑上传到 cf 的存储桶？
