package config

import (
    "encoding/json"
    "io/ioutil"
    "log"
)

type Config struct {
    Addr             string             `json:"addr"`  
    WechatEnabled    bool               `json:"wechat_enabled"`
    SmsEnabled       bool               `json:"sms_enabled"`  
    EmailEnabled     bool               `json:"email_enabled"`  
    Wechat           WechatConfig       `json:"wechat"`
    Sms              SmsConfig          `json:"sms"`  
    Tel              TelConfig          `json:"tel"`  
    Jianzhou         JianzhouConfig     `json:"jianzhou"`  
    Sendcloud        SendcloudConfig    `json:"sendcloud"`  
    Email            EmailConfig        `json:"email"`
    BlackList        []string           `json:"blacklist"`  
    Mask             []string           `json:"mask"`  
}

type WechatConfig struct {
    CorpID           string    `json:"corpid"`
    AlarmSecret      string    `json:"alarm_secret"`
    AlarmAgentID     string    `json:"alarm_agentid"`
    DeploySecret     string    `json:"deploy_secret"`
    DeployAgentID    string    `json:"deploy_agentid"`
    Toparty          string    `json:"toparty"`
    Grouping         bool      `json:"grouping"`
}

type SmsConfig struct {
    AccID    string    `json:"accid"`
    Token    string    `json:"token"`
    AppID    string    `json:"appid"`
    TplID    int       `json:"tplid"`
}

type TelConfig struct {
    AccID         string    `json:"accid"`
    Token         string    `json:"token"`
    AppID         string    `json:"appid"`
    DisplayNum    string    `json:"display_num"`
    MediaName     string    `json:"media_name"`
    PlayTimes     string    `json:"play_times"`
}

type JianzhouConfig struct {
    Account     string    `json:"account"`
    Password    string    `json:"password"`
}

type SendcloudConfig struct {
    User     string    `json:"user"`
    Key      string    `json:"key"`
    TplID    string    `json:"tplid"`
}

type EmailConfig struct {
    IP          string    `json:"ip"`
    Port        int       `json:"port"`
    Username    string    `json:"username"`
    Password    string    `json:"password"`
    OPs         string    `json:"ops"`
}

var CFG *Config

func init() {
    CFG = LoadConfig()
    if CFG == nil {
        panic("fail to load application configuration")
    }
}

// Load configuration from config.json
func LoadConfig() *Config {
    cfg := new(Config)

    buf, err := ioutil.ReadFile("config.json")
    if err != nil {
        log.Printf("fail to read config.json: %s", err.Error())
        return nil
    }
    if err := json.Unmarshal(buf, cfg); err != nil {
        log.Printf("fail to unmarshal config.json: %s", err.Error())
        return nil 
    }
    log.Printf("load application configuration successfully: %v", cfg)

    return cfg
}

