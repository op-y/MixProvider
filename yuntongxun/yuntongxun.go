package yuntongxun

import (
    "bytes"
    "crypto/md5"
    "encoding/base64"
    "encoding/hex"
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "strconv"
    "strings"
    "time"

    "mix-provider/filter"

    "gopkg.in/yaml.v2"
)

type YuntongxunMedia struct {
    L    []MediaInfo    `yaml:"media"`
}

type MediaInfo struct {
    Key            string    `yaml:"key"`
    Filename       string    `yaml:"filename"`
    Description    string    `yaml:"description"`
}

type YuntongxunSmsData struct {
    To            string      `json:"to"`
    AppID         string      `json:"appId"`
    TemplateID    string      `json:"templateId"`
    Datas         []string    `json:"datas"`
}

type YuntongxunTelData struct {
    AppID         string    `json:"appId"`
    To            string    `json:"to"`
    MediaName     string    `json:"mediaName"`
    DisplayNum    string    `json:"displayNum"`
    PlayTimes     string    `json:"playTimes"`
}

var MediaList *YuntongxunMedia

func init() {
    MediaList = LoadYuntongxunMedia()
    if MediaList == nil {
        panic("fail to load media configuration")
    }
}

// Load falcon/yuntongxun media mapping from yuntongxun.yaml
func LoadYuntongxunMedia() *YuntongxunMedia {
    medialist := new(YuntongxunMedia)

    buf, err := ioutil.ReadFile("yuntongxun.yaml")
    if err != nil {
        log.Printf("fail to read yuntongxun.yaml: %s", err.Error())
        return nil
    }
    if err := yaml.Unmarshal(buf, medialist); err != nil {
        log.Printf("fail to unmarshal medialist: %s", err.Error())
        return nil
    }
    log.Printf("load medialist successfully: %v", medialist)

    return medialist
}

// Select media file to play
func SelectYuntongxunMedia(content string) (string, bool) {
    for _, media  := range MediaList.L {
        if strings.Contains(content, media.Key) {
            return media.Filename, true
        }
    }
    return "", false
}

// The entrance of Yuntongxun SMS
func YuntongxunSmsGo(accountSID string, applicationID string, token string, templateID int, tos string, content string) error {
    baseURL := "https://app.cloopen.com:8883"

    // filter number in blacklist
    tos = filter.BlackListFilter(tos)
    if tos == "" {
        log.Println("no receiver remain!")
        return nil
    }

    // generate signature
    timestamp := time.Now().Format("20060102150405")
    params := accountSID + token + timestamp
    md5Ctx := md5.New()
    md5Ctx.Write([]byte(params))
    signature := strings.ToUpper(hex.EncodeToString(md5Ctx.Sum(nil)))

    // generate reuqeust url
    requestURL := baseURL+"/2013-12-26/Accounts/"+accountSID+"/SMS/TemplateSMS?sig="+signature

    // prepare request body
    contents := []string{content}
    data := &YuntongxunSmsData{To: tos, AppID: applicationID, TemplateID: strconv.Itoa(templateID), Datas: contents}
    dataJSON, err := json.Marshal(data)
    if err != nil {
        log.Printf("fail to prepare request data for %s: %s", tos, err.Error())
        return err
    }

    // caculate content length
    contentLength := len(dataJSON) 
    
    // generate authorization
    authorizationString := accountSID + ":" + timestamp
    authorization := base64.StdEncoding.EncodeToString([]byte(authorizationString))

    // send
    timeout := time.Second * 5
    client := &http.Client{Timeout: timeout}
    request, err := http.NewRequest("POST", requestURL, bytes.NewReader(dataJSON))
    request.Header.Add("Accept", "application/json")
    request.Header.Add("Content-Type", "application/json;charset=utf-8")
    request.Header.Add("Content-Length",strconv.Itoa(contentLength))
    request.Header.Add("Authorization", authorization)
    response, err := client.Do(request)
    defer response.Body.Close()
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Printf("fail to read response data for %s: %s", tos, err.Error())
        return err
    }
    log.Printf("send message to %s successfully: %s", tos, string(body))
    return nil
}

// The entrance of Yuntongxun tel
func YuntongxunTelGo(accountSID string, applicationID string, token string, displayNum, defaultName, playTimes, tos string, content string) error {
    baseURL := "https://app.cloopen.com:8883"

    // filter numbers in blacklist
    tos = filter.BlackListFilter(tos)
    if tos == "" {
        log.Println("no receiver remain!")
        return nil
    }

    // select media file
    mediaName := defaultName
    name, found := SelectYuntongxunMedia(content)
    if found {
        mediaName = name
    }

    // generate signature
    timestamp := time.Now().Format("20060102150405")
    params := accountSID + token + timestamp
    md5Ctx := md5.New()
    md5Ctx.Write([]byte(params))
    signature := strings.ToUpper(hex.EncodeToString(md5Ctx.Sum(nil)))

    // generate reuqeust url
    requestURL := baseURL+"/2013-12-26/Accounts/"+accountSID+"/Calls/LandingCalls?sig="+signature

    // generate authorization
    authorizationString := accountSID + ":" + timestamp
    authorization := base64.StdEncoding.EncodeToString([]byte(authorizationString))

    // call each receiver in a loop
    toList := strings.Split(tos, ",")
    for _, to := range toList {
        // prepare request body
        data := &YuntongxunTelData{AppID: applicationID, To: to, MediaName: mediaName, DisplayNum: displayNum, PlayTimes: playTimes}
        log.Printf("data is: %v", data)
        dataJSON, err := json.Marshal(data)
        if err != nil {
            log.Printf("fail to prepare request data for %s: %s", to, err.Error())
            return err
        }

        // caculate content length
        contentLength := len(dataJSON) 
        
        // call
        timeout := time.Second * 5
        client := &http.Client{Timeout: timeout}
        request, err := http.NewRequest("POST", requestURL, bytes.NewReader(dataJSON))
        request.Header.Add("Accept", "application/json")
        request.Header.Add("Content-Type", "application/json;charset=utf-8")
        request.Header.Add("Content-Length",strconv.Itoa(contentLength))
        request.Header.Add("Authorization", authorization)
        response, err := client.Do(request)
        defer response.Body.Close()
        body, err := ioutil.ReadAll(response.Body)
        if err != nil {
            log.Printf("fail to read response data for %s: %s", to, err.Error())
            return err
        }
        log.Printf("call %s successfully: %s", to, string(body))
    }
    return nil
}

