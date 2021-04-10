create table commodities
(
	id int auto_increment,
	name varchar(255) not null unique,
	description text not null,
	created_at timestamp not null,
	updated_at timestamp not null,
	constraint commodities_pk
		primary key (id)
);

