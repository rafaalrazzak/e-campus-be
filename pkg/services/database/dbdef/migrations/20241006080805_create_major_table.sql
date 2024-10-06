-- +goose Up
-- +goose StatementBegin
CREATE TABLE majors (
    code varchar(5) primary key,
    name varchar(50) not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.majors;
-- +goose StatementEnd
