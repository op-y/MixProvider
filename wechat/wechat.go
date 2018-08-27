package wechat

import (
    "bytes"
    "encoding/json"
    "errors"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "strings"
    "sync"
    "time"

    "mix-provider/config"
    "mix-provider/utils"
    "mix-provider/yuntongxun"

    "gopkg.in/yaml.v2"
)

type Contacts struct {
    User    []UserInfo    `yaml:"user"`
}

type UserInfo struct {
    Name       string    `yaml:"name"`
    Tel        string    `yaml:"tel"`
    Team       string    `yaml:"team"`
    ToParty    string    `yaml:"toparty"`
}

type Token struct {
    AccessToken    string    `json:"access_token"`
    ExpiresIn      int       `json:"expires_in"`
}

type Application struct {
    CorperateID    string
    Secret         string
    AgentID        string    
    AppToken       *Token
    Lock           *sync.RWMutex
}

type WechatData struct {
    ToParty        string      `json:"toparty"`
    AgentID        string      `json:"agentid"`
    MessageType    string      `json:"msgtype"`
    Text           *Warning    `json:"text"`
}

type Warning struct {
    Content    string    `json:"content"`
}

type WechatResponse struct {
    ErrCode         int       `json:"errcode"`
    ErrMsg          string    `json:"errmsg"`
    InvalidUser     string    `json:"invaliduser"`
    InvalidParty    string    `json:"invalidparty"`
    InvalidTag      string    `json:"invalidtag"`
}

var UnknownTypeError = errors.New("unknown wechat message type")
var WechatContracts *Contacts
var AlarmApplication *Application
var DeployApplication *Application


func init() {
    WechatContracts = LoadWechatContacts()
    if WechatContracts == nil {
        panic("fail to load wechat contacts mapping configuration")
    }

    AlarmApplication = &Application{
        CorperateID:config.CFG.Wechat.CorpID,
        Secret:config.CFG.Wechat.AlarmSecret, 
        AgentID:config.CFG.Wechat.AlarmAgentID,
        AppToken:nil,
        Lock:new(sync.RWMutex)}

    DeployApplication = &Application{
        CorperateID:config.CFG.Wechat.CorpID,
        Secret:config.CFG.Wechat.DeploySecret, 
        AgentID:config.CFG.Wechat.DeployAgentID,
        AppToken:nil,
        Lock:new(sync.RWMutex)}

    if err := AlarmApplication.GetWechatToken(); err != nil {
        log.Printf("fail to initialize alarm application token: %s", err.Error())
        panic("fail to initialize alarm application token")
    }

    if err := DeployApplication.GetWechatToken(); err != nil {
        log.Printf("fail to initialize deploy application token: %s", err.Error())
        panic("fail to initialize deploy application token")
    }

    log.Printf("init alarm application: %v", AlarmApplication)
    log.Printf("init deploy application: %v", DeployApplication)
}

// Load falcon/wechat users mapping from wechat.yaml
func LoadWechatContacts() *Contacts {
    contacts := new(Contacts)

    buf, err := ioutil.ReadFile("wechat.yaml")
    if err != nil {
        log.Printf("fail to read wechat.yaml: %s", err.Error())
        return nil
    }
    if err := yaml.Unmarshal(buf, contacts); err != nil {
        log.Printf("fail to unmarshal contacts: %s", err.Error())
        return nil
    }
    log.Printf("load contacts successfully: %v", contacts)

    return contacts
}

// Get wechat group id(s) which message will be send to
func GetWechatParty(tos string) string {
    s := utils.NewSet()
    toList := strings.Split(tos, ",")
    for _, to := range toList {
        for _, user := range WechatContracts.User {
            if to == user.Tel {
                s.Add(user.ToParty)
            }
        }
    }
    return s.ToString()
}

// Get wechat access token through wechat API
func (app *Application) GetWechatToken() error {
    accessTokenURL := "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid="+app.CorperateID+"&corpsecret="+app.Secret

    timeout := time.Second * 5
    client := &http.Client{Timeout: timeout}
    accessTokenResponse, err := client.Get(accessTokenURL)
    if err != nil {
        log.Printf("fail to get wechat access token: %s", err.Error())
        return err
    }
    defer accessTokenResponse.Body.Close()

    body, err := ioutil.ReadAll(accessTokenResponse.Body)
    if err != nil {
        log.Printf("fail to read response: %s", err.Error())
        return err
    }

    accessToken := &Token{}
    if err := json.Unmarshal(body, accessToken); err != nil {
        log.Printf("fail to unmarshal access token: %s", err.Error())
        return err
    }

    app.Lock.Lock()
    app.AppToken = accessToken
    app.Lock.Unlock()
    return nil
}

// Send wechat message through wechat API
func (app *Application) SendWechatMessage(toparty string, content string) (*WechatResponse, error) {
    sendMessageURL := "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token="+app.AppToken.AccessToken
    
    encodedContent := url.QueryEscape(content)
    warning := &Warning{Content: encodedContent}
    data := &WechatData{ToParty: toparty, AgentID: app.AgentID, MessageType: "text", Text: warning}

    dataJSON, err := json.Marshal(data)
    if err != nil {
        log.Printf("fail to marshal request data: %s", err.Error())
        return nil, err
    }
    decodedDataJSON, err := url.QueryUnescape(string(dataJSON))
    if err != nil {
        log.Printf("fail to unescape json data: %s", err.Error())
        return nil, err
    }

    timeout := time.Second * 5
    client := &http.Client{Timeout: timeout}
    sendMessageRequest, err := http.NewRequest("POST", sendMessageURL, bytes.NewReader([]byte(decodedDataJSON)))
    if err != nil {
        log.Printf("fail to create request: %s", err.Error())
        return nil, err
    }

    sendMessageResponse, err := client.Do(sendMessageRequest)
    if err != nil {
        log.Printf("fail to call the wechat API: %s", err.Error())
        return nil, err
    }
    defer sendMessageResponse.Body.Close()

    body, err := ioutil.ReadAll(sendMessageResponse.Body)
    if err != nil {
        log.Printf("fail to read wechat response data: %s", err.Error())
        return nil, err
    }

    wechatResponse := &WechatResponse{}
    if err := json.Unmarshal(body, wechatResponse); err != nil {
        log.Printf("fail to unmarshal wechat send message reponse: %s", err.Error())
        return nil, err
    }

    return wechatResponse, nil
}

// The entrance of wechat
func (app *Application) WechatGo(tos string, content string) error {
    toparty := config.CFG.Wechat.Toparty
    if config.CFG.Wechat.Grouping {
        party := GetWechatParty(tos)
        if "" != party {
            toparty = party
        }  
    }

    // send wechat message
    response, err := app.SendWechatMessage(toparty, content)
    if err != nil {
        log.Printf("fail to send wechat message: %s", err.Error())
        return err
    }
    log.Printf("call wechat API successfully: %d %s", response.ErrCode, response.ErrMsg)

    // retry to send wechat message if token is invalid
    switch response.ErrCode {
    case 0:
        log.Printf("send wechat message successfully: %d %s", response.ErrCode, response.ErrMsg)
        return nil
    case 40014, 42001:
        log.Printf("wechat access token expire, try to update it: %d %s", response.ErrCode, response.ErrMsg)

        if err := app.GetWechatToken(); err != nil {
            log.Printf("fail to update wechat access token: %s", err.Error())
            return err
        }

        responseRetry, err := app.SendWechatMessage(toparty, content)
        if err != nil {
            log.Printf("fail to REsend wechat message: %s", err.Error())
            return err
        }

        log.Printf("call wechat API successfully: %d %s", responseRetry.ErrCode, responseRetry.ErrMsg)
        return nil
    case 45009:
        log.Printf("API freq out of limit, try to send SMS: %d %s", response.ErrCode, response.ErrMsg)
        if config.CFG.SmsEnabled {
            warning := content + "[Out of Wechat Limit]"
            if err := yuntongxun.YuntongxunSmsGo(config.CFG.Sms.AccID, config.CFG.Sms.AppID, config.CFG.Sms.Token, config.CFG.Sms.TplID, tos, warning); err != nil {
                log.Printf("fail to wechat2sms: %s", err.Error())
                return err
            }
            log.Printf("wechat2sms successfully")
            return nil
        }
        return nil
    default:
        log.Printf("unexpect wechat API response: %d %s", response.ErrCode, response.ErrMsg)
        return nil
    }
    return nil
}

