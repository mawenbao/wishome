use wishomedb;

create table if not exists users(
    id int not null auto_increment,
    name varchar(15) not null,
    email varchar(50) not null,
    password char(32) not null,
    salt char(255) not null,
    primary key(id)
) AUTO_INCREMENT = 10;

alter table users add unique index users_name(name);
alter table users add unique index users_email(email);

