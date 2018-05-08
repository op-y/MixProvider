package delayer

import (
    "fmt"
    "time"

    "mix-provider/config"
    "mix-provider/email"
    "mix-provider/utils"
)

type Collector struct {
    S    []*utils.Message
}

var C *Collector

func init() {
    C = NewCollector()
}

// Create new Collector
func NewCollector() *Collector {
    return &Collector{S: make([]*utils.Message, 0)}
}

func (c *Collector) Add(m *utils.Message) {
    c.S = append(c.S, m)
}

func (c *Collector) Clear() {
    c.S = make([]*utils.Message, 0)
}


// Send email then clear Message slice
func (c *Collector) CheckTime(emailEnabled, debug bool) {
    m := time.Now().Minute() 
    if debug || (emailEnabled && 1 == m) {
        var content string
        for _, msg := range c.S {
            s :=  fmt.Sprintf("%s %s %s <br />", msg.Endpoint, msg.Content, msg.Timestamp)
            content += s
        }
        email.EmailGo(config.CFG.Email.Username, config.CFG.Email.Password, config.CFG.Email.IP, config.CFG.Email.Port, config.CFG.Email.OPs, "Problem Summary", content)
        c.Clear()
    }
}

// Start delayer goroutine for P3 alarm request handler
func StartDelayer(emailEnabled bool, finish <-chan bool) {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()

    // waiting
DELAY:
    for {
        select {
        case <-finish:
            break DELAY
        case <-ticker.C:
            C.CheckTime(emailEnabled, false)
        default:
            time.Sleep(time.Second * 10)
        }
    }
}

