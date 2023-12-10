CREATE SCHEMA rest_app;

set search_path=rest_app;

create table users(id uuid not null primary key, email varchar(255));

insert into users values('2f3d92a9-f4db-497b-bf16-367bf2ad7e20'::uuid, 'boss@company.com');