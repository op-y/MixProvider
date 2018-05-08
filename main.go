package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"

    "mix-provider/config"
    "mix-provider/delayer"
    "mix-provider/server"
)

// main
func main() {
    log.Printf("=====SYSTEM STARTUP=====")

    sysCh := make(chan os.Signal, 1)
    signal.Notify(sysCh, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
    defer close(sysCh)

    // initialize logger
    log.SetFlags(log.LstdFlags | log.Lshortfile)

    // start delayer
    chDelayer := make(chan bool, 1)
    go delayer.StartDelayer(config.CFG.EmailEnabled, chDelayer)

    // start HTTP server
    chHTTPServer := make(chan bool, 1)
    go server.StartHTTPServer(chHTTPServer)

    // waiting
MAIN:
    for {
        select {
        case <-sysCh:
            log.Printf("system signal: %v", sysCh)
            chDelayer <- true
            close(chDelayer)
            chHTTPServer <- true
            close(chHTTPServer)
            break MAIN
        }
    }

    log.Printf("=====SYSTEM STOPPED=====")
}

