// Copyright 2019 HenryYee.
//
// Licensed under the AGPL, Version 3.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.gnu.org/licenses/agpl-3.0.en.html
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package lib

import (
	"Yearning-go/src/model"
	"crypto/tls"
	"fmt"
	"github.com/labstack/echo/v4"
	"gopkg.in/gomail.v2"
	"net/http"
	"strings"
)

type UserInfo struct {
	ToUser  string
	User    string
	Pawd    string
	Smtp    string
	PubName string
}

var TemoplateTestMail = `
<html>
<body>
	<div style="text-align:center;">
		<h1>Yearning 2.0</h1>
		<h2>此邮件是测试邮件！</h2>
	</div>
</body>
</html>
`

var TmplRejectMail = `
<html>
<body>
<h1>Yearning 工单驳回通知</h1>
<br><p>工单号: %s</p>
<br><p>发起人: %s</p>
<br><p>地址: <a href="%s">%s</a></p>
<br><p>状态: 驳回</p>
<br><p>驳回说明: %s</p>
</body>
</html>
`

var TmplMail = `
<html>
<body>
<h1>Yearning 工单%s通知</h1>
<br><p>工单号: %s</p>
<br><p>发起人: %s</p>
<br><p>地址: <a href="%s">%s</a></p>
<br><p>状态: %s</p>
</body>
</html>
`

var TmplTestDing = `
# Yearning 测试！
`

var TmplReferDing = `# Yearning工单提交通知 #  \n \n  **工单编号:**  %s \n \n **提交人员:**  <font color=\"#78beea\">%s</font> \n \n **审核人员:** <font color=\"#fe8696\">%s</font> \n \n **平台地址:** http://%s \n \n **工单说明:**  %s \n \n **状态:** <font color=\"#1abefa\">已提交</font> \n \n `

var TmplRejectDing = `# Yearning工单驳回通知 #  \n \n  **工单编号:**  %s \n \n **提交人员:**  <font color=\"#78beea\">%s</font> \n \n **审核人员:** <font color=\"#fe8696\">%s</font> \n \n **平台地址:** http://%s \n \n **工单说明:**  %s \n \n **状态:** <font color=\"#df117e\">驳回</font>  \n \n **驳回说明:**  %s `
var TmplSuccessDing = `# Yearning工单执行通知 #  \n \n  **工单编号:**  %s \n \n **提交人员:**  <font color=\"#78beea\">%s</font> \n \n **审核人员:** <font color=\"#fe8696\">%s</font> \n \n **平台地址:** http://%s \n \n **工单说明:**  %s \n \n **状态:** <font color=\"#3fd2bd\">执行成功</font>`
var TmplFailedDing = `# Yearning工单执行通知 #  \n \n  **工单编号:**  %s \n \n **提交人员:**  <font color=\"#78beea\">%s</font> \n \n **审核人员:** <font color=\"#fe8696\">%s</font> \n \n **平台地址:** http://%s \n \n **工单说明:**  %s \n \n **状态:** <font color=\"#ea2426\">执行失败</font>`
var TmplPerformDing = `# Yearning工单转交通知 #  \n \n  **工单编号:**  %s \n \n **提交人员:**  <font color=\"#78beea\">%s</font> \n \n **审核人员:** <font color=\"#fe8696\">%s</font> \n \n **平台地址:** http://%s \n \n **工单说明:**  %s \n \n **状态:** <font color=\"#de4943\">等待执行人执行</font>`
var TmplBackDing = `# Yearning工单执行通知 #  \n \n  **工单编号:**  %s \n \n **提交人员:**  <font color=\"#78beea\">%s</font> \n \n **审核人员:** <font color=\"#fe8696\">%s</font> \n \n **平台地址:** http://%s \n \n **工单说明:**  %s \n \n **状态:** <font color=\"#ea2426\">已撤回</font>`

var TmplQueryRefer = `# Yearning查询申请通知 #  \n \n  **工单编号:**  %s \n \n **提交人员:**  <font color=\"#78beea\">%s</font> \n \n **审核人员:** <font color=\"#fe8696\">%s</font> \n \n **平台地址:** http://%s \n \n **工单说明:**  %s \n \n **状态:** <font color=\"#1abefa\">已提交</font>`
var TmplSuccessQuery = `# Yearning查询申请通知 #  \n \n  **工单编号:**  %s \n \n **提交人员:**  <font color=\"#78beea\">%s</font> \n \n **审核人员:** <font color=\"#fe8696\">%s</font> \n \n **平台地址:** http://%s \n \n **状态:** <font color=\"#3fd2bd\">同意</font>`
var TmplRejectQuery = `# Yearning查询申请通知 #  \n \n  **工单编号:**  %s \n \n **提交人员:**  <font color=\"#78beea\">%s</font> \n \n **审核人员:** <font color=\"#fe8696\">%s</font> \n \n **平台地址:** http://%s \n \n **状态:** <font color=\"#df117e\">已驳回</font>`

var TmplReferDDLWechat = `# %s打算对数据库%s做DDL
> 感觉不对劲，快去找%s
> SQL为%s
`

const YEARNING_ADDRESS = "http://new-yearning.duolainc.com"

func SendMail(c echo.Context, mail model.Message, tmpl string) {
	m := gomail.NewMessage()
	m.SetHeader("From", mail.User)
	m.SetHeader("To", mail.ToUser)
	m.SetHeader("Subject", "Yearning消息推送!")
	m.SetBody("text/html", tmpl)
	d := gomail.NewDialer(mail.Host, mail.Port, mail.User, mail.Password)
	if mail.Ssl {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		c.Logger().Error(err.Error())
		return
	}
}

func SendDingMsg(c echo.Context, msg model.Message, sv string) {
	//请求地址模板

	//创建一个请求
	var mx string

	if msg.PushType {
		mx = fmt.Sprintf(`{"msgtype": "markdown", "markdown": {"content": "%s"}}`, sv)
	} else {
		mx = fmt.Sprintf(`{"msgtype": "markdown", "markdown": {"title": "Yearning sql审计平台", "text": "%s"}}`, sv)
	}

	req, err := http.NewRequest("POST", msg.WebHook, strings.NewReader(mx))
	if err != nil {
		c.Logger().Error(err.Error())
		return
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	//设置请求头
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	//发送请求
	resp, err := client.Do(req)

	if err != nil {
		c.Logger().Error(err.Error())
		return
	}

	//关闭请求
	defer resp.Body.Close()
}

func MessagePush(c echo.Context, workid string, t uint, reject string) {
	var user model.CoreAccount
	var o model.CoreSqlOrder
	var ding, mail, wc string
	model.DB().Where("work_id =?", workid).First(&o)
	model.DB().Select("email").Where("username =?", o.Username).First(&user)
	s := model.GloMessage
	s.ToUser = user.Email

	if t == 0 {
		ding = fmt.Sprintf(TmplRejectDing, o.WorkId, o.Username, o.Assigned, YEARNING_ADDRESS, o.Text, reject)
		mail = fmt.Sprintf(TmplRejectMail, o.WorkId, o.Username, YEARNING_ADDRESS, YEARNING_ADDRESS, reject)
	}

	if t == 1 {
		ding = fmt.Sprintf(TmplSuccessDing, o.WorkId, o.Username, o.Assigned, YEARNING_ADDRESS, o.Text)
		mail = fmt.Sprintf(TmplMail, "执行", o.WorkId, o.Username, YEARNING_ADDRESS, YEARNING_ADDRESS, "执行成功")
	}

	if t == 2 {
		model.DB().Select("email").Where("username =?", o.Assigned).First(&user)
		s.ToUser = user.Email
		ding = fmt.Sprintf(TmplReferDing, o.WorkId, o.Username, o.Assigned, YEARNING_ADDRESS, o.Text)
		mail = fmt.Sprintf(TmplMail, "提交", o.WorkId, o.Username, YEARNING_ADDRESS, YEARNING_ADDRESS, "已提交")
		if o.Type == 0 {
			// ddl
			wc = fmt.Sprintf(TmplReferDDLWechat, o.Username, o.Source, o.Assigned, o.SQL)
		}

	}

	if t == 4 {
		ding = fmt.Sprintf(TmplFailedDing, o.WorkId, o.Username, o.Assigned, YEARNING_ADDRESS, o.Text)
		mail = fmt.Sprintf(TmplMail, "执行", o.WorkId, o.Username, YEARNING_ADDRESS, YEARNING_ADDRESS, "执行失败")
	}

	if t == 5 {
		model.DB().Select("email").Where("username =?", o.Executor).First(&user)
		s.ToUser = user.Email
		ding = fmt.Sprintf(TmplPerformDing, o.WorkId, o.Username, o.Assigned, YEARNING_ADDRESS, o.Text)
		mail = fmt.Sprintf(TmplMail, "转交", o.WorkId, o.Username, YEARNING_ADDRESS, YEARNING_ADDRESS, "等待执行人执行")
	}

	if t == 6 {
		ding = fmt.Sprintf(TmplBackDing, o.WorkId, o.Username, o.Assigned, YEARNING_ADDRESS, o.Text)
		mail = fmt.Sprintf(TmplMail, "提交", o.WorkId, o.Username, YEARNING_ADDRESS, YEARNING_ADDRESS, "已撤销")
	}

	if model.GloOther.Query {
		var op model.CoreQueryOrder
		model.DB().Select("work_id,username,text,assigned").Where("work_id =?", workid).First(&op)
		model.DB().Select("email").Where("username =?", op.Username).First(&user)
		s.ToUser = user.Email
		if t == 6 {
			model.DB().Select("email").Where("username =?", op.Assigned).First(&user)
			s.ToUser = user.Email
			ding = fmt.Sprintf(TmplQueryRefer, op.WorkId, op.Username, op.Assigned, YEARNING_ADDRESS, op.Text)
			mail = fmt.Sprintf(TmplMail, "查询申请", op.WorkId, op.Username, YEARNING_ADDRESS, YEARNING_ADDRESS, "已提交")
		}
		if t == 7 {
			ding = fmt.Sprintf(TmplSuccessQuery, op.WorkId, op.Username, op.Assigned, YEARNING_ADDRESS)
			mail = fmt.Sprintf(TmplMail, "查询申请", op.WorkId, op.Username, YEARNING_ADDRESS, YEARNING_ADDRESS, "已同意")
		}
		if t == 8 {
			ding = fmt.Sprintf(TmplRejectQuery, op.WorkId, op.Username, op.Assigned, YEARNING_ADDRESS)
			mail = fmt.Sprintf(TmplMail, "查询申请", op.WorkId, op.Username, YEARNING_ADDRESS, YEARNING_ADDRESS, "已驳回")
		}
	}

	if model.GloMessage.Mail {
		if user.Email != "" {
			go SendMail(c, s, mail)
		}
	}
	if model.GloMessage.Ding {
		if model.GloMessage.WebHook != "" {
			go SendDingMsg(c, s, ding)

		}
	}

	if wc != "" {
		var leader model.CoreAccount
		model.DB().Where("rule = ? and department like ?", "admin", "%"+user.Department+"%").First(&leader)
		if leader.Username != "" {
			fmt.Println(leader.Username)
			fmt.Println(o.SQL)
			go SentWechatMarkDown("PenZai", wc, GetWechatAppToken())
		}
	}

	fmt.Println(wc)
}
