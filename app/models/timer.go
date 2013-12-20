package models

import (
    "time"
)

type ActionTimer struct {
    Action string
    StartTime time.Time
    StopTime time.Time
}

type ActionTimerResult struct {
    RemoteAddr,
    Action string
    TotalTime time.Duration
    HitCount int
}

