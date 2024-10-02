package ecampus

// Role type as an enum
type Role string

// Define the possible roles
const (
	Admin   Role = "Admin"
	Lecture Role = "Lecture"
	Student Role = "Student"
)

// Major type as an enum
type Major string

// Define the possible majors
const (
	TI Major = "Teknik Informatika"
	SI Major = "Sistem Informasi"
	BD Major = "Bisnis Digital"
)

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     Role   `json:"role"`
	Major    Major  `json:"major"`
	Year     int    `json:"year"`
	Phone    string `json:"phone"`
	Group    int    `json:"group"`
}
