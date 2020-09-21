package main

import "time"

type config struct {
	EmailFrom         string   `yaml:"emailFrom"`
	EmailFromPassword string   `yaml:"emailFromPassword"`
	EmailFromSMTP     string   `yaml:"emailFromSMTP"`
	EmailToSubject    string   `yaml:"emailToSubject"`
	EmailTo           []string `yaml:"emailTo"`
	Account           account  `yaml:"account"`
}

type account struct {
	Username  string `yaml:"username"`
	Phone     string `yaml:"phone"`
	UserAgent string `yaml:"userAgent"`
	Cookies   string `yaml:"cookies"`
}

type signJson struct {
	ErrorCode int      `json:"error_code"`
	ErrorMsg  string   `json:"error_msg"`
	Data      signData `json:"data"`
	Index     int
	Account   account
	Time      time.Time
}

type signData struct {
	AddPoint    int    `json:"add_point"`
	CheckinNum  string `json:"checkin_num"`
	Point       int    `json:"point"`
	Exp         int    `json:"exp"`
	Gold        int    `json:"gold"`
	Prestige    int    `json:"prestige"`
	Rank        int    `json:"rank"`
	Slogan      string `json:"slogan"`
	Cards       string `json:"cards"`
	CanContract int    `json:"can_contract"`
}
