package utils

import (
    "strings"
)

type Message struct {
    Priority     string
    Status       string
    Endpoint     string
    Content      string
    Timestamp    string
}

// Build a new Message
func BuildMessage(content string) *Message {
    msg := new(Message)
    fields := strings.Split(content, "][")
    for idx, field := range fields {
        s := strings.Trim(field, "[]")
        switch idx {
        case 0:
            msg.Priority = s
        case 1:
            msg.Status = s
        case 2:
            msg.Endpoint = s
        case 3:
        case 4:
            msg.Content = s
        case 5:
            l := strings.Split(s, " ")
            t := l[1:]
            ts := strings.Join(t, "T")
            msg.Timestamp = ts
        default:
        }
    }
    return msg
}

