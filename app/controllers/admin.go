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

type TimerResult struct {
    Action string `json:"action"`
    AverageTime string `json:"avgtime"` // format %.3f ms
    HitCount int `json:"hit"`
}

type TimerResultByAction []TimerResult
func (s TimerResultByAction) Len() int { return len(s) }
func (s TimerResultByAction) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s TimerResultByAction) Less(i, j int) bool { return s[i].Action < s[j].Action }

type TimerResultByAvgtime []TimerResult
func (s TimerResultByAvgtime) Len() int { return len(s) }
func (s TimerResultByAvgtime) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s TimerResultByAvgtime) Less(i, j int) bool { return s[i].AverageTime > s[j].AverageTime }

type TimerResultByHitcount []TimerResult
func (s TimerResultByHitcount) Len() int { return len(s) }
func (s TimerResultByHitcount) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s TimerResultByHitcount) Less(i, j int) bool { return s[i].HitCount > s[j].HitCount }

func fetchTimerResults(sortFieldNum int) []TimerResult {
    allActionTimerResults := caching.GetAllActionTimerResults()
    timerResults := make([]TimerResult, len(allActionTimerResults))
    for i, atr := range allActionTimerResults {
        timerResults = append(
            timerResults[:i],
            TimerResult{
                atr.Action,
                fmt.Sprintf("%.3f ms", atr.TotalTime.Seconds() * 1000 / float64(atr.HitCount)),
                atr.HitCount,
            },
        )
    }

    // sort results, except first timer result TOTAL
    switch sortFieldNum {
    case 0:
        sort.Sort(TimerResultByAction(timerResults)[1:])
    case 2:
        sort.Sort(TimerResultByHitcount(timerResults)[1:])
    default:
        sort.Sort(TimerResultByAvgtime(timerResults)[1:])
    }

    return timerResults
}

func (c Admin) Home() revel.Result {
    revel.WARN.Printf("admin signed in from %s", c.Request.RemoteAddr)
    moreNavbarLinks := []models.NavbarLink{
    }
    return c.Render(moreNavbarLinks)
}

func (c Admin) GetTimerResults(sort int) revel.Result {
    results := fetchTimerResults(sort)
    return c.RenderJson(results)
}

