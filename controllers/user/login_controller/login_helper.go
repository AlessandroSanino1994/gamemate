package loginController

import (
	"math/rand"

	"sanino/gamemate/configurations"
	"sanino/gamemate/constants"
	"sanino/gamemate/controllers/shared"
	"sanino/gamemate/models/user/requests/login"

	"strconv"
)

//insertIntoArchives (without check of previous insertions, only error reporting)
//inserts a new User into the archives, doing the salty & hashy work.
func insertIntoArchives(RegTry loginRequests.UserRegistration) (int64, error) {
	salt := rand.Intn(constants.MAX_NUMBER_SALT)
	saltedPass := RegTry.Password + strconv.Itoa(salt)
	hash_pwd := controllerSharedFuncs.ConvertToHexString(saltedPass)
	stmtQuery, err := configurations.ArchivesPool.Prepare("INSERT INTO users (userID, username, hash_pwd, hash_salt) VALUES (NULL, ?, UNHEX(?), ?)")
	if err != nil {
		return -1, err
	}
	defer stmtQuery.Close()

	result, err := stmtQuery.Exec(RegTry.Username, hash_pwd, salt)
	if err != nil { //did not exec query (syntax)
		return -1, err
	}

	affected_rows, err := result.RowsAffected()
	if affected_rows <= 0 || err != nil { //Did not insert
		return -1, err
	}

	insertID, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return insertID, nil
}

//checkLogin checks if a user pass combination is valid for the specified auth try.
//
//Returns true if login is valid, false otherwise and report errors.
func checkLogin(AuthTry loginRequests.UserAuth) (bool, int64, error) {
	var password_hash string
	var userID int64
	var salt int

	stmtQuery, err := configurations.ArchivesPool.Prepare("SELECT HEX(hash_pwd), hash_salt, userID FROM users WHERE username = ?")
	if err != nil {
		return false, -1, err
	}
	defer stmtQuery.Close()

	result, err := stmtQuery.Query(AuthTry.Username)
	if err != nil {
		return false, -1, err
	}
	defer result.Close()

	if !result.Next() { //not found
		return false, -1, nil
	}

	err = result.Scan(&password_hash, &salt, &userID)
	if err != nil {
		return false, -1, err
	}
	salted_pwd := AuthTry.Password + strconv.Itoa(salt)
	salted_hash := controllerSharedFuncs.ConvertToHexString(salted_pwd)
	//fmt.Println("0x" + password_hash)
	//fmt.Println(salted_hash)
	if salted_hash == password_hash {
		return true, userID, nil
	}
	return false, -1, nil
}

//userInArchives checks into the archives if a user is already registered.
//username is valid for both username and email auth.
//
//Returns true if the user is already in the archives, false otherwise.
//NOTE: should i return cachable values??? or get it with another request?
func isRegistered(username string) (bool, error) {
	stmtQuery, err := configurations.ArchivesPool.Prepare("SELECT COUNT(*) AS num_rows FROM users WHERE username = ?")
	if err != nil { //cannot check, consider error
		return true, err
	}
	defer stmtQuery.Close()

	rows, err := stmtQuery.Query(username)
	if err != nil { //cannot check, consider error
		return true, err
	}
	defer stmtQuery.Close()

	var num_rows int
	rows.Scan(&num_rows)
	return num_rows > 0, nil
}

func updateCacheNewUserSession(userID int64, username string) (string, error) {
	return controllerSharedFuncs.UpdateCacheNewSession(constants.LOGGED_USERS_SET, constants.CACHE_REFRESH_INTERVAL, userID, "username", username)
}
