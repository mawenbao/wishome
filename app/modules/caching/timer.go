package caching

import (
    "github.com/robfig/revel"
    "github.com/robfig/revel/cache"
    "github.com/mawenbao/wishome/app"
    "github.com/mawenbao/wishome/app/models"
)

const (
    TOTAL_TIMER = "TOTAL"
)

func SetActionTimerResult(atr *models.ActionTimerResult) bool {
    var totalResult *models.ActionTimerResult
    atrs := GetAllActionTimerResults()

    found := false
    for i, _ := range atrs {
        if atrs[i].Action == atr.Action {
            atrs[i] = *atr
            found = true
        }
        if atrs[i].Action == TOTAL_TIMER {
            totalResult = &atrs[i]
        }
        if found {
            break
        }
    }
    if !found {
        atrs = append(atrs, *atr)
    }

    // set total result
    totalResult.HitCount += 1
    totalResult.TotalTime += atr.TotalTime

    err := cache.Set(app.CACHE_TIMER, atrs, cache.FOREVER)
    if nil != err {
        revel.ERROR.Printf("failed to set action timer in cache: %s", err)
        return false
    }
    return true
}

func GetActionTimerResult(action string) *models.ActionTimerResult {
    atrs := GetAllActionTimerResults()
    for _, atr := range atrs {
        if atr.Action == action {
            return &atr
        }
    }
    revel.INFO.Printf("failed to get action timer for action %s from cache", action)
    return nil
}

func GetAllActionTimerResults() []models.ActionTimerResult {
    atrs := make([]models.ActionTimerResult, 1)
    if err := cache.Get(app.CACHE_TIMER, &atrs); nil != err {
        revel.INFO.Println("failed to get all timer results from cache, will init a new one")
        atrs[0] = models.ActionTimerResult { Action: TOTAL_TIMER }
        cache.Set(app.CACHE_TIMER, atrs, cache.FOREVER)
    }
    return atrs
}

