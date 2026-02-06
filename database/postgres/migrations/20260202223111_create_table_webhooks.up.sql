create table if not exists bot_webhooks_updates
(
    update_id  bigint                              not null
        constraint bot_webhooks_updates_pk
            unique,
    bot_name   varchar(30),
    data       json      default '{}'::jsonb       not null,
    created_at timestamp default current_timestamp not null
);

comment on table bot_webhooks_updates is 'Таблица с webhooks из telegram бота';

create index bot_webhooks_updates_bot_name_created_at_index
    on bot_webhooks_updates (bot_name asc, created_at desc);

