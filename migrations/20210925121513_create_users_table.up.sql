create table users
(
    id bigserial
        constraint users_pk
            primary key,
    telegram_id bigint not null,
    created_at timestamp not null,
    updated_at timestamp not null
);

create unique index users_telegram_id_uindex
    on users (telegram_id);
