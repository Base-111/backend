CREATE TABLE IF NOT EXISTS admin.prompts
(
    id        bigserial primary key,
    name      varchar(255),
    text      text not null,
    is_system boolean
)