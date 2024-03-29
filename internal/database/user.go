package database

import (
	"database/sql"
	"errors"
	"forum/internal/models"
)

type UserRepoImpl struct {
	db *sql.DB
}

func CreateNewUserDB(db *sql.DB) *UserRepoImpl {
	return &UserRepoImpl{db}
}

func (userObj *UserRepoImpl) CreateUserRepo(user *models.User) (int64, error) {
	result, err := userObj.db.Exec(
		`INSERT INTO users (firstName, secondName, usernames, email, password) VALUES (?, ?, ?, ?, ?);`,
		user.FirstName, user.SecondName, user.Username, user.Email, user.Password)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId() // return the auto generated ID of the last added user
}

func (userObj *UserRepoImpl) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := userObj.db.QueryRow(
		`SELECT id, firstName, secondName, usernames, email, password FROM users WHERE email = ?`,
		email).Scan(&user.UserID, &user.FirstName, &user.SecondName, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("element with EMAIL not found")
		}
		return nil, errors.New("error fetching element")
	}

	return user, nil
}

func (userObj *UserRepoImpl) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	err := userObj.db.QueryRow(
		`SELECT id, firstName, secondName, usernames, email, password FROM users WHERE usernames = ?`,
		username).Scan(&user.UserID, &user.FirstName, &user.SecondName, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("element with USERNAME not found")
		}
		return nil, errors.New("error fetching element")
	}
	return user, nil
}

func (userObj *UserRepoImpl) GetUserByUserID(userID int) (*models.User, error) {
	user := &models.User{}
	err := userObj.db.QueryRow(
		`SELECT id, firstName, secondName, usernames, email, password FROM users WHERE id = ?`,
		userID).Scan(&user.UserID, &user.FirstName, &user.SecondName, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (userObj *UserRepoImpl) CreateSession(session *models.Session) error {
	_, err := userObj.db.Exec(
		`INSERT INTO sessions (user_id, token, exp_time) VALUES (?, ?, ?);`,
		session.UserID, session.Token, session.ExpTime)
	if err != nil {
		return err
	}
	return nil
}

func (userObj *UserRepoImpl) UpdateSession(session *models.Session) error {
	if _, err := userObj.db.Exec(
		`UPDATE sessions SET token = ?, exp_time = ? WHERE user_id = ?`,
		session.Token, session.ExpTime, session.UserID); err != nil {
		return err
	}
	return nil
}

func (userObj *UserRepoImpl) GetSessionByUserID(userID int) (*models.Session, error) {
	session := &models.Session{}
	if err := userObj.db.QueryRow(
		`SELECT user_id, token, exp_time FROM sessions WHERE user_id = ?`,
		userID).Scan(&session.UserID, &session.Token, &session.ExpTime); err == nil {
		return nil, err
	}
	return session, nil
}

func (userObj *UserRepoImpl) GetSessionByToken(token string) (*models.Session, error) {
	session := &models.Session{}
	if err := userObj.db.QueryRow(
		`SELECT user_id, token, exp_time FROM sessions WHERE token = ?`,
		token).Scan(&session.UserID, &session.Token, &session.ExpTime); err != nil {
		return nil, err
	}
	return session, nil
}

func (userObj *UserRepoImpl) DeleteSessionByToken(token string) error {
	if _, err := userObj.db.Exec(`
	DELETE FROM sessions WHERE token = ?`, token); err != nil {
		return err
	}
	return nil
}

func (userObj *UserRepoImpl) DeleteSessionByUserID(userID int) error {
	if _, err := userObj.db.Exec(`
	DELETE FROM sessions WHERE user_id = ?`, userID); err != nil {
		return err
	}
	return nil
}
