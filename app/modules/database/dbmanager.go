package database

import (
    "database/sql"
    "github.com/coopernurse/gorp"
    "github.com/robfig/revel"
    "github.com/mawenbao/wishome/app"
    "github.com/mawenbao/wishome/app/models"
)

var (
    MyDbManager = DbManager{} // global db handle
)

func init() {
    // init global db handle
    revel.OnAppStart(func() {
        if !MyDbManager.Init() {
            revel.ERROR.Panicf("failed to init global db manager")
        }
    })
}

type DbManager struct {
    DbMap *gorp.DbMap
}

func (dbmgr *DbManager) Init() bool {
    dbmgr.DbMap = &gorp.DbMap{Dialect: gorp.MySQLDialect{}}
    // add tables
    dbmgr.DbMap.AddTableWithName(models.User{}, app.TABLE_USERS).SetKeys(true, "id")
    // open db(pooled)
    db, err := sql.Open(app.MyGlobal.DbDriver, app.MyGlobal.DbSpec)
    if nil != err {
        revel.ERROR.Printf("failed to get a new database connection: %s", err)
        return false
    }
    dbmgr.DbMap.Db = db
    return true
}

func (dbmgr *DbManager) Db() *gorp.DbMap {
    if nil == dbmgr || nil == dbmgr.DbMap {
        revel.ERROR.Panicf("db manger has not been initialized")
    }
    return dbmgr.DbMap
}

func (dbmgr *DbManager) Close() bool {
    if nil == dbmgr.DbMap || nil == dbmgr.DbMap.Db {
        revel.ERROR.Println("DbMap or DbMap.Db is nil")
        return false
    }
    err := dbmgr.DbMap.Db.Close()
    if nil != err {
        revel.ERROR.Printf("failed to close db connection: %s", err)
        return false
    }
    return true
}

