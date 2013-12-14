use wishomedb;

drop table if exists users;
create table if not exists users(
    id int not null auto_increment,
    name varchar(15) not null,
    email varchar(50) not null,
    email_verified boolean not null default 0,
    password char(32) not null,
    salt char(255) not null,
    primary key(id)
) AUTO_INCREMENT = 10;

#drop index users_name on users;
#drop index users_email on users;

alter table users add unique index users_name(name);
alter table users add unique index users_email(email);

