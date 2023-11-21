package user

type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"`
	Level    Level  `json:"level"`
}

type Level string

const (
	RegularUser Level = "Regular User"
	AdminUser   Level = "Administrator User"
)

type Follow struct {
	From string `json:"from" binding:"required"`
	To   string `json:"to" binding:"required"`
}
