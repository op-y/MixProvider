package utils

import (
    "strings"
)

type Set struct {
    S    map[string]bool
}

// Create new Set
func NewSet() *Set {
    return &Set{S: make(map[string]bool)}
}

func (s *Set) Add(element string) *Set {
    s.S[element] = true
    return s
}

func (s *Set) Exist(element string) bool {
    _, exist := s.S[element]
    return exist
}

func (s *Set) Clear() {
    s.S = make(map[string]bool)
}

func (s *Set) ToSlice() []string {
    count := len(s.S)
    if count == 0 {
        return  []string{}
    }

    result := make([]string, count)
    i := 0
    for element := range s.S {
        result[i] = element
        i++
    }
    return result
}

func (s *Set) ToString() string {
    count := len(s.S)
    if count == 0 {
        return ""
    }

    list := s.ToSlice() 

    return strings.Join(list, "|")
}

