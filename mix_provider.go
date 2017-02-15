package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "os"
)

type Config struct {
    Addr          string    `json:"addr"`  
    Log           string    `json:"log"`  
    WeixinEnabled bool      `json:"weixin_enabled"`
    WeixinUrl     string    `json:"weixin_url"`  
    SmsEnabled    bool      `json:"sms_enabled"`  
    SmsUrl        string    `json:"sms_url"`  
    MailEnabled   bool      `json:"mail_enabled"`  
    MailUrl       string    `json:"mail_url"`  
}

var (
    config *Config
    logger *(log.Logger)
)

func sendWeixin(weixinUrl string, content string, tos string) (string, error) {
    data := url.Values{}
    data.Set("content", content)
    data.Set("tos", tos)
    resp, err := http.PostForm(weixinUrl, data)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    result := string(body)
    return result, nil
}

func sendSMS(smsUrl string, content string, tos string) (string, error) {
    data := url.Values{}
    data.Set("source", "1")
    data.Set("to", tos)
    data.Set("templateId", "154110")
    data.Set("datas", content)
    resp, err := http.PostForm(smsUrl, data)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    result := string(body)
    return result, nil
}

func sendMail(mailUrl string, content string, subject string, tos string) (string, error) {
    data := url.Values{}
    data.Set("content", content)
    data.Set("subject", subject)
    data.Set("tos", tos)
    resp, err := http.PostForm(mailUrl, data)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    result := string(body)
    return result, nil
}

func hello(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello World!")
}

func health(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "ok")
}

func msg(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        var content string
        var tos     string

        r.ParseForm()

        if len(r.Form.Get("content"))==0 {
            logger.Printf("!!!!==msg== <result:%s>", "content is empty")
            fmt.Fprintf(w, "content is empty")
            return
        } else {
            content = r.Form.Get("content")
        }

        if len(r.Form.Get("tos"))==0 {
            logger.Printf("!!!!==msg== <result:%s>", "tos is empty")
            fmt.Fprintf(w, "tos is empty")
            return
        } else {
            tos = r.Form.Get("tos")
        }

        logger.Printf("==msg==>>>> <content:%s><tos:%s>", content, tos)
        if config.WeixinEnabled {
            result, err := sendWeixin(config.WeixinUrl, content, tos)
            if err != nil {
                logger.Printf("<<<<==weixin== <result:%s>", err.Error())
                fmt.Fprintf(w, err.Error())
            } else {
                logger.Printf("<<<<==weixin== <result:%s>", result)
                fmt.Fprintf(w, result)
            }
        }
        if config.SmsEnabled {
            result, err := sendSMS(config.SmsUrl, content, tos)
            if err != nil {
                logger.Printf("<<<<==sms== <result:%s>", err.Error())
                fmt.Fprintf(w, err.Error())
            } else {
                logger.Printf("<<<<==sms== <result:%s>", result)
                fmt.Fprintf(w, result)
            }
        }
    } else {
        logger.Printf("!!!!==msg== <result:%s>", "wrong method")
        fmt.Fprintf(w, "wrong method")
    }
}

func mail(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        var content string
        var subject string
        var tos     string

        r.ParseForm()

        if len(r.Form.Get("content"))==0 {
            logger.Printf("!!!!==mail== <result:%s>", "content is empty")
            fmt.Fprintf(w, "content is empty")
            return
        } else {
            content = r.Form.Get("content")
        }

        if len(r.Form.Get("subject"))==0 {
            logger.Printf("!!!!==mail== <result:%s>", "subject is empty")
            fmt.Fprintf(w, "subject is empty")
            return
        } else {
            subject = r.Form.Get("subject")
        }

        if len(r.Form.Get("tos"))==0 {
            logger.Printf("!!!!==mail== <result:%s>", "tos is empty")
            fmt.Fprintf(w, "tos is empty")
            return
        } else {
            tos = r.Form.Get("tos")
        }

        logger.Printf("==mail==>>>> <content:%s><subject:%s><tos:%s>", content, subject, tos)
        if config.MailEnabled {
            result, err := sendMail(config.MailUrl, content, subject, tos)
            if err != nil {
                logger.Printf("<<<<==mail== <result:%s>", err.Error())
                fmt.Fprintf(w, err.Error())
            } else {
                logger.Printf("<<<<==mail== <result:%s>", result)
                fmt.Fprintf(w, result)
            }
        }
    } else {
        logger.Printf("!!!!==mail== <result:%s>", "wrong method")
        fmt.Fprintf(w, "wrong method")
    }
}

func main() {
    // Read configuration file cfg.json
    buf, err := ioutil.ReadFile("cfg.json")
    if err != nil {
        fmt.Printf("read cfg.json failed!")
        return
    }

    // Load configuration from cfg.json
    if err := json.Unmarshal(buf, &config); err != nil {
        fmt.Printf("load cfg.json error %v", err)
        return
    }

    // initialize logger
    file, err := os.OpenFile(config.Log, os.O_CREATE|os.O_APPEND|os.O_RDWR,0644)
    if err != nil {
        fmt.Println("open log file failed!")
        return
    }
    defer file.Close()
    logger = log.New(file, "", log.LstdFlags)
    logger.Println("=====SYSTEM STARTUP=====")

    // initialize HTTP server
    http.HandleFunc("/", hello)
    http.HandleFunc("/health", health)
    http.HandleFunc("/msg", msg)
    http.HandleFunc("/mail", mail)
    if err := http.ListenAndServe(config.Addr, nil); err != nil {
        fmt.Println("HTTP server start failed!")
    }
}
