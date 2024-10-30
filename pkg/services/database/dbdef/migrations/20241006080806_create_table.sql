-- +goose Up

-- Create the user_role enum type
CREATE TYPE user_role AS ENUM('admin', 'lecture', 'student');

-- +goose StatementBegin
CREATE TABLE departments (
                             code VARCHAR(255) PRIMARY KEY,
                             name VARCHAR(255) NOT NULL,
                             description TEXT,
                             created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                             updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE users (
                       id BIGSERIAL PRIMARY KEY,
                       nim_nip VARCHAR(100) NOT NULL,
                       name VARCHAR(255) NOT NULL,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       password TEXT NOT NULL,
                       role user_role NOT NULL,
                       department_code VARCHAR(255),
                       entry_year INT,
                       status VARCHAR(50),
                       address TEXT,
                       photo_url TEXT,
                       created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
                       deleted_at TIMESTAMP,

                       CONSTRAINT fk_department
                           FOREIGN KEY (department_code)
                               REFERENCES departments(code)
                               ON DELETE SET NULL
);

CREATE TABLE academic_years (
                                id BIGSERIAL PRIMARY KEY,
                                year INT NOT NULL,
                                semester INT NOT NULL,
                                is_active BOOLEAN NOT NULL DEFAULT FALSE,
                                start_date DATE NOT NULL,
                                end_date DATE NOT NULL,
                                description TEXT,
                                created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE courses (
                         id BIGSERIAL PRIMARY KEY,
                         code VARCHAR(20) NOT NULL UNIQUE,
                         name VARCHAR(255) NOT NULL,
                         credits INT NOT NULL,
                         semester INT NOT NULL,
                         department_code VARCHAR(255) NOT NULL REFERENCES departments(code),
                         description TEXT,
                         is_active BOOLEAN NOT NULL DEFAULT TRUE,
                         created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                         updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE course_prerequisites (
                                      course_id BIGINT REFERENCES courses(id),
                                      prerequisite_id BIGINT REFERENCES courses(id),
                                      created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                      PRIMARY KEY (course_id, prerequisite_id)
);

CREATE TABLE study_plans (
                             id BIGSERIAL PRIMARY KEY,
                             student_id BIGINT NOT NULL REFERENCES users(id),
                             academic_year_id BIGINT NOT NULL REFERENCES academic_years(id),
                             status VARCHAR(20) NOT NULL CHECK (status IN ('draft', 'submitted', 'approved', 'rejected')),
                             max_credits INT NOT NULL,
                             total_credits INT NOT NULL,
                             gpa DECIMAL(4,2),
                             advisor_id BIGINT NOT NULL REFERENCES users(id),
                             notes TEXT,
                             submitted_at TIMESTAMP,
                             approved_at TIMESTAMP,
                             rejected_at TIMESTAMP,
                             created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                             updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE study_plan_details (
                                    id BIGSERIAL PRIMARY KEY,
                                    study_plan_id BIGINT NOT NULL REFERENCES study_plans(id),
                                    course_id BIGINT NOT NULL REFERENCES courses(id),
                                    status VARCHAR(20) NOT NULL CHECK (status IN ('enrolled', 'completed', 'dropped')),
                                    grade DECIMAL(4,2),
                                    graded_by BIGINT REFERENCES users(id),
                                    graded_at TIMESTAMP,
                                    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE class_schedules (
                                 id BIGSERIAL PRIMARY KEY,
                                 course_id BIGINT NOT NULL REFERENCES courses(id),
                                 lecturer_id BIGINT NOT NULL REFERENCES users(id),
                                 academic_year_id BIGINT NOT NULL REFERENCES academic_years(id),
                                 day_of_week INT NOT NULL CHECK (day_of_week BETWEEN 1 AND 7),
                                 start_time TIME NOT NULL,
                                 end_time TIME NOT NULL,
                                 room VARCHAR(50) NOT NULL,
                                 quota INT NOT NULL,
                                 enrolled INT NOT NULL DEFAULT 0,
                                 created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                 updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE assignments (
                             id BIGSERIAL PRIMARY KEY,
                             class_schedule_id BIGINT NOT NULL REFERENCES class_schedules(id),
                             title VARCHAR(255) NOT NULL,
                             description TEXT,
                             due_date TIMESTAMP NOT NULL,
                             max_score DECIMAL(5,2) NOT NULL,
                             weight DECIMAL(5,2) NOT NULL,
                             type VARCHAR(20) NOT NULL CHECK (type IN ('homework', 'quiz', 'project', 'exam')),
                             instructions TEXT,
                             attachment_url TEXT,
                             created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                             updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE assignment_submissions (
                                        id BIGSERIAL PRIMARY KEY,
                                        assignment_id BIGINT NOT NULL REFERENCES assignments(id),
                                        student_id BIGINT NOT NULL REFERENCES users(id),
                                        submission_url TEXT NOT NULL,
                                        score DECIMAL(5,2),
                                        feedback TEXT,
                                        status VARCHAR(20) NOT NULL CHECK (status IN ('submitted', 'late', 'graded')),
                                        submitted_at TIMESTAMP NOT NULL,
                                        graded_by BIGINT REFERENCES users(id),
                                        graded_at TIMESTAMP,
                                        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                        updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE attendance (
                            id BIGSERIAL PRIMARY KEY,
                            class_schedule_id BIGINT NOT NULL REFERENCES class_schedules(id),
                            student_id BIGINT NOT NULL REFERENCES users(id),
                            status VARCHAR(20) NOT NULL CHECK (status IN ('present', 'absent', 'excused')),
                            date DATE NOT NULL,
                            created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                            updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- Adding indexes for better query performance
CREATE INDEX idx_users_nim_nip ON users(nim_nip);
CREATE INDEX idx_users_department_code ON users(department_code);
CREATE INDEX idx_users_entry_year ON users(entry_year);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_academic_years_year ON academic_years(year);
CREATE INDEX idx_academic_years_is_active ON academic_years(is_active);
CREATE INDEX idx_courses_department ON courses(department_code);
CREATE INDEX idx_courses_code ON courses(code);
CREATE INDEX idx_study_plans_student ON study_plans(student_id);
CREATE INDEX idx_study_plans_academic_year ON study_plans(academic_year_id);
CREATE INDEX idx_study_plan_details_study_plan ON study_plan_details(study_plan_id);
CREATE INDEX idx_class_schedules_course ON class_schedules(course_id);
CREATE INDEX idx_class_schedules_academic_year ON class_schedules(academic_year_id);
CREATE INDEX idx_assignments_class_schedule ON assignments(class_schedule_id);
CREATE INDEX idx_assignment_submissions_assignment ON assignment_submissions(assignment_id);
CREATE INDEX idx_attendance_class_schedule ON attendance(class_schedule_id);
CREATE INDEX idx_attendance_student ON attendance(student_id);

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS attendance;
DROP TABLE IF EXISTS assignment_submissions;
DROP TABLE IF EXISTS assignments;
DROP TABLE IF EXISTS class_schedules;
DROP TABLE IF EXISTS study_plan_details;
DROP TABLE IF EXISTS study_plans;
DROP TABLE IF EXISTS course_prerequisites;
DROP TABLE IF EXISTS courses;
DROP TABLE IF EXISTS academic_years;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS departments;

-- Drop the user_role enum type
DROP TYPE IF EXISTS user_role;
-- +goose StatementEnd