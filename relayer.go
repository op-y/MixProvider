package main

import (
    "io/ioutil"
    "net/http"
    "net/url"
)

func wechatRelay(wechatURL string, tos string, content string) (string, error) {
    data := url.Values{}
    data.Set("content", content)
    data.Set("tos", tos)
    response, err := http.PostForm(wechatURL, data)
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

func mailRelay(mailURL string, tos string, subject string, content string) (string, error) {
    data := url.Values{}
    data.Set("content", content)
    data.Set("subject", subject)
    data.Set("tos", tos)
    response, err := http.PostForm(mailURL, data)
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
