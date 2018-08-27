package server

import (
    "fmt"
    "log"
    "net/http"

    "mix-provider/config"
    "mix-provider/delayer"
    "mix-provider/email"
    "mix-provider/filter"
    "mix-provider/utils"
    "mix-provider/wechat"
    "mix-provider/yuntongxun"
)

// Handler("/")
func Hello(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Hello, MixProvider!")
}

// Handler("/health")
func Health(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "ok")
}

// Handler("/email")
func Email(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        var tos     string
        var subject string
        var content string

        r.ParseForm()

        if len(r.Form.Get("tos"))==0 {
            log.Printf("!!!!==email== <result:%s>", "tos is empty")
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprintf(w, "tos is empty")
            return
        } else {
            tos = r.Form.Get("tos")
        }

        if len(r.Form.Get("subject"))==0 {
            log.Printf("!!!!==email== <result:%s>", "subject is empty")
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprintf(w, "subject is empty")
            return
        } else {
            subject = r.Form.Get("subject")
        }

        if len(r.Form.Get("content"))==0 {
            log.Printf("!!!!==email== <result:%s>", "content is empty")
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprintf(w, "content is empty")
            return
        } else {
            content = r.Form.Get("content")
        }

        log.Printf("==email==>>>> <tos:%s><subject:%s><content:%s>", tos, subject, content)

        // send email
        if ! config.CFG.EmailEnabled {
            log.Printf("<<<<==email== <result:%s>", "email is not enabled")
            w.WriteHeader(http.StatusOK)
            fmt.Fprintf(w, "email is not enabled!")
            return
        }

        if err := email.EmailGo(config.CFG.Email.Username, config.CFG.Email.Password, config.CFG.Email.IP, config.CFG.Email.Port, tos, subject, content); err != nil {
            log.Printf("<<<<==email== <result:%s>", err.Error())
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, err.Error())
        }
        log.Println("<<<<==email== ok!")
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "done")
    } else {
        log.Printf("!!!!==email== <result:%s>", "wrong method")
        w.WriteHeader(http.StatusMethodNotAllowed)
        fmt.Fprintf(w, "wrong method")
    }
}

// Handler("/msg")
func Message(w http.ResponseWriter, r *http.Request) {
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

        // send sms
        if ! config.CFG.SmsEnabled {
            log.Printf("<<<<==sms== <result:%s>", "sms is not enabled")
            w.WriteHeader(http.StatusOK)
            fmt.Fprintf(w, "sms is not enabled!")
            return
        }

        if err := yuntongxun.YuntongxunSmsGo(config.CFG.Sms.AccID, config.CFG.Sms.AppID, config.CFG.Sms.Token, config.CFG.Sms.TplID, tos, content); err != nil {
            log.Printf("<<<<==sms== <result:%s>", err.Error())
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, err.Error())
            return
        }
        log.Println("<<<<==sms== ok!")
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "done")
    } else {
        log.Printf("!!!!==msg== <result:%s>", "wrong method")
        w.WriteHeader(http.StatusMethodNotAllowed)
        fmt.Fprintf(w, "wrong method")
        return
    }
}

// Handler("/deploy")
func Deploy(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        var tos     string
        var content string

        r.ParseForm()

        if len(r.Form.Get("tos"))==0 {
            log.Printf("!!!!==deploy== <result:%s>", "tos is empty")
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprintf(w, "tos is empty")
            return
        } else {
            tos = r.Form.Get("tos")
        }

        if len(r.Form.Get("content"))==0 {
            log.Printf("!!!!==deploy== <result:%s>", "content is empty")
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprintf(w, "content is empty")
            return
        } else {
            content = r.Form.Get("content")
        }

        log.Printf("==deploy==>>>> <tos:%s><content:%s>", tos, content)

        // send wechat
        if ! config.CFG.WechatEnabled {
            log.Printf("<<<<==wechat== <result:%s>", "wechat is not enabled")
            w.WriteHeader(http.StatusOK)
            fmt.Fprintf(w, "wechat is not enabled!")
            return
        }

        if err := wechat.DeployApplication.WechatGo(tos, content); err != nil {
            log.Printf("<<<<==wechat== <result:%s>", err.Error())
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, err.Error())
            return
        }
        log.Printf("<<<<==wechat== <result:%s>", "done")
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "done")
        return
    } else {
        log.Printf("!!!!==deploy== <result:%s>", "wrong method")
        w.WriteHeader(http.StatusMethodNotAllowed)
        fmt.Fprintf(w, "wrong method")
        return
    }
}

// Handler("/alarm")
func Alarm(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        var tos     string
        var content string

        r.ParseForm()

        if len(r.Form.Get("tos"))==0 {
            log.Printf("!!!!==alarm== <result:%s>", "tos is empty")
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprintf(w, "tos is empty")
            return
        } else {
            tos = r.Form.Get("tos")
        }

        if len(r.Form.Get("content"))==0 {
            log.Printf("!!!!==alarm== <result:%s>", "content is empty")
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprintf(w, "content is empty")
            return
        } else {
            content = r.Form.Get("content")
        }

        log.Printf("==alarm==>>>> <tos:%s><content:%s>", tos, content)

        priority := filter.PriorityFilter(content)
        switch priority {
        case 0:
            if err := P0(tos, content); err != nil {
                log.Printf("<<<<==P0== <result:%s>", err.Error())
                w.WriteHeader(http.StatusInternalServerError)
                fmt.Fprintf(w, err.Error())
            }
            log.Printf("<<<<==P0== <result:%s>", "P0 done")
            w.WriteHeader(http.StatusOK)
            fmt.Fprintf(w, "done")
            return
        case 1:
            if err := P1(tos, content); err != nil {
                log.Printf("<<<<==P1== <result:%s>", err.Error())
                w.WriteHeader(http.StatusInternalServerError)
                fmt.Fprintf(w, err.Error())
            }
            log.Printf("<<<<==P1== <result:%s>", "P1 done")
            w.WriteHeader(http.StatusOK)
            fmt.Fprintf(w, "done")
            return
        case 2:
            if err := P2(tos, content); err != nil {
                log.Printf("<<<<==P2== <result:%s>", err.Error())
                w.WriteHeader(http.StatusInternalServerError)
                fmt.Fprintf(w, err.Error())
            }
            log.Printf("<<<<==P2== <result:%s>", "P2 done")
            w.WriteHeader(http.StatusOK)
            fmt.Fprintf(w, "done")
            return
        case 3:
            if err := P3(tos, content); err != nil {
                log.Printf("<<<<==P3== <result:%s>", err.Error())
                w.WriteHeader(http.StatusInternalServerError)
                fmt.Fprintf(w, err.Error())
            }
            log.Printf("<<<<==P3== <result:%s>", "P3 done")
            w.WriteHeader(http.StatusOK)
            fmt.Fprintf(w, "done")
            return
        case 4:
            if err := P4(tos, content); err != nil {
                log.Printf("<<<<==P4== <result:%s>", err.Error())
                w.WriteHeader(http.StatusInternalServerError)
                fmt.Fprintf(w, err.Error())
            }
            log.Printf("<<<<==P4== <result:%s>", "P4 done")
            w.WriteHeader(http.StatusOK)
            fmt.Fprintf(w, "done")
            return
        case 5:
            if err := P5(tos, content); err != nil {
                log.Printf("<<<<==P5== <result:%s>", err.Error())
                w.WriteHeader(http.StatusInternalServerError)
                fmt.Fprintf(w, err.Error())
            }
            log.Printf("<<<<==P5== <result:%s>", "P5 done")
            w.WriteHeader(http.StatusOK)
            fmt.Fprintf(w, "done")
            return
        case 6:
            if err := P6(tos, content); err != nil {
                log.Printf("<<<<==P6== <result:%s>", err.Error())
                w.WriteHeader(http.StatusInternalServerError)
                fmt.Fprintf(w, err.Error())
            }
            log.Printf("<<<<==P6== <result:%s>", "P6 done")
            w.WriteHeader(http.StatusOK)
            fmt.Fprintf(w, "done")
            return
        default:
            log.Printf("<<<<==alarm== <result:%s>", "unknown priority")
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprintf(w, "unknown priority")
            return
        }
    } else {
        log.Printf("!!!!==wechat== <result:%s>", "wrong method")
        w.WriteHeader(http.StatusMethodNotAllowed)
        fmt.Fprintf(w, "wrong method")
        return
    }
}

// Process P0 alarm message
func P0(tos, content string) error {
    var telErr error
    var smsErr error

    msg := utils.BuildMessage(content)
    if "PROBLEM" ==  msg.Status {
        telErr = yuntongxun.YuntongxunTelGo(config.CFG.Tel.AccID, config.CFG.Tel.AppID, config.CFG.Tel.Token, config.CFG.Tel.DisplayNum, config.CFG.Tel.MediaName, config.CFG.Tel.PlayTimes, tos, content)
        if telErr != nil {
            log.Printf("<<<<==tel== <result:%s>", telErr.Error())
        }
        log.Printf("<<<<==tel== <result:%s>", "ok")
    }

    if config.CFG.SmsEnabled {
        smsErr = yuntongxun.YuntongxunSmsGo(config.CFG.Sms.AccID, config.CFG.Sms.AppID, config.CFG.Sms.Token, config.CFG.Sms.TplID, tos, content)
        if smsErr != nil {
            log.Printf("<<<<==sms== <result:%s>", smsErr.Error())
        }
        log.Printf("<<<<==sms== <result:%s>", "ok")
    }
    log.Printf("<<<<==sms== <result:%s>", "sms is not enabled")

    if telErr != nil {
        return telErr
    }
    return nil
}

// Process P1 alarm message
func P1(tos, content string) error {
    if config.CFG.WechatEnabled {
        if err := wechat.AlarmApplication.WechatGo(tos, content); err != nil {
            log.Printf("<<<<==wechat== <result:%s>", err.Error())
            return err
        }
        log.Printf("<<<<==wechat== <result:%s>", "ok")
        return nil
    }
    log.Printf("<<<<==wechat== <result:%s>", "wechat is not enabled")
    return nil
}

// Process P2 alarm message
func P2(tos, content string) error {
    if config.CFG.SmsEnabled {
        if err := yuntongxun.YuntongxunSmsGo(config.CFG.Sms.AccID, config.CFG.Sms.AppID, config.CFG.Sms.Token, config.CFG.Sms.TplID, tos, content); err != nil {
            log.Printf("<<<<==sms== <result:%s>", err.Error())
            return err
        }
        log.Printf("<<<<==sms== <result:%s>", "ok")
        return nil
    }
    log.Printf("<<<<==sms== <result:%s>", "sms is not enabled")
    return nil
}

// Process P3 alarm message
func P3(tos, content string) error {
    msg := utils.BuildMessage(content)
    if "PROBLEM" ==  msg.Status {
        delayer.C.Add(msg)
        log.Printf("<<<<==delayer== <result:%s>", "ok")
        return nil
    }
    log.Printf("<<<<==delayer== <result:%s>", "ignore")
    return nil
}

// Process P4 alarm message
func P4(tos, content string) error {
    return nil
}

// Process P5 alarm message
func P5(tos, content string) error {
    return nil
}

// Process P6 alarm message
func P6(tos ,content string) error {
    return nil
}

// Start HTTP server goroutine
func StartHTTPServer(finish <-chan bool) {
    http.HandleFunc("/", Hello)
    http.HandleFunc("/health", Health)
    http.HandleFunc("/email", Email)
    http.HandleFunc("/msg", Message)
    http.HandleFunc("/deploy", Deploy)
    http.HandleFunc("/alarm", Alarm)

    server := &http.Server{Addr: config.CFG.Addr}
    
    // start server in new goroutine
    go func() {
        if err := server.ListenAndServe(); err != nil {
            log.Printf("HTTP server error: %s", err.Error())
        }
    } ()

    // waiting
SERVER:
    for {
        select {
        case <-finish:
            if err := server.Shutdown(nil); err != nil {
                panic(err)
            }
            break SERVER
        }
    }

    log.Println("HTTP server stopped")
}

