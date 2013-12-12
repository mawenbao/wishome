package database

import (
    "database/sql"
    "github.com/coopernurse/gorp"
    "github.com/robfig/revel"
    "github.com/mawenbao/wishome/app/models"
)

var (
    DbDriver string
    DbSpec string
)

type DbManager struct {
    DbMap *gorp.DbMap
}

func InitDbConfig() {
    var found bool
    if DbDriver, found = revel.Config.String("db.driver"); !found {
        revel.ERROR.Fatal("Config error: db.driver not defined")
    }
    if DbSpec, found = revel.Config.String("db.spec"); !found {
        revel.ERROR.Fatal("Config error: db.spec not defined")
    }

}

func NewDbManager() *DbManager {
    dbmgr := new(DbManager)
    if !dbmgr.Init() {
        revel.ERROR.Printf("failed to init new db manager")
        return nil
    }
    return dbmgr
}

func (dbmgr *DbManager) Init() bool {
    dbmgr.DbMap = &gorp.DbMap{Dialect: gorp.MySQLDialect{}}
    // add tables
    dbmgr.DbMap.AddTableWithName(models.User{}, "users").SetKeys(true, "id")
    // open new db connection(pooled)
    //dbmgr.DbMap.Db, err := sql.Open(DbDriver, DbSpec)
    db, err := sql.Open(DbDriver, DbSpec)
    if nil != err {
        revel.ERROR.Printf("failed to get a new database connection: %s", err)
        return false
    }
    dbmgr.DbMap.Db = db
    return true
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

func IsNameExists(dbmap *gorp.DbMap, name string) bool {
    count, err := dbmap.SelectInt("select count(*) from users where name=?", name)
    if nil != err {
        revel.ERROR.Printf("db query failed: %s", err)
        panic(err)
    }
    return 0 != count
}

func IsEmailExists(dbmap *gorp.DbMap, email string) bool {
    count, err := dbmap.SelectInt("select count(*) from users where email=?", email)
    if nil != err {
        revel.ERROR.Printf("db query failed: %s", err)
    }
    return 0 != count
}

func FindUserByName(dbmap *gorp.DbMap, name string) (u *models.User) {
    err := dbmap.SelectOne(u, "select * from users where name=?", name)
    if nil != err {
        revel.ERROR.Printf("failed to select user by name %s: %s", name, err)
    }
    return
}

func FindUserByEmail(dbmap *gorp.DbMap, email string) (u *models.User) {
    err := dbmap.SelectOne(u, "select * from users where email=?", email)
    if nil != err {
        revel.ERROR.Printf("failed to select user by email %s: %s", email, err)
    }
    return
}

func FindUserByID(dbmap *gorp.DbMap, id int32) (u *models.User) {
    err := dbmap.SelectOne(u, "select * from users where id=?", id)
    if nil != err {
        revel.ERROR.Printf("failed to select user by id %d: %s", id, err)
    }
    return
}

func SaveUser(dbmap *gorp.DbMap, u models.User) bool {
    err := dbmap.Insert(&u)
    if nil != err {
        revel.ERROR.Printf("failed to insert user %s: %s", u, err)
        return false
    }
    return true
}

func UpdateUser(dbmap *gorp.DbMap, u models.User) bool {
    count, err := dbmap.Update(&u)
    if nil != err {
        revel.ERROR.Printf("failed to update user %s: %s", u, err)
        return false
    } else if 1 != count {
        revel.ERROR.Printf("updated more than 1 user given user %s", u)
        return false
    }
    return true
}

// remove user by id
func RemoveUser(dbmap *gorp.DbMap, u models.User) bool {
    count, err := dbmap.Delete(&u)
    if nil != err {
        revel.ERROR.Printf("failed to delete user %s: %s", u, err)
        return false
    }
    if 1 != count {
        revel.ERROR.Printf("deleted more than 1 user given user %s id %d", u, u.ID)
        return false
    }
    return true
}

