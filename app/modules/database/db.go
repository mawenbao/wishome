package database

import (
    "fmt"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
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
    db, err := sql.Open(app.MyGlobal[app.CONFIG_DB_DRIVER].(string), app.MyGlobal[app.CONFIG_DB_SPEC].(string))
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

var isNameExistsSql = fmt.Sprintf("select count(*) from %s where name=?", app.TABLE_USERS)
func IsNameExists(name string) bool {
    count, err := MyDbManager.Db().SelectInt(isNameExistsSql, name)
    if nil != err {
        revel.ERROR.Printf("db query failed: %s", err)
        panic(err)
    }
    return 0 != count
}

var isEmailExistsSql = fmt.Sprintf("select count(*) from %s where email=?", app.TABLE_USERS)
func IsEmailExists(email string) bool {
    count, err := MyDbManager.Db().SelectInt(isEmailExistsSql, email)
    if nil != err {
        revel.ERROR.Printf("db query failed: %s", err)
        panic(err)
    }
    return 0 != count
}

var isNameEmailExistsSql = fmt.Sprintf("select count(*) from %s where name=? and email=?", app.TABLE_USERS)
func IsNameEmailExists(name, email string) bool {
    count, err := MyDbManager.Db().SelectInt(isNameEmailExistsSql, name, email)
    if nil != err {
        revel.ERROR.Printf("db query failed: %s", err)
        panic(err)
    }
    return 0 != count
}

var isEmailVerifiedSql = fmt.Sprintf("select count(*) from %s where name=? and email_verified!=0", app.TABLE_USERS)
func IsEmailVerified(name string) bool {
    count, err := MyDbManager.Db().SelectInt(isEmailVerifiedSql, name)
    if nil != err {
        revel.ERROR.Printf("db query failed: %s", err)
        panic(err)
    }
    return 0 != count;
}

var findUserByNameSql = fmt.Sprintf("select * from %s where name=?", app.TABLE_USERS)
func FindUserByName(name string) (u *models.User) {
    u = new(models.User)
    err := MyDbManager.Db().SelectOne(u, findUserByNameSql, name)
    if nil != err {
        revel.ERROR.Fatalf("failed to select user by name %s: %s", name, err)
        return nil
    }
    return
}

var findUserByEmailSql = fmt.Sprintf("select * from %s where email=?", app.TABLE_USERS)
func FindUserByEmail(email string) (u *models.User) {
    u = new(models.User)
    err := MyDbManager.Db().SelectOne(u, findUserByEmailSql, email)
    if nil != err {
        revel.ERROR.Printf("failed to select user by email %s: %s", email, err)
        return nil
    }
    return
}

var findUserByIDSql = fmt.Sprintf("select * from %s where id=?", app.TABLE_USERS)
func FindUserByID(id int32) (u *models.User) {
    u = new(models.User)
    err := MyDbManager.Db().SelectOne(u, findUserByIDSql, id)
    if nil != err {
        revel.ERROR.Printf("failed to select user by id %d: %s", id, err)
        return nil
    }
    return
}

func SaveUser(u models.User) bool {
    err := MyDbManager.Db().Insert(&u)
    if nil != err {
        revel.ERROR.Printf("failed to insert user %s: %s", u, err)
        return false
    }
    return true
}

func UpdateUser(u models.User) bool {
    count, err := MyDbManager.Db().Update(&u)
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
func RemoveUser(u models.User) bool {
    count, err := MyDbManager.Db().Delete(&u)
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

