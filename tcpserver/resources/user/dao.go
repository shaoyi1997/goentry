package user

import (
	"database/sql"

	"git.garena.com/shaoyihong/go-entry-task/tcpserver/services"

	"git.garena.com/shaoyihong/go-entry-task/common/logger"

	"git.garena.com/shaoyihong/go-entry-task/common/pb"
)

type IUserDAO interface {
	GetByUsername(username string) (*pb.User, error)
	UpdateNickname(username, nickname string) error
	UpdateProfileImage(username, imageUrl string) error
}

type UserDAO struct {
	db        *sql.DB
	tableName string
}

func NewUserDAO(database services.Database) IUserDAO {
	createUserTable(database)
	return &UserDAO{
		db:        database.Db,
		tableName: database.DatabaseName + ".User",
	}
}

func createUserTable(database services.Database) {
	query := `CREATE TABLE IF NOT EXISTS ` + database.DatabaseName + ".User (" +
		`user_id INT UNSIGNED AUTO_INCREMENT,
		username VARCHAR(64) NOT NULL DEFAULT '',
		password VARCHAR(64) NOT NULL DEFAULT '',	
		nickname VARCHAR(64) NOT NULL DEFAULT '',
		profile_image VARCHAR(128),
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY(user_id),
		UNIQUE KEY(username)
	)ENGINE=InnoDB DEFAULT CHARSET=utf8;`
	if _, err := database.Db.Exec(query); err != nil {
		logger.ErrorLogger.Panicln("Failed to create user table")
	}
}

func (dao *UserDAO) GetByUsername(username string) (*pb.User, error) {
	query := "SELECT user_id, username, nickname, password, profile_image FROM " + dao.tableName + " WHERE username = ?"
	user := new(pb.User)
	err := dao.db.QueryRow(query, username).Scan(&user.UserId, &user.Username, &user.Nickname, &user.Password, &user.ProfileImage)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (dao *UserDAO) UpdateNickname(username, nickname string) error {
	_, err := dao.db.Exec("UPDATE "+dao.tableName+" SET nickname = ? WHERE username = ?", nickname, username)
	if err != nil {
		logger.ErrorLogger.Println("Failed to update user nickname: ", err)
		return err
	}
	return nil
}

func (dao *UserDAO) UpdateProfileImage(username, imageUrl string) error {
	_, err := dao.db.Exec("UPDATE"+dao.tableName+"SET profile_image = ? WHERE username = ?", imageUrl, username)
	if err != nil {
		logger.ErrorLogger.Println("Failed to update user profile image: ", err)
		return err
	}
	return nil
}
