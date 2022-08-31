create table if not exists user_data
(
    id      varchar(36)   not null
        primary key,
    name    varchar(100)  null,
    surname varchar(100)  null,
    age     smallint      null,
    gender  char          null,
    hobbies varchar(1000) null,
    city    varchar(100)  null
);

create table if not exists user_access
(
    login    varchar(20)  not null
        primary key,
    password text not null,
    user_id  varchar(36)  null,
    constraint user_access_user_data_id_fk
        foreign key (user_id) references user_data (id)
            on update cascade on delete cascade
);

create table if not exists friends
(
    created_time timestamp not null default now(),
    user_id   varchar(36) not null,
    friend_id varchar(36) not null,
    primary key (user_id, friend_id),
    constraint friend_id_fk
        foreign key (friend_id) references user_data (id)
            on update cascade on delete cascade,
    constraint user_id_fk
        foreign key (user_id) references user_data (id)
            on update cascade on delete cascade
);