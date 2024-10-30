-- +goose Up
-- +goose StatementBegin

-- Insert initial data into departments table (IDs prefixed with 101)
INSERT INTO departments (code, name, description)
VALUES
    ('CS', 'Computer Science', 'Department of Computer Science'),
    ('BA', 'Business Administration', 'Department of Business Administration'),
    ('EE', 'Electrical Engineering', 'Department of Electrical Engineering');

-- Insert initial data into academic_years table (IDs prefixed with 201)
INSERT INTO academic_years (id, year, semester, is_active, start_date, end_date, description)
VALUES
    (2010000000001, 2024, 1, true, '2024-08-01', '2024-12-20', 'First semester of 2024'),
    (2010000000002, 2024, 2, false, '2025-01-10', '2025-05-30', 'Second semester of 2024');

-- Insert initial data into users table (IDs prefixed with 301)
INSERT INTO users (id, nim_nip, name, email, password, role, department_code, entry_year, status, address, created_at)
VALUES
    (3010000000001, '123456', 'Alice Johnson', 'alice@example.com', 'hashed_password', 'student', 'CS', 2022, 'active', '123 Main St', NOW()),
    (3010000000002, '789101', 'Bob Smith', 'bob@example.com', 'hashed_password', 'lecture', 'BA', 2023, 'active', '456 Elm St', NOW()),
    (3010000000003, '111213', 'Charlie Brown', 'charlie@example.com', 'hashed_password', 'admin', 'EE', 2020, 'active', '789 Oak St', NOW());

-- Insert initial data into courses table (IDs prefixed with 401)
INSERT INTO courses (id, code, name, credits, semester, department_code, description, is_active)
VALUES
    (4010000000001, 'CS101', 'Introduction to Computer Science', 3, 1, 'CS', 'Basics of Computer Science', true),
    (4010000000002, 'BA201', 'Marketing Basics', 3, 1, 'BA', 'Introduction to Marketing', true),
    (4010000000003, 'EE301', 'Circuit Analysis', 4, 1, 'EE', 'Fundamentals of Circuit Analysis', true);

-- Insert initial data into course_prerequisites table (IDs prefixed with 501)
INSERT INTO course_prerequisites (course_id, prerequisite_id)
VALUES
    (4010000000001, 4010000000002),
    (4010000000003, 4010000000001);

-- Insert initial data into study_plans table (IDs prefixed with 601)
INSERT INTO study_plans (id, student_id, academic_year_id, status, max_credits, total_credits, gpa, advisor_id, notes, submitted_at)
VALUES
    (6010000000001, 3010000000001, 2010000000001, 'submitted', 24, 18, 3.5, 3010000000002, 'Approved by advisor', NOW()),
    (6010000000002, 3010000000002, 2010000000002, 'approved', 20, 16, 3.8, 3010000000003, 'Needs improvement in practicals', NOW());

-- Insert initial data into study_plan_details table (IDs prefixed with 701)
INSERT INTO study_plan_details (id, study_plan_id, course_id, status, grade, graded_by, graded_at)
VALUES
    (7010000000001, 6010000000001, 4010000000001, 'completed', 3.7, 3010000000002, NOW()),
    (7010000000002, 6010000000002, 4010000000002, 'enrolled', NULL, NULL, NULL);

-- Insert initial data into class_schedules table (IDs prefixed with 801)
INSERT INTO class_schedules (id, course_id, lecturer_id, academic_year_id, day_of_week, start_time, end_time, room, quota, enrolled)
VALUES
    (8010000000001, 4010000000001, 3010000000002, 2010000000001, 2, '09:00', '10:30', 'Room A1', 30, 25),
    (8010000000002, 4010000000002, 3010000000003, 2010000000001, 4, '11:00', '12:30', 'Room B2', 25, 20);

-- Insert initial data into assignments table (IDs prefixed with 901)
INSERT INTO assignments (id, class_schedule_id, title, description, due_date, max_score, weight, type, instructions)
VALUES
    (9010000000001, 8010000000001, 'Homework 1', 'First homework assignment', '2024-09-15 23:59', 100, 10, 'homework', 'Complete the questions in Chapter 1'),
    (9010000000002, 8010000000002, 'Quiz 1', 'First quiz', '2024-09-20 10:00', 50, 20, 'quiz', 'Multiple-choice quiz covering initial topics');

-- Insert initial data into assignment_submissions table (IDs prefixed with 1001)
INSERT INTO assignment_submissions (id, assignment_id, student_id, submission_url, score, feedback, status, submitted_at, graded_by, graded_at)
VALUES
    (1001000000001, 9010000000001, 3010000000001, 'https://example.com/submissions/1', 85.0, 'Good work!', 'graded', NOW(), 3010000000002, NOW()),
    (1001000000002, 9010000000002, 3010000000002, 'https://example.com/submissions/2', 90.0, 'Excellent!', 'graded', NOW(), 3010000000003, NOW());

-- Insert initial data into attendance table (IDs prefixed with 1101)
INSERT INTO attendance (id, class_schedule_id, student_id, status, date)
VALUES
    (1101000000001, 8010000000001, 3010000000001, 'present', '2024-09-01'),
    (1101000000002, 8010000000001, 3010000000002, 'absent', '2024-09-01'),
    (1101000000003, 8010000000002, 3010000000003, 'present', '2024-09-01');

-- +goose StatementEnd
