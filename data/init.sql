drop database if exists wishomedb;
create database if not exists wishomedb;

grant all privileges on wishomedb.* to wishome@localhost identified by 'wishome';
grant all privileges on wishomedb.* to wishome@127.0.0.1 identified by 'wishome';
flush privileges;

