# wishome

A simple web site built with revel.

## Dependencies

*  A local smtp server

## TODO

1. <del>send email for confirmation after user signup</del>
2. <del>reset password via email</del>
3. <del>change db handle from gorp.DbMap to database.DbManager in all the database related validate functions(validators/user.go) and database helper functions(database/db.go)</del> rewrite db logic, use a global db handle now with database/sql connection pool support
4. captcha verification
5. <del>use message translation in modules/validators package too</del>
6. <del>move session related logic from controller to modules/session</del>

