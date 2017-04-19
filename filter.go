package main

import (
    "strings"
)

func priorityFilter(content string, priority string) bool {
    return strings.Contains(content, priority)
}

func numberFilter(tos string, skipNumber string) bool {
    if tos == skipNumber {
        return false
    } else {
        return true
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
