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
    WechatEnabled    bool      `json:"wechat_enabled"`
    WechatChannel    string    `json:"wechat_channel"`
    WechatRelayerURL string    `json:"wechat_relayer_url"`  
    SMSEnabled       bool      `json:"sms_enabled"`  
    SMSChannel       string    `json:"sms_channel"`  
    SMSBlacklist     []string  `json:"sms_blacklist"`  
}

var config *Config

func hello(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello World!")
}

func health(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "ok")
}

func wechat(w http.ResponseWriter, r *http.Request) {
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

        log.Printf("==wechat==>>>> <tos:%s><content:%s>", tos, content)

        hasWechat, hasSMS := priorityFilter(content)

        // Wechat
        if hasWechat && config.WechatEnabled {
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
        if hasSMS && config.SMSEnabled {
            tos = blacklistFilter(tos, config.SMSBlacklist)
            if len(tos)== 0  {
                log.Printf("<<<<==wechat-sms== <result:%s>", "No number remain after passing the priority filter")
                fmt.Fprintf(w, "No number remain after passing the priority filter!")
                return
            }

            switch config.SMSChannel {
            case "yuntongxun":
                result, err := yuntongxunGo(tos, content)
                if err != nil {
                    log.Printf("<<<<==wechat-sms== <result:%s>", err.Error())
                    fmt.Fprintf(w, err.Error())
                } else {
                    log.Printf("<<<<==wechat-sms== <result:%s>", result)
                    fmt.Fprintf(w, result)
                }
            case "sendcloud":
                result, err := sendcloudGo(tos, content)
                if err != nil {
                    log.Printf("<<<<==wechat-sms== <result:%s>", err.Error())
                    fmt.Fprintf(w, err.Error())
                } else {
                    log.Printf("<<<<==wechat-sms== <result:%s>", result)
                    fmt.Fprintf(w, result)
                }
            default:
                log.Printf("<<<<==wechat-sms== <result:%s>", "No such a SMS channel")
                fmt.Fprintf(w, "No such a SMS channel!")
            }
        } else {
            log.Printf("<<<<==wechat-sms== <result:%s>", "SMS is not enabled")
            fmt.Fprintf(w, "SMS is not enabled!")
        }
    } else {
        log.Printf("!!!!==wechat== <result:%s>", "wrong method")
        w.WriteHeader(http.StatusMethodNotAllowed)
        fmt.Fprintf(w, "wrong method")
    }
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

        log.Printf("==sms==>>>> <tos:%s><content:%s>", tos, content)

        // SMS
        if ! config.SMSEnabled {
            log.Printf("<<<<==sms== <result:%s>", "SMS is not enabled")
            fmt.Fprintf(w, "SMS is not enabled!")
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
    http.HandleFunc("/wechat", wechat)
    http.HandleFunc("/msg", msg)
    if err := http.ListenAndServe(config.Addr, nil); err != nil {
        log.Printf("HTTP server start failed!")
    }
}
