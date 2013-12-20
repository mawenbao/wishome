package controllers

import (
    "fmt"
    "sort"
    "github.com/robfig/revel"
    "github.com/mawenbao/wishome/app/models"
    "github.com/mawenbao/wishome/app/modules/caching"
)

type Admin struct {
    *revel.Controller
}

// sortFieldNum is the field index in TimerJsonResult
// sortOrder: 0 asc, 1 desc
func fetchTimerResults(sortFieldNum int, sortOrder int) []models.TimerJsonResult {
    allActionTimerResults := caching.GetAllActionTimerResults()
    if nil == allActionTimerResults {
        revel.INFO.Printf("timer cache not inited")
        return nil
    }

    timerResults := make([]models.TimerJsonResult, len(allActionTimerResults))
    for i, atr := range allActionTimerResults {
        result := models.TimerJsonResult {
            atr.Action,
            fmt.Sprintf("%.3f ms", atr.TotalTime.Seconds() * 1000 / float64(atr.HitCount)),
            atr.HitCount,
        }
        timerResults = append(timerResults[:i], result)
    }

    // sort results, except first timer result TOTAL
    var sortTgt sort.Interface = models.TimerJsonResultByAvgtime(timerResults)
    switch sortFieldNum {
    case 0:
        sortTgt = models.TimerJsonResultByAction(timerResults)
    case 2:
        sortTgt = models.TimerJsonResultByHitcount(timerResults)
    }
    // desc
    if 1 == sortOrder {
        sortTgt = sort.Reverse(sortTgt)
    }

    sort.Sort(sortTgt)
    // put TOTAL timer result in the end
    return timerResults
}

func (c Admin) Home() revel.Result {
    revel.WARN.Printf("admin signed in from %s", GetRemoteAddr(c.Controller))
    moreNavbarLinks := []models.NavbarLink{
    }
    return c.Render(moreNavbarLinks)
}

func (c Admin) GetTimerResults(sortField int, sortOrder int) revel.Result {
    results := fetchTimerResults(sortField, sortOrder)
    return c.RenderJson(results)
}

