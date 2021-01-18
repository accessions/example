package html

import (
	"datacollectoraiot/models/rsp"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"os"
	"path/filepath"
)
var templatePath = "./template"

func GetGenerateHtml(date string, emailRsp []rsp.EmailRsp, emailTrends rsp.EmailTrendTemplateRsp, emailTrendHours rsp.EmailHoursTrendTemplateRsp) {

	contenstTmp, err := template.ParseFiles(filepath.Join(templatePath, "index.html"))
	if err != nil {
		log.Fatalf("读取模版文件失败, %s", err.Error())
	}
	fileName := filepath.Join(templatePath, "newindex.html")
	generateStaticHtml(contenstTmp, fileName, gin.H{"date":date, "num":0.44, "list": emailRsp, "emailTrends": emailTrends, "emailTrendHours": emailTrendHours})
}

func generateStaticHtml(template *template.Template, fileName string, templateRsp map[string]interface{}) {
	if exist(fileName) {
		err := os.Remove(fileName)
		if err != nil {
			log.Fatalf("移除文件失败, %s", err.Error())
		}
	}
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("打开文件失败, %s", err.Error())
	}
	defer file.Close()
	_ = template.Execute(file, &templateRsp)
}

func exist(fileName string) bool {

	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}
