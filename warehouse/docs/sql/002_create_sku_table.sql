create table warehouse_db.skus
(
    id         bigint auto_increment
        primary key,
    sku        varchar(255) not null,
    name       text         not null,
    created_at timestamp    not null,
    updated_at timestamp    not null
);

