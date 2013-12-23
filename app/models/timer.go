package models

import (
    "time"
    "strconv"
    "strings"
    "github.com/robfig/revel"
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

// sort helper
func compareAvgtime(a, b string) bool {
    ai, err := strconv.ParseFloat(strings.Split(a, " ")[0], 32)
    if nil != err {
        revel.ERROR.Printf("failed to parse timer avgtime %s", a)
        return true
    }
    bi, err := strconv.ParseFloat(strings.Split(b, " ")[0], 32)
    if nil != err {
        revel.ERROR.Printf("failed to parse timer avgtime %s", a)
        return true
    }
    return ai <= bi
}

// sort desc
type TimerJsonResultByAction []TimerJsonResult
func (s TimerJsonResultByAction) Len() int { return len(s) }
func (s TimerJsonResultByAction) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s TimerJsonResultByAction) Less(i, j int) bool { return s[i].Action <= s[j].Action }

type TimerJsonResultByAvgtime []TimerJsonResult
func (s TimerJsonResultByAvgtime) Len() int { return len(s) }
func (s TimerJsonResultByAvgtime) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s TimerJsonResultByAvgtime) Less(i, j int) bool { return compareAvgtime(s[i].AverageTime, s[j].AverageTime) }

type TimerJsonResultByHitcount []TimerJsonResult
func (s TimerJsonResultByHitcount) Len() int { return len(s) }
func (s TimerJsonResultByHitcount) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s TimerJsonResultByHitcount) Less(i, j int) bool { return s[i].HitCount <= s[j].HitCount }

