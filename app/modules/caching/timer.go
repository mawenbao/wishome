package caching

import (
    "sync"
    "github.com/mawenbao/wishome/app/models"
)

// in memory cache
// @TODO save results in a buffer and flush them
// to the results list every 3 or more seconds
var timerRWLock sync.RWMutex
var TimerResults = []models.ActionTimerResult{}

func SetActionTimerResult(atr *models.ActionTimerResult) {
    timerRWLock.Lock()
    defer timerRWLock.Unlock()

    found := false
    for i, _ := range TimerResults {
        // both are lower case
        if TimerResults[i].Action == atr.Action {
            TimerResults[i] = *atr
            found = true
            break
        }
    }
    if !found {
        TimerResults = append(TimerResults, *atr)
    }
}

func GetActionTimerResult(action string) *models.ActionTimerResult {
    timerRWLock.RLock()
    defer timerRWLock.RUnlock()

    for _, atr := range TimerResults {
        if atr.Action == action {
            return &atr
        }
    }
    return nil
}

func GetAllActionTimerResults() []models.ActionTimerResult {
    timerRWLock.RLock()
    defer timerRWLock.RUnlock()

    tgtResults := make([]models.ActionTimerResult, len(TimerResults))
    copy(tgtResults, TimerResults)
    return tgtResults
}

