package models

import (
    "fmt"
    "github.com/mawenbao/wishome/app/modules/common"
)

type User struct {
    ID int32 `db:"id"`
    Name string `db:"name"`
    Email string `db:"email"`
    Password string `db:"password"`
    PassSalt string `db:"salt"` // 255 byte
}

func (u *User) String() string {
    return fmt.Sprintf("User(%s)", u.Name)
}

func (u *User) EncryptPass() {
    u.PassSalt = newPassSalt(u)
    u.Password = common.MD5Sum(u.Password + u.PassSalt)
}

// salt = [random(50) + email + random int(100) + name + random int(105)][:255]
func newPassSalt(u *User) string {
    return fmt.Sprintf(
        "%s%s%s%s%s",
        common.NewRandomString(50),
        u.Email,
        common.NewRandomString(100),
        u.Name,
        common.NewRandomString(105),
    )[:255]
}

