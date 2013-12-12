package modules

import (
    "github.com/robfig/revel"
    "github.com/mawenbao/wishome/app/modules/database"
)

func init() {
    revel.OnAppStart(database.InitDbConfig)
}

