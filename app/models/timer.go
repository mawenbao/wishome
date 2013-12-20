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

type TimerJsonResult struct {
    Action string `json:"action"`
    AverageTime string `json:"avgtime"` // format %.3f ms
    HitCount int `json:"hit"`
}

// sort desc
type TimerJsonResultByAction []TimerJsonResult
func (s TimerJsonResultByAction) Len() int { return len(s) }
func (s TimerJsonResultByAction) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s TimerJsonResultByAction) Less(i, j int) bool { return s[i].Action <= s[j].Action }

type TimerJsonResultByAvgtime []TimerJsonResult
func (s TimerJsonResultByAvgtime) Len() int { return len(s) }
func (s TimerJsonResultByAvgtime) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s TimerJsonResultByAvgtime) Less(i, j int) bool { return s[i].AverageTime <= s[j].AverageTime }

type TimerJsonResultByHitcount []TimerJsonResult
func (s TimerJsonResultByHitcount) Len() int { return len(s) }
func (s TimerJsonResultByHitcount) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s TimerJsonResultByHitcount) Less(i, j int) bool { return s[i].HitCount <= s[j].HitCount }

