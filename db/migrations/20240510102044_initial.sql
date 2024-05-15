-- +goose Up
-- +goose StatementBegin
create table file_entities (
    id integer primary key autoincrement unique,
    parent_id integer not null references files(id) on delete cascade,
    name text not null,
    message_id int,
    file_id text,
    file_size uint32,
    created_at datetime not null default current_timestamp,
    updated_at datetime not null default current_timestamp
);

create index entity_parent_name on file_entities(parent_id, name);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table file_entities;
-- +goose StatementEnd
