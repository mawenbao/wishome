package models

import (
    "fmt"
    "github.com/mawenbao/wishome/app/modules/common"
)

type User struct {
    ID int32 `db:"id"`
    Name string `db:"name"`
    Email string `db:"email"`
    EmailVerified bool `db:"email_verified"`
    Password string `db:"password"`
    PassSalt string `db:"salt"` // 255 byte
}

// make record of user reseting password
type ResetPassRecord struct {
    ID int32 `db:"id"`
    UserID int32 `db:"user_id"`
    ExpireTime int64 `db:"expire_time"`
    FinishTime int64 `db:"finish_time"`
}

func (u *User) String() string {
    return fmt.Sprintf("User(%s)", u.Name)
}

func (u *User) EncryptPass() {
    u.PassSalt = newPassSalt(u)
    u.Password = common.MD5Sum(u.Password + u.PassSalt)
}

func (u *User) IsValid() bool {
    return "" != u.Name && "" != u.Email
}

func (u *User) IsSecured() bool {
    return u.IsValid() && u.ID > 0 && "" != u.PassSalt && "" != u.Password
}

// salt = [readable random(50) + email + raw random(100) + name + raw random(105)][:255]
func newPassSalt(u *User) string {
    return fmt.Sprintf(
        "%s%s%s%s%s",
        common.NewReadableRandom(50),
        u.Email,
        common.NewRawRandom(100),
        u.Name,
        common.NewRawRandom(105),
    )[:255]
}

