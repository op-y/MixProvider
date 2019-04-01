package jianzhou

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "strings"
    "time"

    "mix-provider/filter"
)

// The entrance of Jianzhou SMS
func JianzhouGo(account string, password string, tos string, content string) error {
    if ok := filter.ContentFilter(content); ok {
        log.Println("message not pass the content filter!")
        return nil
    }

    tos = filter.BlackListFilter(tos)
    if tos == "" {
        log.Println("no receiver remain!")
        return nil
    }

    destmobile := filter.ReplaceCommas(tos)
    if destmobile == "" {
        log.Println("no receiver remain!")
        return nil
    }

    msgText := "【签名】" + content

    // generate reuqeust url
    requestURL := "http://www.jianzhou.sh.cn/JianzhouSMSWSServer/http/sendBatchMessage"

    // prepare request body
    data := fmt.Sprintf("account=%s&password=%s&destmobile=%s&msgText=%s", account, password, destmobile, msgText)
    log.Printf("request data is: %s", data)

    response, err := http.PostForm(requestURL, url.Values{"account": {account}, "password": {password},"msgText": {msgText},"destmobile": {destmobile}})

    // send
    timeout := time.Second * 5
    client := &http.Client{Timeout: timeout}
    request, err := http.NewRequest("POST", requestURL, strings.NewReader(data))
    if err != nil {
        log.Printf("fail to build request: %s", err.Error())
        return err
    }
    request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    response, err := client.Do(request)

    if err != nil {
        log.Printf("fail to call jianzhou API: %s", err.Error())
        return err
    }
    defer response.Body.Close()

    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Printf("fail to read response data for %s: %s", destmobile, err.Error())
        return err
    }
    log.Printf("send message to %s successfully: %s", destmobile, string(body))
    return nil
}
