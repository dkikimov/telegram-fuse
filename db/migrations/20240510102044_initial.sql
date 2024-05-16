-- +goose Up
-- +goose StatementBegin
create table file_entities (
    id integer primary key autoincrement unique,
    parent_id integer not null references files(id) on delete cascade,
    name text not null,
    message_id int default -1,
    file_id text default '-1',
    file_size uint32 default 0,
    created_at datetime not null default current_timestamp,
    updated_at datetime not null default current_timestamp
);

create index entity_parent_name on file_entities(parent_id, name);

insert into file_entities (id, parent_id, name) values (0, 0, 'root');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table file_entities;
-- +goose StatementEnd
