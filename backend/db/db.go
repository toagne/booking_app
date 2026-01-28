package db

import (
	"database/sql"
	"log"
	"time"
	"os"
	"github.com/go-sql-driver/mysql"
	"net/http"
	"fmt"
	"encoding/json"
	"strings"
)

var db *sql.DB

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

func getTeamId(name string) (int, error) {
	var id int
	err := db.QueryRow(`SELECT id FROM teams WHERE name = ?`, name).Scan(&id)
	return id, err
}

func createTeamsTable(data struct{Matches []Match}) {
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS teams (
			id INT AUTO_INCREMENT NOT NULL,
			name VARCHAR(50),
			PRIMARY KEY (id),
			UNIQUE KEY uniq_team (name)
		)
	`); err != nil {
		log.Fatal("Could not create matches table: ", err)
	}

	for _, match := range data.Matches {
		if match.Round == "Matchday 20" {
			break
		}
		_, err := db.Exec(`INSERT INTO teams (name)
			VALUES (?)`,
			match.Team1,
		)
		if err != nil {
			if !strings.Contains(err.Error(), "Duplicate entry") {
				log.Fatal("add matches: ", err)
			}
		}
	}
}

func createMatchesTable(data struct{Matches []Match}) {
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS matches (
			id INT AUTO_INCREMENT NOT NULL,
			round VARCHAR(50) NOT NULL,
			date VARCHAR(20) NOT NULL,
			time VARCHAR(20) NOT NULL,
			team1 VARCHAR(50) NOT NULL,
			team1_id INT NOT NULL,
			team2 VARCHAR(50) NOT NULL,
			team2_id INT NOT NULL,
			PRIMARY KEY (id),
			UNIQUE KEY uniq_match (round, date, time, team1, team2)
		)
	`); err != nil {
		log.Fatal("Could not create matches table: ", err)
	}

	for _, match := range data.Matches {
		id1, err := getTeamId(match.Team1)
		if err != nil {
			log.Fatalf("team1 not found: %v", err)
		}
		id2, err := getTeamId(match.Team2)
		if err != nil {
			log.Fatalf("team2 not found: %v", err)
		}

		_, err = db.Exec(`INSERT INTO matches (round, date, time, team1, team1_id, team2, team2_id)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			match.Round, match.Date, match.Time, match.Team1, id1, match.Team2, id2,
		)
		if err != nil {
			if !strings.Contains(err.Error(), "Duplicate entry") {
				log.Fatal("add matches: ", err)
			}
		}
	}
}

func createUsersTable() {
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT NOT NULL,
			username VARCHAR(100) NOT NULL,
			password TEXT NOT NULL,
			PRIMARY KEY (id),
			UNIQUE KEY uniq_user (username)
		)
	`); err != nil {
		log.Fatal("Could not create users table: ", err)
	}
}

func createBookingsTable() {
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS bookings (
			id INT AUTO_INCREMENT NOT NULL,
			user_id INT NOT NULL,
			match_id INT NOT NULL,
			n_of_tickets INT NOT NULL,
			PRIMARY KEY (id)
		)
	`); err != nil {
		log.Fatal("Could not create bookings table: ", err)
	}
}

func createTables() {
	resp, err := http.Get("https://raw.githubusercontent.com/openfootball/football.json/master/2025-26/it.1.json")
	if err != nil {
		log.Fatalf("Error fetching URL: %v", err)
	}
	defer resp.Body.Close()

	var data struct {
		Matches []Match
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}

	createTeamsTable(data)
	createMatchesTable(data)
	createUsersTable()
	createBookingsTable()
}

func InitDb() {
	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DB_USER")
	cfg.Passwd = os.Getenv("DB_PASSWORD")
	cfg.Net = "tcp"
	cfg.Addr = "db"
	cfg.DBName = "booking_db"

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	for i := 1; i <= 10; i++ {
		if err := db.Ping(); err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	fmt.Println("Database Connected!")

	createTables()
}

func GetMatchesByMatchday(matchday string) ([]Match, error) {
	rows, err := db.Query(`
		SELECT date, time, team1, team2
		FROM matches
		WHERE round = ?
	`, matchday)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []Match

	for rows.Next() {
		var m Match
		if err := rows.Scan(
			&m.Date,
			&m.Time,
			&m.Team1,
			&m.Team2,
		); err != nil {
			return nil, err
		}
		matches = append(matches, m)
	}

	return matches, rows.Err()
}

func GetMatchesByTeam(teamId int) ([]Match, error) {
	rows, err := db.Query(`
		SELECT round, date, time, team1, team2
		FROM matches
		WHERE team1_id LIKE ? OR team2_id like ?
	`, teamId, teamId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []Match

	for rows.Next() {
		var m Match
		if err := rows.Scan(
			&m.Round,
			&m.Date,
			&m.Time,
			&m.Team1,
			&m.Team2,
		); err != nil {
			return nil, err
		}
		matches = append(matches, m)
	}

	return matches, rows.Err()
}

func GetMatchByMatchId(MatchId int) (Match, error) {
	rows, err := db.Query(`
		SELECT round, date, time, team1, team2
		FROM matches
		WHERE id = ?
	`, MatchId)
	if err != nil {
		return Match{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		return Match{}, sql.ErrNoRows
	}

	var m Match
	if err := rows.Scan(
		&m.Round,
		&m.Date,
		&m.Time,
		&m.Team1,
		&m.Team2,
	); err != nil {
		return Match{}, err
	}

	return m, rows.Err()
}

func AddUser(username string, password string) error {
	if _, err := db.Exec(`
		INSERT INTO users (username, password) VALUES (?, ?)
	`, username, password); err != nil {
		if !strings.Contains(err.Error(), "Duplicate entry") {
			log.Print("add user: ", err)
			return err
		}
	}
	return nil
}

func GetUserByEmail(email string) (User, error) {
	var user User
	err := db.QueryRow(`SELECT * FROM users WHERE username=?`, email).Scan(
		&user.Id,
		&user.Username,
		&user.HashedPassword,
	)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("no user with username %v\n", email)
	case err != nil:
		log.Fatalf("query error: %v\n", err)
	default:
		log.Printf("username %v in database\n", email)
	}
	return user, nil
}

func AddBooking(userId int, matchId int, nOfTickets int) (int, error) {
	res, err := db.Exec(`
		INSERT INTO bookings (user_id, match_id, n_of_tickets)
		VALUES (?, ?, ?)
	`, userId, matchId, nOfTickets)
	if err != nil {
		log.Print("add booking: ", err)
		return 0, err
	}

	bookingId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(bookingId), nil
}

func GetBookingInfo(bookingId int) (*Booking, error) {
	var booking Booking
	err := db.QueryRow(`
	SELECT bookings.id,
			bookings.user_id,
			users.username,
			matches.round,
			matches.date,
			matches.time,
			matches.team1,
			matches.team2,
			bookings.n_of_tickets
		FROM bookings
		JOIN matches
		ON bookings.match_id = matches.id
		JOIN users
		ON bookings.user_id = users.id
		WHERE bookings.id = ?
	`, bookingId).Scan(&booking.Id,
						&booking.UserId,
						&booking.Username,
						&booking.Match.Round,
						&booking.Match.Date,
						&booking.Match.Time,
						&booking.Match.Team1,
						&booking.Match.Team2,
						&booking.Tickets,
	)

	if err != nil {
		return nil, err
	}

	return &booking, nil
}