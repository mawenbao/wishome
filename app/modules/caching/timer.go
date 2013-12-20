package caching

import (
    "github.com/robfig/revel"
    "github.com/robfig/revel/cache"
    "github.com/mawenbao/wishome/app"
    "github.com/mawenbao/wishome/app/models"
)

func SetActionTimerResult(atr *models.ActionTimerResult) bool {
    atrs := GetAllActionTimerResults()
    // init timer cache
    if nil == atrs {
        err := cache.Set(app.CACHE_TIMER, []models.ActionTimerResult{*atr}, cache.FOREVER)
        if nil != err {
            revel.ERROR.Printf("failed to set action timer in cache: %s", err)
            return false
        }
        return true
    }

    found := false
    for i, _ := range atrs {
        // both are lower case
        if atrs[i].Action == atr.Action {
            atrs[i] = *atr
            found = true
        }
        if found {
            break
        }
    }
    if !found {
        atrs = append(atrs, *atr)
    }

    err := cache.Set(app.CACHE_TIMER, atrs, cache.FOREVER)
    if nil != err {
        revel.ERROR.Printf("failed to set action timer in cache: %s", err)
        return false
    }
    return true
}

func GetActionTimerResult(action string) *models.ActionTimerResult {
    atrs := GetAllActionTimerResults()
    if nil == atrs {
        return nil
    }
    for _, atr := range atrs {
        if atr.Action == action {
            return &atr
        }
    }
    revel.INFO.Printf("failed to get action timer for action %s from cache", action)
    return nil
}

// the main entrance
func GetAllActionTimerResults() []models.ActionTimerResult {
    if !app.MyGlobal.AdminTimer {
        return []models.ActionTimerResult{}
    }

    atrs := make([]models.ActionTimerResult, 1)
    if err := cache.Get(app.CACHE_TIMER, &atrs); nil != err {
        revel.INFO.Println("failed to get all timer results from cache, will init a new one")
        return nil
    }
    return atrs
}

