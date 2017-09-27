package main

import (
    "strings"
)

func priorityFilter(content string) (bool, bool) {
    p0 := "[P0]"
    p1 := "[P1]"
    p2 := "[P2]"

    hasWechat := strings.Contains(content, p0)
    hasSMS := strings.Contains(content, p1)
    hasBoth := strings.Contains(content, p2)

    if hasWechat {
        return true, false
    } else if hasSMS {
        return false, true
    } else if hasBoth {
        return true, true
    } else {
        return true, false
    }
}

func blacklistFilter(tos string, blacklist []string) string {
    if len(blacklist) == 0 {
        return tos
    }

    remainList := []string{}
    toList := strings.Split(tos, ",")
    for _, number := range toList {
        match := false
        for _, blackNumber := range blacklist {
            if number == blackNumber {
                match = true
                break
            }
        }
        if ! match {
            remainList = append(remainList, number)
        }
    }

    return strings.Join(remainList, ",")
}
