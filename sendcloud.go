package main

import (
    "crypto/md5"
    "encoding/hex"
    "io/ioutil"
    "net/http"
    "net/url"
    "strings"
)

func sendcloudGo(tos string, content string) (string, error) {
    smsUrl := "http://host:port/api/path"
    smsKey := "yourKey"
    smsUser := "yourUsername"
    smsTemplateId := "yourTemplateId"

    // replace com with c0m
    message := strings.Replace(content, ".com", ".c0m", -1)

    vars := "{\"content\":\"" + message + "\"}"

    params := smsKey + "&" + "phone=" + tos + "&" + "smsUser=" + smsUser + "&" + "templateId=" + smsTemplateId + "&" + "vars=" + vars + "&" + smsKey

    md5Ctx := md5.New()
    md5Ctx.Write([]byte(params))
    signature := hex.EncodeToString(md5Ctx.Sum(nil))

    data := url.Values{}
    data.Set("smsUser", smsUser)
    data.Set("templateId", smsTemplateId)
    data.Set("phone", tos)
    data.Set("vars", vars)
    data.Set("signature", signature)

    response, err := http.PostForm(smsUrl, data)
    if err != nil {
        return "", err
    }
    defer response.Body.Close()
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return "", err
    }
    result := string(body)
    return result, nil
}
