package controllers

import (
    "time"
    "strings"
    "github.com/robfig/revel"
    "github.com/mawenbao/wishome/app"
    "github.com/mawenbao/wishome/app/routes"
    "github.com/mawenbao/wishome/app/models"
    "github.com/mawenbao/wishome/app/modules/caching"
)

const (
    TIMER_KEY = "timer"
)

func LoadInterceptors() {
    // only use action timer when admin timer is on
    if app.MyGlobal.AdminTimer {
        revel.WARN.Printf("action timer enabled")
        revel.InterceptFunc(startActionTimer, revel.BEFORE, revel.ALL_CONTROLLERS)
        revel.InterceptFunc(stopActionTimer, revel.FINALLY, revel.ALL_CONTROLLERS)
    }

    // check admin ip
    revel.InterceptMethod(Admin.checkAdminIP, revel.BEFORE)
}

func startActionTimer(c *revel.Controller) revel.Result {
    // init a new timer for the action
    timer := &models.ActionTimer {
        Action: c.Action,
        StartTime: time.Now(),
    }

    c.Args[TIMER_KEY] = timer
    return c.Result
}

func stopActionTimer(c *revel.Controller) revel.Result {
    currTimer := c.Args[TIMER_KEY].(*models.ActionTimer)
    if nil == currTimer {
        revel.ERROR.Print("failed to get current timer for action %s", c.Action)
        return nil
    }
    currTimer.StopTime = time.Now()
    runTime := currTimer.StopTime.Sub(currTimer.StartTime)

    // use lower case action name
    currAction := strings.ToLower(c.Action)
    timer := caching.GetActionTimerResult(currAction)
    if nil == timer {
        timer = &models.ActionTimerResult {
            RemoteAddr: c.Request.RemoteAddr,
            Action: currAction,
            TotalTime: runTime,
            HitCount: 1,
        }
    } else {
        timer.HitCount += 1
        timer.TotalTime += runTime
    }

    if !caching.SetActionTimerResult(timer) {
        revel.ERROR.Printf("failed to set action timer result in cache for %s", c.Action)
        return nil
    }

    return c.Result
}

// check if user is admin
func (c Admin) checkAdminIP() revel.Result {
    if !app.MyGlobal.IsAdminIP(c.Request.RemoteAddr) {
        c.Flash.Error(c.Message("error.require.signin"))
        revel.WARN.Printf("%s is not in the admin ip list", c.Request.RemoteAddr)
        return c.Redirect(routes.User.Signin())
    }
    return c.Result
}

