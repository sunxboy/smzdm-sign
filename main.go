package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"
)

var (
	conf    *config = &config{}
	loc     *time.Location
	startAt time.Time
)

func (c *account) signIn() ([]byte, error) {
	ts := time.Now().In(loc).UnixNano()
	url := fmt.Sprintf("https://zhiyou.smzdm.com/user/checkin/jsonp_checkin?callback=jQuery112409568846254764496_%d&_=%d", ts, ts)

	content, err := c.handle("GET", url, "https://www.smzdm.com/", nil)
	if err != nil {
		log.Errorf("smzdm sign-in error: %v", err)
		return nil, err
	}

	return content, nil
}

func (c *account) handle(method, url, referer string, body *io.Reader) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Warningf("could not access url: %s, get error: %v \n", url, err)
		return nil, err
	}

	header := http.Header{}
	if referer == "" {
		referer = "https://www.smzdm.com"
	}
	header.Add("Referer", referer)
	header.Add("User-Agent", c.UserAgent)
	header.Add("Cookie", c.Cookies)
	req.Header = header
	resp, err := client.Do(req)
	if err != nil {
		log.Warningf("invalid header when access url: %s, get error: %v", url, err)
		return nil, err
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Warningf("cannot get response when access url: %s, get error: %v", url, err)
		return nil, err
	}
	return content, nil
}

func parseAsMailBody(responseData []byte) ([]byte, error) {
	jsonData := &signJson{
		Index:   1,
		Account: conf.Account,
		Time:    time.Now().In(loc),
	}
	reg := regexp.MustCompile(`^jQuery\d+_\d+\((.*?)\)$`)
	jsonStr := reg.ReplaceAll(responseData, []byte(`$1`))
	err := json.Unmarshal(jsonStr, jsonData)
	if err != nil {
		errContent := fmt.Sprintf("sign-in passed, but got unknown response: %v", string(responseData))
		sendMailError(errContent)
		return nil, err
	}

	if jsonData.ErrorCode != 0 {
		errContent := fmt.Sprintf("sign-in passed, but is error response: %v", jsonStr)
		sendMailError(errContent)
		return nil, err
	}

	mailBody, err := formatMailBody(*jsonData)
	if err != nil {
		errContent := fmt.Sprintf("gen mail content failed: %v", err)
		sendMailError(errContent)
		return nil, err
	}
	return mailBody, nil
}

func formatMailBody(json signJson) ([]byte, error) {
	var mailContent bytes.Buffer
	tpl, err := template.New("mail.ghtml").Funcs(template.FuncMap{
		"join":     strings.Join,
		"dateTime": time.Now().In(loc).Format,
	}).ParseFiles("mail.ghtml")
	if err != nil {
		log.Errorf("read mail tpl failed: %v", err)
		return nil, err
	}
	err = tpl.Execute(&mailContent, struct {
		Conf       config
		Content    string
		StartAt    time.Time
		EndAt      time.Time
		SignResult signJson
	}{
		Conf:       *conf,
		Content:    "Hello SMZDM!!" + time.Now().In(loc).Format(time.Kitchen),
		StartAt:    startAt,
		EndAt:      time.Now().In(loc),
		SignResult: json,
	})
	if err != nil {
		log.Errorf("parse mail tpl failed: %v", err)
		return nil, err
	}
	return mailContent.Bytes(), nil
}

func sendMailError(content string) {
	contentTpl := `From: Smzdm-Auto-Sign<{{.Conf.EmailFrom}}>
To: {{join .Conf.EmailTo ", "}}
Subject: {{.Conf.EmailToSubject}}_{{dateTime "20060102"}}
Content-Type: text/html; charset=utf-8

{{.Content}}
`
	tpl, _ := template.New("mail").Funcs(template.FuncMap{
		"join":     strings.Join,
		"dateTime": time.Now().In(loc).Format,
	}).Parse(contentTpl)

	var mailBody bytes.Buffer
	tpl.Execute(&mailBody, struct {
		Conf    config
		Content string
	}{
		Conf:    *conf,
		Content: content,
	})

	log.Info(mailBody.String())

	sendMail(mailBody.Bytes())
}

func sendMail(content []byte) {
	smtpHost := []byte(conf.EmailFromSMTP)[0:strings.LastIndex(conf.EmailFromSMTP, ":")]
	// 认证
	auth := smtp.PlainAuth("", conf.EmailFrom, conf.EmailFromPassword, string(smtpHost))
	// 发送
	err := smtp.SendMail(conf.EmailFromSMTP, auth, conf.EmailFrom, conf.EmailTo, content)
	if err != nil {
		log.Errorf("send mail failed: %v", err)
	}
}

func main() {
	s1 := gocron.NewScheduler(loc)
	s1.Every(1).Day().At("01:20").Do(smzdmSignIn)
	<-s1.StartAsync()
}

func smzdmSignIn() {
	startAt = time.Now().In(loc)
	zone, offset := startAt.Zone()
	log.Infof("current time zone: %v, offset:  %v", zone, offset)
	z := conf.Account
	content, err := z.signIn()
	if err != nil {
		errContent := fmt.Sprintf("smzdm sign-in failed: %v", err)
		sendMailError(errContent)
		return
	}

	mailBody, err := parseAsMailBody(content)
	if err == nil {
		sendMail(mailBody)
	}
}

func init() {
	loc, _ = time.LoadLocation("Asia/Chongqing")

	configPath := "config.yaml"
	configFile, err := os.Open(configPath)
	if err != nil {
		log.Errorf("config.yaml cannot open: %v", err)
		os.Exit(1)
	}
	defer configFile.Close()

	configStr, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Errorf("config.yaml cannot read: %v", err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(configStr, conf)
	if err != nil {
		log.Errorf("config.yaml format error: %v", err)
		os.Exit(1)
	}

	//setup log

	// with Json Formatter
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	file, err := os.OpenFile(LOG_FILE, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	log.SetOutput(file)
}

const LOG_FILE = "/tmp/smzdm-sign.log"
