package database

import (
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "github.com/robfig/revel"
    "github.com/mawenbao/wishome/app"
    "github.com/mawenbao/wishome/app/models"
)

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

