create table messages
(
    id bigserial
        constraint message_pk
            primary key,
    group_id bigint not null,
    message_id bigint not null,
    source varchar not null,
    message_at timestamp not null,
    created_at timestamp not null,
    updated_at timestamp not null
);

create index messages_group_id_index ON messages (group_id);
