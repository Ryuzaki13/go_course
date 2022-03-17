package db

import (
	"awesomeProject2/utils"
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	Login    string
	Password string
	Name     string
	Role     string
}

type Session struct {
	Hash   string
	User   User
	Date   string
	Device string
}

var query map[string]*sql.Stmt

func prepare() error {
	query = make(map[string]*sql.Stmt)

	//	Link.Exec(`CREATE TABLE "Session" (
	//    "Hash" varchar(64) primary key,
	//    "User" varchar not null,
	//    "Date" timestamp not null,
	//    "Device" varchar not null);
	//ALTER TABLE "Session" ADD CONSTRAINT "Session_fk0" FOREIGN KEY ("User") REFERENCES "User"("Login");
	//ALTER TABLE "User" ADD COLUMN "Role" varchar NOT NULL DEFAULT("");
	//`)

	var e error
	query["SessionInsert"], e = Link.Prepare(`INSERT INTO "Session"
		("Hash", "User", "Date", "Device")
		VALUES ($1, $2, CURRENT_TIMESTAMP, $3)`)
	if e != nil {
		fmt.Println(e)
		return e
	}

	query["SessionSelect"], e = Link.Prepare(`
SELECT "Hash", "Login", "Name", "Role", "Date", "Device"
FROM "Session" AS s
INNER JOIN "User" AS u
ON u."Login"=s."User"`)
	if e != nil {
		fmt.Println(e)
		return e
	}

	query["UserInsert"], e = Link.Prepare(`INSERT INTO "User"
		("Login", "Password", "Name", "Role")
		VALUES ($1, $2, $3, $4)`)
	if e != nil {
		fmt.Println(e)
		return e
	}

	query["UserSelect"], e = Link.Prepare(
		`SELECT "Login", "Password", "Name" FROM "User" ORDER BY "Name"`)
	if e != nil {
		fmt.Println(e)
		return e
	}

	query["UserLogin"], e = Link.Prepare(
		`SELECT "Name", "Role" FROM "User" WHERE "Login"=$1 AND "Password"=$2`)
	if e != nil {
		fmt.Println(e)
		return e
	}

	return nil
}

func (u *User) Insert() bool {
	stmt, ok := query["UserInsert"]
	if !ok {
		return false
	}
	_, e := stmt.Exec(u.Login, u.Password, u.Name, u.Role)
	if e != nil {
		Logger.Println(e)
		return false
	}

	return true
}

func (u *User) LogIn() bool {
	stmt, ok := query["UserLogin"]
	if !ok {
		return false
	}
	row := stmt.QueryRow(u.Login, u.Password)
	e := row.Scan(&u.Name, &u.Role)
	if e != nil {
		Logger.Println(e)
		return false
	}

	return true
}

func CreateSession(user *User, m map[string]Session) (string, bool) {
	stmt, ok := query["SessionInsert"]
	if !ok {
		Logger.Println("SessionInsert not found")
		return "", false
	}

	hash, e := utils.GenerateHash(user.Login)
	if e != nil {
		Logger.Println(e)
		return "", false
	}

	_, e = stmt.Exec(hash, user.Login, "")
	if e != nil {
		Logger.Println(e)
		return "", false
	}

	if m != nil {
		m[hash] = Session{
			Hash: hash,
			User: User{
				Login:    user.Login,
				Password: "",
				Name:     user.Name,
				Role:     user.Role,
			},
			Date:   time.Now().String()[:19],
			Device: "",
		}
	}

	return hash, true
}

func LoadSession(m map[string]Session) {
	stmt, ok := query["SessionSelect"]
	if !ok {
		Logger.Println("SessionSelect not found")
		return
	}

	rows, e := stmt.Query()
	if e != nil {
		Logger.Println(e)
		return
	}

	defer rows.Close()

	for rows.Next() {
		session := Session{}
		e = rows.Scan(&session.Hash, &session.User.Login, &session.User.Name, &session.User.Role, &session.Date, &session.Device)
		if e != nil {
			Logger.Println(e)
			return
		}

		m[session.Hash] = session
	}
}

func (u *User) SelectAll() []User {
	rows, e := Link.Query(`SELECT "Login", "Name" FROM "User" ORDER BY "Name"`)
	if e != nil {
		fmt.Println(e)
		return nil
	}

	defer rows.Close()

	users := make([]User, 0)

	for rows.Next() {
		e = rows.Scan(&u.Login, &u.Name)
		if e != nil {
			fmt.Println(e)
			return nil
		}

		users = append(users, User{
			Login:    u.Login,
			Password: "",
			Name:     u.Name,
		})
	}

	return users
}

func deferTx(tx *sql.Tx) {
	r := recover()
	if r != nil {
		_ = tx.Rollback()
	} else {
		_ = tx.Commit()
	}
}
