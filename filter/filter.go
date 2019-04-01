package filter

import (
    "strings"

    "mix-provider/config"
)

// Judge the priority by the message content
func PriorityFilter(content string) int {
    // P0: realtime tel + sms
    if strings.HasPrefix(content, "[P0]") {
        return 0;
    }
    // P1: realtime wechat
    if strings.HasPrefix(content, "[P1]") {
        return 1;
    }
    // P2: realtime sms
    if strings.HasPrefix(content, "[P2]") {
        return 2;
    }
    // P3: summary wechat and email per hour
    if strings.HasPrefix(content, "[P3]") {
        return 3;
    }
    // P4: retain
    if strings.HasPrefix(content, "[P4]") {
        return 4;
    }
    // P5: retain
    if strings.HasPrefix(content, "[P5]") {
        return 5;
    }
    // P6: retain
    if strings.HasPrefix(content, "[P6]") {
        return 6;
    }
    // default: wechat
    return 1;
}

// Filter the numbers in the blacklist
func BlackListFilter(tos string) string {
    if len(config.CFG.BlackList) == 0 {
        return tos
    }

    remainList := []string{}
    toList := strings.Split(tos, ",")
    for _, number := range toList {
        match := false
        for _, blackNumber := range config.CFG.BlackList {
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

// Replace commas by semicolons
func ReplaceCommas(tos string) string {
    return strings.Replace(tos, ",", ";", -1)
}

// Filter the content in the mask keyword list
func ContentFilter(content string) bool {
    if strings.Contains(content, "[OK]") {
        for _, keyword := range config.CFG.Mask {
            if strings.Contains(content, keyword) {
                return true
            }
        }
    }
    return false
}
