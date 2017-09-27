package main

import (
    "bytes"
    "crypto/md5"
    "encoding/base64"
    "encoding/hex"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "strconv"
    "strings"
    "time"
)

type YuntongxunData struct {
    To         string   `json:"to"`
    AppID      string   `json:"appId"`
    TemplateID string   `json:"templateId"`
    Datas      []string `json:"datas"`
}

func yuntongxunGo(tos string, content string) (string, error) {
    baseURL := "https://ip:port"
    accountSID := "yourAccountSID"
    applicationID := "yourApplicationID"
    templateID := 00000
    token := "yourToken"
    
    // process timestamp
    timestamp := time.Now().Format("20060102150405")

    params := accountSID + token + timestamp

    // generate signature
    md5Ctx := md5.New()
    md5Ctx.Write([]byte(params))
    signature := strings.ToUpper(hex.EncodeToString(md5Ctx.Sum(nil)))

    // generate reuqeust url
    requestURL := baseURL+"/0000-00-00/Accounts/"+accountSID+"/SMS/TemplateSMS?sig="+signature

    // prepare request body
    contents := []string{content}
    data := &YuntongxunData{To: tos, AppID: applicationID, TemplateID: strconv.Itoa(templateID), Datas: contents}

    dataJSON, err := json.Marshal(data)
    if err != nil {
        return "", err
    }

    // caculate content length
    contentLength := len(dataJSON) 
    
    // generate authorization
    authorizationString := accountSID + ":" + timestamp

    authorization := base64.StdEncoding.EncodeToString([]byte(authorizationString))

    // send SMS
    client := &http.Client{}
    request, err := http.NewRequest("POST", requestURL, bytes.NewReader(dataJSON))
    request.Header.Add("Accept", "application/json")
    request.Header.Add("Content-Type", "application/json;charset=utf-8")
    request.Header.Add("Content-Length",strconv.Itoa(contentLength))
    request.Header.Add("Authorization", authorization)
    response, err := client.Do(request)
    defer response.Body.Close()

    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return "", err
    }
    message := string(body)
    return message, nil
}
