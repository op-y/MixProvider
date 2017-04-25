package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
)

type Config struct {
    Addr             string    `json:"addr"`  
    MailEnabled      bool      `json:"mail_enabled"`  
    MailChannel      string    `json:"mail_channel"`  
    MailRelayerURL   string    `json:"mail_relayer_url"`  
    WechatEnabled    bool      `json:"wechat_enabled"`
    WechatChannel    string    `json:"wechat_channel"`
    WechatRelayerURL string    `json:"wechat_relayer_url"`  
    SMSEnabled       bool      `json:"sms_enabled"`  
    SMSChannel       string    `json:"sms_channel"`  
    SMSPriority      string    `json:"sms_priority"`  
    SMSSkip          string    `json:"sms_skip"`  
    SMSBlacklist     []string  `json:"sms_blacklist"`  
}

var config *Config

func hello(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello World!")
}

func health(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "ok")
}

func msg(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        var tos     string
        var content string

        r.ParseForm()

        if len(r.Form.Get("tos"))==0 {
            log.Printf("!!!!==msg== <result:%s>", "tos is empty")
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprintf(w, "tos is empty")
            return
        } else {
            tos = r.Form.Get("tos")
        }

        if len(r.Form.Get("content"))==0 {
            log.Printf("!!!!==msg== <result:%s>", "content is empty")
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprintf(w, "content is empty")
            return
        } else {
            content = r.Form.Get("content")
        }

        log.Printf("==msg==>>>> <tos:%s><content:%s>", tos, content)

        // Wechat
        if config.WechatEnabled {
            switch config.WechatChannel {
            case "relayer":
                result, err := wechatRelay(config.WechatRelayerURL, tos, content)
                if err != nil {
                    log.Printf("<<<<==wechat== <result:%s>", err.Error())
                    fmt.Fprintf(w, err.Error())
                } else {
                    log.Printf("<<<<==wechat== <result:%s>", result)
                    fmt.Fprintf(w, result)
                }
            case "wechat":
                result, err := wechatGo(content)
                if err != nil {
                    log.Printf("<<<<==wechat== <result:%s>", err.Error())
                    fmt.Fprintf(w, err.Error())
                } else {
                    log.Printf("<<<<==wechat== <result:%s>", result)
                    fmt.Fprintf(w, result)
                }
            default:
                log.Printf("<<<<==wechat== <result:%s>", "No such a wechat channel")
                fmt.Fprintf(w, "No such a wechat channel")
            }
        } else {
            log.Printf("<<<<==wechat== <result:%s>", "Wechat is not enabled")
            fmt.Fprintf(w, "Wechat is not enabled!")
        }

        // SMS
        if ! config.SMSEnabled {
            log.Printf("<<<<==sms== <result:%s>", "SMS is not enabled")
            fmt.Fprintf(w, "SMS is not enabled!")
            return
        }

        if ! priorityFilter(content, config.SMSPriority) {
            log.Printf("<<<<==sms== <result:%s>", "Do not pass the priority filter")
            fmt.Fprintf(w, "Do not pass the priority filter!")
            return
        }

        if ! numberFilter(tos, config.SMSSkip) {
            log.Printf("<<<<==sms== <result:%s>", "Do not pass the number filter")
            fmt.Fprintf(w, "Do not pass the number filter!")
            return
        }

        tos = blacklistFilter(tos, config.SMSBlacklist)
        if len(tos)== 0  {
            log.Printf("<<<<==sms== <result:%s>", "No number remain after passing the priority filter")
            fmt.Fprintf(w, "No number remain after passing the priority filter!")
            return
        }

        switch config.SMSChannel {
        case "yuntongxun":
            result, err := yuntongxunGo(tos, content)
            if err != nil {
                log.Printf("<<<<==sms== <result:%s>", err.Error())
                fmt.Fprintf(w, err.Error())
            } else {
                log.Printf("<<<<==sms== <result:%s>", result)
                fmt.Fprintf(w, result)
            }
        case "sendcloud":
            result, err := sendcloudGo(tos, content)
            if err != nil {
                log.Printf("<<<<==sms== <result:%s>", err.Error())
                fmt.Fprintf(w, err.Error())
            } else {
                log.Printf("<<<<==sms== <result:%s>", result)
                fmt.Fprintf(w, result)
            }
        default:
            log.Printf("<<<<==sms== <result:%s>", "No such a SMS channel")
            fmt.Fprintf(w, "No such a SMS channel!")
        }
    } else {
        log.Printf("!!!!==msg== <result:%s>", "wrong method")
        w.WriteHeader(http.StatusMethodNotAllowed)
        fmt.Fprintf(w, "wrong method")
    }
}

func mail(w http.ResponseWriter, r *http.Request) {
    if ! config.MailEnabled {
        log.Printf("<<<<==mail== <result: %s>", "Mail is not enabled")
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "Mail is not enabled!")
        return
    }

    if r.Method == "POST" {
        var tos     string
        var subject string
        var content string

        r.ParseForm()

        if len(r.Form.Get("tos"))==0 {
            log.Printf("!!!!==mail== <result:%s>", "tos is empty")
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprintf(w, "tos is empty")
            return
        } else {
            tos = r.Form.Get("tos")
        }

        if len(r.Form.Get("subject"))==0 {
            log.Printf("!!!!==mail== <result:%s>", "subject is empty")
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprintf(w, "subject is empty")
            return
        } else {
            subject = r.Form.Get("subject")
        }

        if len(r.Form.Get("content"))==0 {
            log.Printf("!!!!==mail== <result:%s>", "content is empty")
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprintf(w, "content is empty")
            return
        } else {
            content = r.Form.Get("content")
        }

        log.Printf("==mail==>>>> <tos:%s><subject:%s><content:%s>", tos, subject, content)

        // select mail channel
        switch config.MailChannel {
        case "relayer":
            result, err := mailRelay(config.MailRelayerURL, tos, subject, content)
            if err != nil {
                log.Printf("<<<<==mail== <result:%s>", err.Error())
                w.WriteHeader(http.StatusInternalServerError)
                fmt.Fprintf(w, err.Error())
            } else {
                log.Printf("<<<<==mail== <result:%s>", result)
                w.WriteHeader(http.StatusOK)
                fmt.Fprintf(w, result)
            }
        default:
            log.Printf("!!!!==mail== <result:%s>", "No such mail channel")
            w.WriteHeader(http.StatusOK)
            fmt.Fprintf(w, "No such a mail channel!")
        }
    } else {
        log.Printf("!!!!==mail== <result:%s>", "wrong method")
        w.WriteHeader(http.StatusMethodNotAllowed)
        fmt.Fprintf(w, "wrong method")
    }
}

func main() {
    // initialize logger
    log.SetFlags(log.LstdFlags | log.Lshortfile)

    // load configuration
    buf, err := ioutil.ReadFile("cfg.json")
    if err != nil {
        log.Printf("read cfg.json failed!")
        return
    }
    if err := json.Unmarshal(buf, &config); err != nil {
        log.Printf("load cfg.json error %v", err)
        return
    }

    log.Printf("=====SYSTEM STARTUP=====")

    // initialize HTTP server
    http.HandleFunc("/", hello)
    http.HandleFunc("/health", health)
    http.HandleFunc("/msg", msg)
    http.HandleFunc("/mail", mail)
    if err := http.ListenAndServe(config.Addr, nil); err != nil {
        log.Printf("HTTP server start failed!")
    }
}
