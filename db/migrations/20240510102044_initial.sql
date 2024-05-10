-- +goose Up
-- +goose StatementBegin
create table files (
    id uint64 primary key unique,
    parent_id uint64 not null references files(id) on delete cascade,
    name text not null,
    size int not null,
    message_id int not null,
    file_id text not null,
    file_size uint32 not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

create index files_parent_name on files(parent_id, name);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table files;
drop index files_parent_name;
-- +goose StatementEnd
