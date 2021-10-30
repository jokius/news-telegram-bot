create table groups
(
    id bigserial
        constraint group_pk
            primary key,
    user_id bigint not null,
    source_name varchar not null,
    group_name varchar not null,
    last_update_at timestamp not null,
    created_at timestamp not null,
    updated_at timestamp not null
);

create index groups_user_id_index ON groups (user_id);
