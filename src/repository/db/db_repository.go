package db

import (
	"strings"

	mysql_db "github.com/LibenHailu/bookstore_oauth_api/src/client/mysql"
	"github.com/LibenHailu/bookstore_oauth_api/src/domain/access_token"
	"github.com/LibenHailu/bookstore_oauth_api/src/utils/errors"
)

const (
	queryGetAccessToken    = "SELECT access_token,user_id,client_id,expires FROM	access_tokens WHERE access_token = ?;"
	queryCreateAccessToken = "INSERT INTO access_tokens(access_token,user_id,client_id,expires) VALUES (?,?,?,?);"
	queryUpdateExpires     = "UPDATE access_token SET expires=? WHERE access_token=?;"
	errorNoRows            = "no rows in result set"
)

type DbRepository interface {
	GetById(string) (*access_token.AccessToken, *errors.RestErr)
	Create(access_token.AccessToken) *errors.RestErr
	UpdateExpirationTime(access_token.AccessToken) *errors.RestErr
}

type dbRepository struct {
}

func NewRepository() DbRepository {
	return &dbRepository{}
}

func (r *dbRepository) GetById(id string) (*access_token.AccessToken, *errors.RestErr) {
	stmt, err := mysql_db.Client.Prepare(queryGetAccessToken)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	var result access_token.AccessToken
	if getErr := stmt.QueryRow(id).Scan(&result.AccessToken, &result.UserId, &result.ClientId, &result.Expires); getErr != nil {
		if strings.Contains(getErr.Error(), errorNoRows) {
			return nil, errors.NewNotFoundError("no access token found with the given id")
		}
		return nil, errors.NewInternalServerError(getErr.Error())
	}
	return &result, nil

}

func (r *dbRepository) Create(at access_token.AccessToken) *errors.RestErr {
	stmt, err := mysql_db.Client.Prepare(queryCreateAccessToken)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	_, saveErr := stmt.Exec(at.AccessToken, at.UserId, at.ClientId, at.Expires)

	if saveErr != nil {
		return errors.NewInternalServerError(saveErr.Error())
	}
	return nil
}

func (r *dbRepository) UpdateExpirationTime(at access_token.AccessToken) *errors.RestErr {
	stmt, err := mysql_db.Client.Prepare(queryUpdateExpires)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(at.Expires, at.AccessToken)

	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	return nil
}
