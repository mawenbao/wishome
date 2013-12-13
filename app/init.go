package app

import (
    "log"
    "time"
    "github.com/robfig/revel"
)

// store global custom configuration in app.conf
var MyGlobal map[string]interface{} = make(map[string]interface{})

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.ActionInvoker,           // Invoke the action.
	}

    // set log flags
    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

    revel.OnAppStart(parseCustomConfig)
}

// parse custom settings in app.conf and save in MyGlobal
// parse failure will result in panic
func parseCustomConfig() {
    keyLen, found := revel.Config.Int(CONFIG_RESETPASS_KEY_LEN)
    if !found {
        keyLen = 32
    } else {
        MyGlobal[CONFIG_RESETPASS_KEY_LEN] = keyLen
    }

    keyLife, found := revel.Config.String(CONFIG_RESETPASS_KEY_LIFE)
    if !found {
        MyGlobal[CONFIG_RESETPASS_KEY_LIFE] = time.Duration(30 * time.Minute)
    } else {
        keyLifeTime, err := time.ParseDuration(keyLife)
        if nil != err {
            revel.ERROR.Printf("failed to parse %s in app.conf", CONFIG_RESETPASS_KEY_LIFE)
            panic("parse app.conf custom config failed")
        }
        MyGlobal[CONFIG_RESETPASS_KEY_LIFE] = keyLifeTime
    }
}

