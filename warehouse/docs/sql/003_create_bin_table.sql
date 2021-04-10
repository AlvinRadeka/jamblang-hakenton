create table warehouse_db.bins
(
    id           bigint auto_increment
        primary key,
    warehouse_id bigint       not null,
    name         varchar(255) not null,
    latitude     float        not null,
    longitude    float        not null,
    created_at   timestamp    not null,
    updated_at   timestamp    not null
);

