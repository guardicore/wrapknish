package main

import "log"

func fatal(format string, args ...interface{}) {
    log.Panicf(format, args...)
}

