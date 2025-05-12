create table if not exists public.contacts
(
    id BIGSERIAL primary key,
    first_name text not null,
    last_name text not null
);