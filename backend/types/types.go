package types

type Match struct {
	Round	string	`json:"round"`
	Date	string	`json:"date"`
	Time	string	`json:"time"`
	Team1	string	`json:"team1"`
	Team2	string	`json:"team2"`
}

type User struct {
	Id				int		`json:"id"`
	Username		string	`json:"username"`
	HashedPassword	string	`json:"password"`
}

type Booking struct {
	Id			int		`json:"id"`
	UserId		int		`json:"user_id"`
	Username	string	`json:"username"`
	Match		Match	`json:"match"`
	Tickets		int		`json:"n_of_tickets"`
}

type RegisterAndLoginPayload struct {
	Email		string `json:"email" binding:"required,email"`
	Password	string `json:"password" binding:"required,min=8"`
}

type BookingPayload struct {
	GameId		int `json:"gameId"`
	NOfTickets	int `json:"tickets"`
}

type Email struct {
	To		string
	Subject	string
	Body	string
}

type UserRepo interface {
	AddUser(email, hashedPassword string) error
	GetUserByEmail(email string) (*User, error)
}

type MatchRepo interface {
	GetMatchByMatchId(matchId int) (*Match, error)
	GetMatchesByTeam(teamId int) (*[]Match, error)
	GetMatchesByMatchday(matchday string) (*[]Match, error)
}

type BookingRepo interface {
	GetMatchByMatchId(MatchId int) (*Match, error)
	AddBooking(userId int, matchId int, nOfTickets int) (int, error)
	GetBookingInfo(bookingId int) (*Booking, error)
}