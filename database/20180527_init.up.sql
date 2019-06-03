SET statement_timeout = 60000;
SET lock_timeout = 60000;

CREATE TABLE payments
(
    type text not null,
    id uuid not null primary key,
    version smallint not null default 0,
    organisation_id uuid not null,
    attributes jsonb not null default json_object('{}')
);

