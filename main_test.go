package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"regexp"
	"testing"
	"time"
)

func Test_SignInConfig(t *testing.T) {
	// given
	account := conf.Account

	// when
	json, error := account.signIn()

	// then
	if error != nil {
		t.Errorf("sign-in failed: %v", error)
	}
	t.Logf("response: %v", json)
}

func Test_SendMailError(t *testing.T) {
	// given
	content := "hello, I'm sunxboy"

	// when
	sendMailError(content)

	// then
	t.Error("receive mail alert")
}

func Test_FormatMailBody(t *testing.T) {
	// given
	json := &signJson{
		ErrorCode: 0,
		ErrorMsg:  "",
		Data: signData{
			AddPoint:    100,
			CheckinNum:  "1022",
			Point:       21,
			Exp:         10000,
			Gold:        100,
			Prestige:    1000,
			Slogan:      "222",
			Rank:        333,
			Cards:       "1212",
			CanContract: 0,
		},
	}

	// when
	content, error := formatMailBody(*json)

	log.Infof("%v", string(content))
	// when
	if error != nil {
		t.Errorf("format mail body error: %v", error)
	}
}

func Test_parseSigninResponseMethod(t *testing.T) {
	// given
	response := `jQuery1124000480060641146729_1600551081132({
    "qiyu_group_info": {
        "groupid": "482066",
        "staffid": ""
    },
    "qyuid": "9a0bad010dd62f97d156fd5ba9d99ccf",
    "qytoken": "MTUvY24wa2tZZzF5ZGc9PV9fX3dVQU1lMWFxVURJK0JFazRadkUvUUE9PQ==",
    "smzdm_id": "4740115153",
    "nickname": "sunxboy",
    "avatar": "\/\/avatarimg.smzdm.com\/default\/4740115153\/5b35bd2cedc84-small.jpg",
    "point": "17238",
    "exp": "70246",
    "gold": "76",
    "silver": "212",
    "shang": {
        "day_has_shang_gold": 0,
        "day_shang_gold_limit": 200,
        "day_shang_per_gold_limit": 50
    },
    "prestige": "5",
    "level": 57,
    "vip_level": 7,
    "has_mobile": 1,
    "can_apply": 0,
    "comment_sync_sina": 0,
    "capabilities": "",
    "is_shenghuojia_vip": 0,
    "is_shenghuojia_common": 0,
    "logo": {
        "goldbl": "",
        "rank": "<div class=\"rank face-stuff-level\" title=\"57\u7ea7\"><a href=\"https:\/\/zhiyou.smzdm.com\/user\/tequan\/\" target=\"_blank\"><i class=\"icon-biaoqing-sun\"><\/i><i class=\"icon-biaoqing-sun\"><\/i><i class=\"icon-biaoqing-sun\"><\/i><i class=\"icon-biaoqing-moon\"><\/i><i class=\"icon-biaoqing-moon\"><\/i><i class=\"icon-biaoqing-star\"><\/i><\/a><\/div>",
        "vip_rank": "<div class=\"rank face-stuff-level\">\r\n                        <a href=\"https:\/\/zhiyou.smzdm.com\/user\/tequan\/\" target=\"_blank\">\r\n                        <img src=\"https:\/\/res.smzdm.com\/h5\/h5_user\/dist\/assets\/level\/7.png?v=1\">\r\n                        <\/a>\r\n                    <\/div>",
        "medal": "<div class=\"icon-medal\"><a title=\"\u7b7e\u52301666\u5929\" href=\"javascript:;\" ><img src=\"https:\/\/eimg.smzdm.com\/202004\/10\/5e9064b27b5c57553.png\" alt=\"\u7b7e\u52301666\u5929\"><\/a><\/div>",
        "zhongce_grade": ""
    },
    "logo_front": {
        "goldbl": 0,
        "medal": {
            "media": 0,
            "living": 0,
            "black5": 0,
            "juweihui": 0,
            "fans": 0,
            "seven": 0,
            "xunzhang618": 0,
            "goodboy": 0,
            "signin201609": 0
        },
        "zhongce_grade": 0,
        "official_auth": {
            "official_auth_icon": "",
            "official_auth_type": "",
            "official_auth_url": ""
        }
    },
    "checkin": {
        "has_checkin": false,
        "slogan": "<div class=\"signIn_data\">\u60a8\u5df2\u7ecf\u8fde\u7eed\u7b7e\u5230<span class=\"red\">1794<\/span>\u5929<\/div>",
        "daily_checkin_num": "1794",
        "set_checkin_url": "https:\/\/zhiyou.smzdm.com\/user\/checkin\/jsonp_checkin",
        "client_has_checkin": false,
        "weixin_has_checkin": false
    },
    "unread": {
        "notice": {
            "num": "0",
            "url": "https:\/\/zhiyou.smzdm.com\/user\/notice"
        },
        "pm": {
            "num": "0",
            "url": "https:\/\/zhiyou.smzdm.com\/user\/pm"
        },
        "comment": {
            "num": 0,
            "url": "https:\/\/zhiyou.smzdm.com\/user\/shoudaopinglun\/"
        }
    },
    "notification": {
        "comment": {
            "latest_id": 0
        },
        "notice": {
            "latest_id": "812129945"
        },
        "ab_test": 0
    },
    "is_business": false,
    "bantips": "",
    "banright": [],
    "login_error_num": 0,
    "close_comment_enter": false,
    "is_anonymous": 0,
    "bgm": [],
    "can_draw": false,
    "is_set_safepass": true,
    "has_guanzhu_dongtai": 0,
    "is_agree_protocol": 1,
    "creation_date": "2015-10-23 13:50:24",
    "avatar_ornament": [],
    "sys_date": "2020-09-20 05:31:21"
})`
	// when
	result, error := parseAsMailBody([]byte(response))

	// then
	if error != nil {
		t.Errorf("parse mail body error: %v", error)
	}

	t.Log(string(result))
}
func Test_parseSigninResponse(t *testing.T) {
	// given
	jsonData := &signJson{
		Index:   1,
		Account: conf.Account,
		Time:    time.Now(),
	}
	response := `jQuery112409568846254764496_1600556526004637000({"error_code":0,"error_msg":"","data":{"add_point":0,"checkin_num":"1795","point":17238,"exp":70246,"gold":76,"prestige":5,"rank":57,"slogan":"<div class=\"signIn_data\">\u4eca\u65e5\u5df2\u9886<span class=\"red\">0<\/span>\u79ef\u5206<\/div>","cards":"1","can_contract":0}})`

	// when
	reg := regexp.MustCompile(`^jQuery\d+_\d+\((.*?)\)$`)
	jsonStr := reg.ReplaceAll([]byte(response), []byte(`$1`))
	err := json.Unmarshal(jsonStr, jsonData)
	if err != nil {
		t.Errorf("unmarshal response error: %v", err)
	}

	if jsonData.ErrorCode != 0 {
		t.Errorf("response json data error: %v", jsonData.ErrorMsg)
	}

	if jsonData.Data.Gold != 76 {
		t.Errorf("incorrect gold: %v", jsonData.Data.Gold)
	}

	if jsonData.Data.Exp != 70246 {
		t.Errorf("incorrect experence: %v", jsonData.Data.Exp)
	}

	body, error := formatMailBody(*jsonData)

	if error != nil {
		t.Errorf("mail tpl error: %v", error)
	}

	t.Log(string(body))
}

func Test_SignInManual(t *testing.T) {
	// given
	account := &account{
		Username:  "sunxboy",
		Phone:     "15669831872",
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.102 Safari/537.36",
		Cookies:   "__ckguid=C88OxlqTFeXaboiwk88v73; device_id=213070643316005511152691610bd45641bf3a92f95a926377a7a149ca; homepage_sug=g; r_sort_type=score; sensorsdata2015jssdkcross=%7B%22distinct_id%22%3A%22174a847ddd565d-032c16b43a6595-316b7004-2073600-174a847ddd6e01%22%2C%22first_id%22%3A%22%22%2C%22props%22%3A%7B%22%24latest_traffic_source_type%22%3A%22%E7%9B%B4%E6%8E%A5%E6%B5%81%E9%87%8F%22%2C%22%24latest_search_keyword%22%3A%22%E6%9C%AA%E5%8F%96%E5%88%B0%E5%80%BC_%E7%9B%B4%E6%8E%A5%E6%89%93%E5%BC%80%22%2C%22%24latest_referrer%22%3A%22%22%7D%2C%22%24device_id%22%3A%22174a847ddd565d-032c16b43a6595-316b7004-2073600-174a847ddd6e01%22%7D; __jsluid_s=78a800afc652be1bdfdf7d68ce806374; Hm_lvt_9b7ac3d38f30fe89ff0b8a0546904e58=1600551117; zdm_qd=%7B%7D; _ga=GA1.2.80725790.1600551116; _gid=GA1.2.816973104.1600551117; sess=NTNmNGZ8MTYwNDQzOTEzN3w0NzQwMTE1MTUzfDA5MmExMzJhYTI0NTk5ZjkzMzE0YjBmNzMzM2FiODBh; user=user%3A4740115153%7C4740115153; _zdmA.uid=ZDMA.AIH-Bwnb1.1600551142.2419200; Hm_lpvt_9b7ac3d38f30fe89ff0b8a0546904e58=1600551142; smzdm_id=4740115153; __gads=ID=2c3157e32c6ca689:T=1600551117:S=ALNI_MY7TzY6SoRt7tkB3lGlIY4tJlLqwg",
	}

	// when
	result, error := account.signIn()

	// then
	if error != nil {
		t.Errorf("smzdm login error: %v", error)
		return
	}
	t.Error(string(result))
}

func Test_SigninConfig(t *testing.T) {
	// given

	// when
	smzdmSignIn()

	t.Logf("%+v", conf.Account)
}
