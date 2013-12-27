# wishome

A simple web site built with revel.

## Dependencies

*  A local smtp server

## TODO

1. <del>send email for confirmation after user signup</del>
2. <del>reset password via email</del>
3. <del>change db handle from gorp.DbMap to database.DbManager in all the database related validate functions(validators/user.go) and database helper functions(database/db.go), rewrite db logic, use a global db handle now with database/sql connection pool support</del>
4. <del>captcha verification</del>, <del>should refresh cpatcha id at server side</del>
5. <del>use message translation in modules/validators package too</del>
6. <del>move session related logic from controller to modules/session</del>
7. <del>Use a custom struct type for custom global config.</del>
8. <del>Move input autofocus related logic from template to javascript.</del>
9. <del>count action response time in admin page</del>

## Issues
1. <del>Fix: resetpass.html template will autofocus on the email input even if name is not set.</del>
2. <del>Fix: admin page controller.Request.RemoteAddr is wrong when using nginx as proxy frontend.</del>
3. <del>Fix: race conditions may happen in caching module of timer results. **USE Redis Instead**.</del> Use in-memory cache now
4. Fix: multiple wishome instances will lead to mismatched captcha storage. Try to use redis to store captcha.
