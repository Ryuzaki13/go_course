package database

import (
	"awesomeProject2/utils"
	"database/sql"
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

var sessionMap map[string]Session
var query map[string]*sql.Stmt

func prepareUser() []string {
	query = make(map[string]*sql.Stmt)
	sessionMap = make(map[string]Session)

	errorList := make([]string, 0)

	var e error
	query["SessionInsert"], e = Link.Prepare(`INSERT INTO "Session" ("Hash", "User", "Date", "Device") VALUES ($1, $2, CURRENT_TIMESTAMP, $3)`)
	if e != nil {
		errorList = append(errorList, e.Error())
	}

	query["SessionSelect"], e = Link.Prepare(`SELECT "Hash", "Login", "Name", "Role", "Date", "Device" FROM "Session" AS s INNER JOIN "User" AS u ON u."Login"=s."User"`)
	if e != nil {
		errorList = append(errorList, e.Error())
	}

	query["UserInsert"], e = Link.Prepare(`INSERT INTO "User" ("Login", "Password", "Name", "Role") VALUES ($1, $2, $3, $4)`)
	if e != nil {
		errorList = append(errorList, e.Error())
	}

	query["UserSelect"], e = Link.Prepare(`SELECT "Login", "Password", "Name" FROM "User" ORDER BY "Name"`)
	if e != nil {
		errorList = append(errorList, e.Error())
	}

	query["UserLogin"], e = Link.Prepare(`SELECT "Name", "Role" FROM "User" WHERE "Login"=$1 AND "Password"=$2`)
	if e != nil {
		errorList = append(errorList, e.Error())
	}

	return errorList
}

func (u *User) Insert() bool {
	stmt, ok := query["UserInsert"]
	if !ok {
		return false
	}
	_, e := stmt.Exec(u.Login, u.Password, u.Name, u.Role)
	if e != nil {
		utils.Logger.Println(e)
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
		utils.Logger.Println(e)
		return false
	}
	return true
}

func CreateSession(user *User) (string, bool) {
	stmt, ok := query["SessionInsert"]
	if !ok {
		utils.Logger.Println("SessionInsert not found")
		return "", false
	}

	hash, e := utils.GenerateHash(user.Login)
	if e != nil {
		utils.Logger.Println(e)
		return "", false
	}

	_, e = stmt.Exec(hash, user.Login, "")
	if e != nil {
		utils.Logger.Println(e)
		return "", false
	}

	if sessionMap != nil {
		sessionMap[hash] = Session{
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
		utils.Logger.Println("SessionSelect not found")
		return
	}

	rows, e := stmt.Query()
	if e != nil {
		utils.Logger.Println(e)
		return
	}

	defer rows.Close()

	for rows.Next() {
		session := Session{}
		e = rows.Scan(&session.Hash, &session.User.Login, &session.User.Name, &session.User.Role, &session.Date, &session.Device)
		if e != nil {
			utils.Logger.Println(e)
			return
		}

		m[session.Hash] = session
	}
}

func GetSession(hash string) *Session {
	session, ok := sessionMap[hash]
	if ok {
		return &session
	}

	return nil
}
