package models

import (
	"time"
)

type Role string

const (
	RoleStudent  Role = "student"
	RoleLecturer Role = "lecturer"
	RoleAdmin    Role = "admin"
)

// Department represents an academic department
type Department struct {
	Name        string    `db:"name" json:"name"`
	Code        string    `db:"code" json:"code"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// AcademicYear represents an academic year period
type AcademicYear struct {
	ID        int64     `db:"id" json:"id"`
	Year      int       `db:"year" json:"year"`
	Semester  int       `db:"semester" json:"semester"`
	IsActive  bool      `db:"is_active" json:"is_active"`
	StartDate time.Time `db:"start_date" json:"start_date"`
	EndDate   time.Time `db:"end_date" json:"end_date"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type BaseUser struct {
	ID             int64      `db:"id" json:"id"`
	NimNip         string     `db:"nim_nip" json:"nim_nip"`
	Name           string     `db:"name" json:"name"`
	Email          string     `db:"email" json:"email"`
	Role           Role       `db:"role" json:"role"` // e.g., "student", "lecturer", "admin"
	DepartmentCode string     `db:"department_code" json:"department_code"`
	EntryYear      int        `db:"entry_year" json:"entry_year,omitempty"`
	Status         string     `db:"status" json:"status"`             // Added status (active/inactive/graduated)
	Address        string     `db:"address" json:"address,omitempty"` // Added address
	PhotoURL       *string    `db:"photo_url" json:"photo_url,omitempty"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt      *time.Time `db:"deleted_at" json:"deleted_at,omitempty"` // Added soft delete
}

// User represents any user in the system (student, lecturer, admin)
type User struct {
	BaseUser        // Embeds BaseUser, inheriting all its fields
	Password string `db:"password" json:"password"` // Only added to User struct
}

// Course represents an academic course
type Course struct {
	ID             int64     `db:"id" json:"id"`
	Code           string    `db:"code" json:"code"`
	Name           string    `db:"name" json:"name"`
	Credits        int       `db:"credits" json:"credits"`
	Semester       int       `db:"semester" json:"semester"`
	DepartmentCode string    `db:"department_code" json:"department_code"`
	Description    string    `db:"description" json:"description"`     // Added description
	Prerequisites  []int64   `db:"prerequisites" json:"prerequisites"` // Added prerequisites
	IsActive       bool      `db:"is_active" json:"is_active"`         // Added active status
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

// StudyPlan represents a student's study plan for a semester
type StudyPlan struct {
	ID             int64      `db:"id" json:"id"`
	StudentID      int64      `db:"student_id" json:"student_id"`
	AcademicYearID int64      `db:"academic_year_id" json:"academic_year_id"`
	Status         string     `db:"status" json:"status"` // draft/submitted/approved/rejected
	MaxCredits     int        `db:"max_credits" json:"max_credits"`
	TotalCredits   int        `db:"total_credits" json:"total_credits"`
	GPA            float64    `db:"gpa" json:"gpa"`               // Added GPA
	AdvisorID      int64      `db:"advisor_id" json:"advisor_id"` // Added academic advisor
	Notes          string     `db:"notes" json:"notes,omitempty"` // Added notes
	SubmittedAt    *time.Time `db:"submitted_at" json:"submitted_at,omitempty"`
	ApprovedAt     *time.Time `db:"approved_at" json:"approved_at,omitempty"`
	RejectedAt     *time.Time `db:"rejected_at" json:"rejected_at,omitempty"` // Added rejection tracking
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`
}

// StudyPlanDetail represents individual courses in a study plan
type StudyPlanDetail struct {
	ID          int64      `db:"id" json:"id"`
	StudyPlanID int64      `db:"study_plan_id" json:"study_plan_id"`
	CourseID    int64      `db:"course_id" json:"course_id"`
	Status      string     `db:"status" json:"status"` // enrolled/completed/dropped
	Grade       *float64   `db:"grade" json:"grade,omitempty"`
	GradedBy    *int64     `db:"graded_by" json:"graded_by,omitempty"` // Added grader reference
	GradedAt    *time.Time `db:"graded_at" json:"graded_at,omitempty"` // Added grading timestamp
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}

// ClassSchedule represents course schedules
type ClassSchedule struct {
	ID             int64     `db:"id" json:"id"`
	CourseID       int64     `db:"course_id" json:"course_id"`
	LecturerID     int64     `db:"lecturer_id" json:"lecturer_id"`
	AcademicYearID int64     `db:"academic_year_id" json:"academic_year_id"`
	DayOfWeek      int       `db:"day_of_week" json:"day_of_week"`
	StartTime      time.Time `db:"start_time" json:"start_time"`
	EndTime        time.Time `db:"end_time" json:"end_time"`
	Room           string    `db:"room" json:"room"`
	Quota          int       `db:"quota" json:"quota"`
	Enrolled       int       `db:"enrolled" json:"enrolled"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

// Assignment represents course assignments
type Assignment struct {
	ID              int64     `db:"id" json:"id"`
	ClassScheduleID int64     `db:"class_schedule_id" json:"class_schedule_id"`
	Title           string    `db:"title" json:"title"`
	Description     string    `db:"description" json:"description"`
	DueDate         time.Time `db:"due_date" json:"due_date"`
	MaxScore        float64   `db:"max_score" json:"max_score"`
	Weight          float64   `db:"weight" json:"weight"` // Percentage weight in final grade
	Type            string    `db:"type" json:"type"`     // homework/quiz/project/exam
	Instructions    string    `db:"instructions" json:"instructions"`
	AttachmentURL   string    `db:"attachment_url" json:"attachment_url,omitempty"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

// AssignmentSubmission represents student assignment submissions
type AssignmentSubmission struct {
	ID            int64      `db:"id" json:"id"`
	AssignmentID  int64      `db:"assignment_id" json:"assignment_id"`
	StudentID     int64      `db:"student_id" json:"student_id"`
	SubmissionURL string     `db:"submission_url" json:"submission_url"`
	Score         *float64   `db:"score" json:"score,omitempty"`
	Feedback      string     `db:"feedback" json:"feedback,omitempty"`
	Status        string     `db:"status" json:"status"` // submitted/late/graded
	GradedBy      *int64     `db:"graded_by" json:"graded_by,omitempty"`
	GradedAt      *time.Time `db:"graded_at" json:"graded_at,omitempty"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`
}

// Attendance represents class attendance records
type Attendance struct {
	ID              int64     `db:"id" json:"id"`
	ClassScheduleID int64     `db:"class_schedule_id" json:"class_schedule_id"`
	StudentID       int64     `db:"student_id" json:"student_id"`
	Status          string    `db:"status" json:"status"` // present/absent/excused
	Date            time.Time `db:"date" json:"date"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}
