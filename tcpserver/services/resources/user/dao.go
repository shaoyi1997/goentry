package user

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"git.garena.com/shaoyihong/go-entry-task/common/logger"
	"git.garena.com/shaoyihong/go-entry-task/common/pb"
	"git.garena.com/shaoyihong/go-entry-task/tcpserver/common"
)

const (
	numOfTables = 10
)

type IUserDAO interface {
	getByUsername(username string) (*pb.User, error)
	updateNickname(username, nickname string) error
	updateProfileImage(username, imageURL string) error
	insert(username, password, nickname, imageURL string) (*pb.User, error)
}

type DAO struct {
	db        *sql.DB
	tableName string
}

func newUserDAO(database common.Database) IUserDAO {
	createUserTable(database)

	return &DAO{
		db:        database.DB,
		tableName: database.DatabaseName + ".User",
	}
}

func createUserTable(database common.Database) {
	for i := 0; i < numOfTables; i++ {
		query := `CREATE TABLE IF NOT EXISTS ` + database.DatabaseName + ".User_" + strconv.Itoa(i) +
			`(user_id INT UNSIGNED AUTO_INCREMENT,
				username VARCHAR(64) NOT NULL,
				password VARCHAR(64) NOT NULL,	
				nickname VARCHAR(64) NOT NULL DEFAULT '',
				profile_image VARCHAR(128) NOT NULL DEFAULT '',
				created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				PRIMARY KEY(user_id),
				UNIQUE KEY(username)
			)ENGINE=InnoDB DEFAULT CHARSET=utf8;`
		if _, err := database.DB.Exec(query); err != nil {
			logger.ErrorLogger.Panicln("Failed to create user table")
		}
	}
}

func (dao *DAO) getTableNameByUsername(username string) string {
	return dao.tableName + "_" + strconv.Itoa(int(username[0]%numOfTables))
}

func (dao *DAO) getByUsername(username string) (*pb.User, error) {
	query := "SELECT user_id, username, nickname, password, profile_image FROM " +
		dao.getTableNameByUsername(username) + " WHERE username = ?"
	user := new(pb.User)

	err := dao.db.QueryRow(query, username).Scan(&user.UserId, &user.Username, &user.Nickname,
		&user.Password, &user.ProfileImage)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errUsernameNotFound
		}

		logger.ErrorLogger.Println("Failed to get user:", err)

		return nil, err
	}

	return user, nil
}

func (dao *DAO) updateNickname(username, nickname string) error {
	_, err := dao.db.Exec("UPDATE "+dao.getTableNameByUsername(username)+
		" SET nickname = ? WHERE username = ?", nickname, username)
	if err != nil {
		logger.ErrorLogger.Println("Failed to update user nickname: ", err)

		return err
	}

	return nil
}

func (dao *DAO) updateProfileImage(username, imageURL string) error {
	_, err := dao.db.Exec("UPDATE "+dao.getTableNameByUsername(username)+
		" SET profile_image = ? WHERE username = ?", imageURL, username)
	if err != nil {
		logger.ErrorLogger.Println("Failed to update user profile image: ", err)

		return err
	}

	return nil
}

func (dao *DAO) insert(username, password, nickname, imageURL string) (*pb.User, error) {
	user := &pb.User{
		UserId:       nil,
		Username:     &username,
		Password:     &password,
		Nickname:     &nickname,
		ProfileImage: &imageURL,
	}

	result, err := dao.db.Exec("INSERT INTO "+dao.getTableNameByUsername(username)+
		` (username, password, nickname, profile_image) VALUES (?, ?, ?, ?)`, username, password, nickname, imageURL)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			err = errUsernameAlreadyExists
		}

		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	unsignedID := uint64(id)
	user.UserId = &unsignedID

	return user, nil
}
