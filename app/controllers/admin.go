package controllers

import (
    "fmt"
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

func fetchTimerResults() []TimerResult {
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
    return timerResults
}

func (c Admin) Home() revel.Result {
    revel.WARN.Printf("admin signed in from %s", c.Request.RemoteAddr)
    moreNavbarLinks := []models.NavbarLink{
    }
    return c.Render(moreNavbarLinks)
}

func (c Admin) GetTimerResults() revel.Result {
    results := fetchTimerResults()
    return c.RenderJson(results)
}

