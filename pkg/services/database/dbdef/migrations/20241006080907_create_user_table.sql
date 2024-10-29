-- +goose Up
-- +goose StatementBegin
CREATE TYPE user_role AS ENUM('ADMIN', 'LECTURE', 'STUDENT');

CREATE TABLE IF NOT EXISTS public."users" (
    id bigserial primary key,
    name varchar(100) not null,
    email varchar(100) unique not null,
    password text not null,
    role user_role not null,
    major varchar(5) not null,
    year int not null,
    phone varchar(25) not null,
    "group" int not null,
    picture text,

    constraint fk_major
    foreign key (major)
    references majors(code)
    );

CREATE INDEX idx_users_email ON public."users"(email);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE public."users";
-- +goose StatementEnd
