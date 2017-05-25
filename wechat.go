package main

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "net/url"
)

type Token struct {
    AccessToken string    `json:"access_token"`
    ExpiresIn   int       `json:"expires_in"`
}

type Warning struct {
    Content     string    `json:"content"`
}

type WechatData struct {
    ToParty     string    `json:"toparty"`
    AgentID     string    `json:"agentid"`
    MessageType string    `json:"msgtype"`
    Text        *Warning  `json:"text"`
}

func wechatGo(content string) (string, error) {
    corperateID := "xxx"
    secret := "iamasecret"

    // get access token
    accessTokenURL := "https://127.0.0.1:80/xxx/gettoken?corpid="+corperateID+"&corpsecret="+secret

    accessTokenResponse, err := http.Get(accessTokenURL)
    if err != nil {
        return "", err
    }
    defer accessTokenResponse.Body.Close()
    body, err := ioutil.ReadAll(accessTokenResponse.Body)
    if err != nil {
        return "", err
    }
    accessToken := &Token{}    
    if err := json.Unmarshal(body, accessToken); err != nil {
        return "", err
    }

    // send wechat message
    sendMessageURL := "https://127.0.0.1:80/message/send"
    
    encodedContent := url.QueryEscape(content)
    warning := &Warning{Content: encodedContent}
    data := &WechatData{ToParty: "3", AgentID: "2", MessageType: "text", Text: warning}
    dataJSON, err := json.Marshal(data)
    if err != nil {
        return "", err
    }
    decodedDataJSON, err := url.QueryUnescape(string(dataJSON))
    if err != nil {
        return "", err
    }

    client := &http.Client{}
    sendMessageRequest, err := http.NewRequest("POST", sendMessageURL, bytes.NewReader([]byte(decodedDataJSON)))
    sendMessageResponse, err := client.Do(sendMessageRequest)
    defer sendMessageResponse.Body.Close()

    body, err = ioutil.ReadAll(sendMessageResponse.Body)
    if err != nil {
        return "", err
    }
    message := string(body)
    return message, nil
}
