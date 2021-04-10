create table warehouse_db.warehouses
(
    id         bigint auto_increment
        primary key,
    name       varchar(255) not null,
    latitude   float        not null,
    longitude  float        not null,
    created_at timestamp    not null,
    updated_at timestamp    not null
);