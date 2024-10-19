package database

type User struct {
	ID       int64  `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
	Role     string `db:"role" json:"role"`
	Major    string `db:"major" json:"major"`
	Year     int    `db:"year" json:"year"`
	Phone    string `db:"phone" json:"phone"`
	Group    int    `db:"group" json:"group"`
}
