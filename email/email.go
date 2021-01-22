package main

import (
	"datacollectoraiot/dao"
	"datacollectoraiot/models/req"
	"datacollectoraiot/models/rsp"
	"datacollectoraiot/pkg/common"
	"datacollectoraiot/pkg/email"
	"datacollectoraiot/pkg/setting"
	"fmt"
	"github.com/wangcheng0509/gpkg/exceptionless"
	"github.com/wangcheng0509/gpkg/try"
	"github.com/xinliangnote/go-util/mail"
	"net/http"
	"net/url"
	"sync"
	"time"
)

/*var uri  = "http://localhost:9000/#/email/template?date=%s";
var upload = "http://127.0.0.1:8080/convert/html2image?u=doctron&p=lampnick&customClip=true&clipWidth=1152&clipHeight=1768&uploadKey=auto_email_%s.png&url=%s";
var download = "https://oss-platform.oss-cn-zhangjiakou.aliyuncs.com/auto_email_%s.png";*/
type Accession struct {
	m sync.Mutex
}

func EmailSendOssData (date string) {

	config := setting.EmailOssConfigs
	_ = exceptionless.Log("开始: "+time.Now().Format("2006-01-02 15:04:05"), false)
	uri  := fmt.Sprintf(config.Uri, date)
	upload := fmt.Sprintf(config.Upload, date, url.PathEscape(uri))
	download := fmt.Sprintf(config.Download, date)
	uploadResp, err := http.Get(upload)
	if err != nil {
		_ = exceptionless.Log("汇报: Email Upload Error: "+ err.Error(), false)
	}
	defer uploadResp.Body.Close()
	downloadResp, err := http.Get(download)
	if err != nil {
		_ = exceptionless.Log("汇报: Email Download Error: "+ err.Error(), false)
	}
	defer downloadResp.Body.Close()

	parse, _ := time.Parse("2006-01-02", date)
	setting.EmailOss.Subject = fmt.Sprintf("%s OSS平台各项指标汇报", parse.Format("2006年01月02日"))
	setting.EmailOss.Body = emailTemplate(parse.Format("2006年01月02日"), uri, download)

	if err = mail.Send(setting.EmailOss); err != nil {
		_ = exceptionless.Log("汇报: Email Send Error: "+ err.Error(), false)
	}

	_ = exceptionless.Log("汇报结束: "+time.Now().Format("2006-01-02 15:04:05"), false)

}

func emailTemplate (date, uri, download string) string {
	return fmt.Sprintf(`

<pre style="color: cadetblue;"> OSS信息汇报如下: %s</pre>
<a href='%s' target='_blank'><img src='%s' /></a>
`,
 date, uri, download)

}

//#region 邮件信息
//go:generate go test  -v ./...
func (access *Accession) EmailHandler(date string) {
	defer access.m.Unlock()
	access.m.Lock()
	startTime := time.Now().Unix()
	var emailReq []*req.EmailReq
	var emailTrendRsp []*rsp.EmailTrendRsp
	var emailHoursTrendRsp []rsp.EmailHoursTrendRsp
	var emailTemplateRsp []rsp.EmailRsp
	var emailTrendTemplateRsp rsp.EmailTrendTemplateRsp
	var emailHoursTrendTemplateRsp rsp.EmailHoursTrendTemplateRsp
	if len(date) <= 0 {
		date = time.Now().Format("2006-01-02")
	}
	try.Try(func() {
		// 指标
		if err := dao.GetStatQuotas(date, &emailReq); err != nil {
			try.Throw(1, err.Error())
		} else {
			emailTemplateRsp = handlerEmailRsp(emailReq)
		}
		// 近7日
		if err := dao.GetTrends(date,30, &emailTrendRsp); err != nil {
			try.Throw(2, err.Error())
		} else {
			for _, data := range emailTrendRsp {
				emailTrendTemplateRsp.Date = append(emailTrendTemplateRsp.Date, data.Date)
				emailTrendTemplateRsp.Onlines = append(emailTrendTemplateRsp.Onlines, data.Onlines)
			}
		}
		// 昨日 昨日+1 昨日+6
		parse, _ := time.Parse("2006-01-02", date)
		day := int(time.Now().Sub(parse).Hours() / 24)
		days := [3]string{time.Now().AddDate(0, 0, -1-day).Format("2006-01-02"), time.Now().AddDate(0, 0, -2-day).Format("2006-01-02"), time.Now().AddDate(0, 0, -7-day).Format("2006-01-02")}
		for i := 0; i < len(days); i++ {
			if err := dao.GetHoursTrends(&emailHoursTrendRsp, days[i]); err != nil {
				try.Throw(3, err.Error())
			} else {
				for _, item := range emailHoursTrendRsp {
					switch i {
					case 0:
						emailHoursTrendTemplateRsp.Date = append(emailHoursTrendTemplateRsp.Date, item.Date)
						emailHoursTrendTemplateRsp.A = append(emailHoursTrendTemplateRsp.A, item.Onlines)
					case 1:
						emailHoursTrendTemplateRsp.B = append(emailHoursTrendTemplateRsp.B, item.Onlines)
					case 2:
						emailHoursTrendTemplateRsp.C = append(emailHoursTrendTemplateRsp.C, item.Onlines)
					}
				}
			}
		}
		email.GetGenerateHtml(date, emailTemplateRsp, emailTrendTemplateRsp, emailHoursTrendTemplateRsp)
		fmt.Printf("--- 耗时: %ds \n", time.Now().Unix() - startTime)
	}).Catch(0, func(exception try.Exception) {
		fmt.Println(exception.Msg)
	}).Finally(func() {})
}

func handlerEmailRsp(emailReq []*req.EmailReq) []rsp.EmailRsp {
	var emailRsp []rsp.EmailRsp
	for _, v := range emailReq {
		day, _ := common.Fix64ToRound(v.Yesterday_count, v.Prev_yesterday_count)
		week, _ := common.Fix64ToRound(v.Week_count, v.Prev_week_count)
		month, _ := common.Fix64ToRound(v.Month_count, v.Prev_month_count)
		emailRsp = append(emailRsp, rsp.EmailRsp{
			Name:  v.Name,
			Count: int64(v.Count),
			Day:   day,
			Week:  week,
			Month: month,
		})
	}
	return emailRsp
}

