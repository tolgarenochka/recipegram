create database recipegram
    with owner postgres;

create sequence public.table_id_seq;

create table public."user"
(
    id          integer default nextval('table_id_seq'::regclass)           not null,
    mail        text                                                        not null,
    pas         text                                                        not null,
    photo       varchar default 'user_photo/default.png'::character varying not null,
    first_name  text                                                        not null,
    second_name text                                                        not null,
    third_name  text
);

alter table public."user"
    add constraint user_id
        primary key (id);






