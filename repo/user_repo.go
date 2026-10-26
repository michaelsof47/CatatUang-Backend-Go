package repo

import (
	"catatuangbackend/models"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type UserInterface interface {
	CheckAvailabilityAccount(createUser models.User) (bool, error)
	InsertNewAccount(createUser models.User, savePath string, hashedPassword []byte) (models.UserResult, error)
	GetUserById(id string) (models.UserResult, error)
	GetUserByPhoneAndEmail(emailOrPhone string) (models.UserResult, error)
	GetUserPasswordById(id string) (models.User, error)
	UpdatePasswordById(hashedPassword []byte, id string) error
	UpdatePasswordBasedOnEmailAndPhone(hashedPassword []byte, emailOrPhone string) error
	UpdateAccountUserInfo(firstName, lastName, email, phone, id string) error
	GetUrlImagePhoto(id string) (string, error)
	UpdatePhotoProfile(savePath, id string) error
}

type UserRepo struct {
	db *sql.DB
}

func UserRepository(db *sql.DB) UserInterface {
	return &UserRepo{db: db}
}

func (h *UserRepo) CheckAvailabilityAccount(createUser models.User) (bool, error) {
	var exists bool
	isAccAvailableQuery := `SELECT EXISTS(
		SELECT * FROM "User" WHERE 
		LOWER(first_name) = LOWER($1) 
		AND LOWER(last_name) = LOWER($2) 
		AND LOWER(email) = LOWER($3)
		AND LOWER(phone) = LOWER($4)
		)`
	err := h.db.QueryRow(isAccAvailableQuery, &createUser.FirstName,
		&createUser.LastName, &createUser.Email, &createUser.Phone).Scan(&exists)
	return exists, err
}

func (h *UserRepo) InsertNewAccount(createUser models.User, savePath string, hashedPassword []byte) (models.UserResult, error) {
	var user models.UserResult
	query := `INSERT INTO "User" (first_name, last_name, url_user_image, email, phone, password) 
			  VALUES ($1, $2, $3, $4, $5, $6)`
	if _, err := h.db.Exec(query, createUser.FirstName, createUser.LastName, savePath, createUser.Email, createUser.Phone, hashedPassword); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return models.UserResult{}, errors.New("Email anda sudah terdaftar")
		}
		return models.UserResult{}, err
	}
	return user, nil
}

func (h *UserRepo) GetUserById(id string) (models.UserResult, error) {
	var user models.UserResult
	query := `SELECT id, first_name, last_name, reward_status, url_user_image, email, phone, password, "createdAt" FROM "User" WHERE id = $1`
	err := h.db.QueryRow(query, id).Scan(&user.Id, &user.FirstName, &user.LastName, &user.RewardStatus, &user.UrlUserImage, &user.Email, &user.Phone, &user.Password, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return models.UserResult{}, errors.New("User tidak ditemukan")
	} else if err != nil {
		return models.UserResult{}, err
	}
	return user, err
}

func (h *UserRepo) GetUserByPhoneAndEmail(emailOrPhone string) (models.UserResult, error) {
	var user models.UserResult
	query := `SELECT id, password FROM "User" WHERE email = $1 OR phone = $1`
	err := h.db.QueryRow(query, emailOrPhone).Scan(&user.Id, &user.Password)
	return user, err
}

func (h *UserRepo) GetUserPasswordById(id string) (models.User, error) {
	var user models.User
	query := `SELECT password FROM "User" WHERE id = $1`
	err := h.db.QueryRow(query, id).Scan(&user.Password)
	if err == sql.ErrNoRows {
		return models.User{}, errors.New("Akun Password tidak ditemukan")
	} else if err != nil {
		return models.User{}, err
	}
	return user, err
}

func (h *UserRepo) UpdatePasswordById(hashedPassword []byte, id string) error {
	query := `UPDATE "User" SET password = $1 WHERE id = $2`
	_, err := h.db.Exec(query, hashedPassword, id)
	return err
}

func (h *UserRepo) UpdatePasswordBasedOnEmailAndPhone(hashedPassword []byte, emailOrPhone string) error {
	query := `UPDATE "User" SET password = $1 WHERE email = $2 OR phone = $2`
	_, err := h.db.Exec(query, hashedPassword, emailOrPhone)
	return err
}

func (h *UserRepo) UpdateAccountUserInfo(firstName, lastName, email, phone, id string) error {
	query := `UPDATE "User" SET first_name = $1, last_name = $2, email = $3, phone = $4 WHERE id = $5`
	_, err := h.db.Exec(query, firstName, lastName, email, phone, id)
	return err
}

func (h *UserRepo) GetUrlImagePhoto(id string) (string, error) {
	var oldImagePath string
	query := `SELECT url_user_image FROM "User" WHERE id = $1`
	err := h.db.QueryRow(query, id).Scan(&oldImagePath)
	if err == sql.ErrNoRows {
		return "", errors.New("Foto Akun tidak ditemukan")
	} else if err != nil {
		return "", err
	}
	return oldImagePath, err
}

func (h *UserRepo) UpdatePhotoProfile(savePath, id string) error {
	query := `UPDATE "User" SET url_user_image = $1 WHERE id = $2`
	_, err := h.db.Exec(query, savePath, id)
	return err
}
