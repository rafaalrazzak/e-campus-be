-- +goose Up
-- +goose StatementBegin

-- Insert initial data into departments table
INSERT INTO departments (code, name, description)
VALUES
    ('CS', 'Computer Science', 'Department of Computer Science'),
    ('BA', 'Business Administration', 'Department of Business Administration'),
    ('EE', 'Electrical Engineering', 'Department of Electrical Engineering');

-- Insert initial data into academic_years table
INSERT INTO academic_years (year, semester, is_active, start_date, end_date, description)
VALUES
    (2024, 1, true, '2024-08-01', '2024-12-20', 'First semester of 2024'),
    (2024, 2, false, '2025-01-10', '2025-05-30', 'Second semester of 2024');

-- Insert initial data into users table
INSERT INTO users (nim_nip, name, email, password, role, department_code, entry_year, status, address, created_at)
VALUES
    ('123456', 'Alice Johnson', 'alice@example.com', 'hashed_password', 'student', 'CS', 2022, 'active', '123 Main St', NOW()),
    ('789101', 'Bob Smith', 'bob@example.com', 'hashed_password', 'lecture', 'BA', 2023, 'active', '456 Elm St', NOW()),
    ('111213', 'Charlie Brown', 'charlie@example.com', 'hashed_password', 'admin', 'EE', 2020, 'active', '789 Oak St', NOW());

-- +goose StatementEnd
