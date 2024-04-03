package common

import (
    "log"
)

func OrExit[T any] (v T, e error) T {
    if e != nil {
        log.Fatalf("Fatal error: %s", e)
    }
    return v
}
