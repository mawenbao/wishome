package database

import (
    "fmt"
    "database/sql"
    "github.com/coopernurse/gorp"
    "github.com/robfig/revel"
    "github.com/mawenbao/wishome/app"
    "github.com/mawenbao/wishome/app/models"
    _ "github.com/go-sql-driver/mysql"
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
    dbmgr.DbMap.AddTableWithName(models.User{}, app.TABLE_USERS).SetKeys(true, "id")
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

var isNameExistsSql = fmt.Sprintf("select count(*) from %s where name=?", app.TABLE_USERS)
func IsNameExists(dbmap *gorp.DbMap, name string) bool {
    count, err := dbmap.SelectInt(isNameExistsSql, name)
    if nil != err {
        revel.ERROR.Printf("db query failed: %s", err)
        panic(err)
    }
    return 0 != count
}

var isEmailExistsSql = fmt.Sprintf("select count(*) from %s where email=?", app.TABLE_USERS)
func IsEmailExists(dbmap *gorp.DbMap, email string) bool {
    count, err := dbmap.SelectInt(isEmailExistsSql, email)
    if nil != err {
        revel.ERROR.Printf("db query failed: %s", err)
        panic(err)
    }
    return 0 != count
}

var isNameEmailExistsSql = fmt.Sprintf("select count(*) from %s where name=? and email=?", app.TABLE_USERS)
func IsNameEmailExists(dbmap *gorp.DbMap, name, email string) bool {
    count, err := dbmap.SelectInt(isNameEmailExistsSql, name, email)
    if nil != err {
        revel.ERROR.Printf("db query failed: %s", err)
        panic(err)
    }
    return 0 != count
}

var isEmailVerifiedSql = fmt.Sprintf("select count(*) from %s where name=? and email_verified!=0", app.TABLE_USERS)
func IsEmailVerified(dbmap *gorp.DbMap, name string) bool {
    count, err := dbmap.SelectInt(isEmailVerifiedSql, name)
    if nil != err {
        revel.ERROR.Printf("db query failed: %s", err)
        panic(err)
    }
    return 0 != count;
}

var findUserByNameSql = fmt.Sprintf("select * from %s where name=?", app.TABLE_USERS)
func FindUserByName(dbmap *gorp.DbMap, name string) (u *models.User) {
    u = new(models.User)
    err := dbmap.SelectOne(u, findUserByNameSql, name)
    if nil != err {
        revel.ERROR.Fatalf("failed to select user by name %s: %s", name, err)
        return nil
    }
    return
}

var findUserByEmailSql = fmt.Sprintf("select * from %s where email=?", app.TABLE_USERS)
func FindUserByEmail(dbmap *gorp.DbMap, email string) (u *models.User) {
    u = new(models.User)
    err := dbmap.SelectOne(u, findUserByEmailSql, email)
    if nil != err {
        revel.ERROR.Printf("failed to select user by email %s: %s", email, err)
        return nil
    }
    return
}

var findUserByIDSql = fmt.Sprintf("select * from %s where id=?", app.TABLE_USERS)
func FindUserByID(dbmap *gorp.DbMap, id int32) (u *models.User) {
    u = new(models.User)
    err := dbmap.SelectOne(u, findUserByIDSql, id)
    if nil != err {
        revel.ERROR.Printf("failed to select user by id %d: %s", id, err)
        return nil
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

