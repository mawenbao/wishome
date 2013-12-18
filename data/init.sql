drop database if exists wishomedb;
create database if not exists wishomedb;

grant all privileges on wishomedb.* to wishome@localhost identified by 'wishome';
flush privileges;

