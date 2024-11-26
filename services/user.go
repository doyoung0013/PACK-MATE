package services

import (
	"fmt"

	"github.com/GDG-on-Campus-KHU/SC4_BE/db"
	"github.com/GDG-on-Campus-KHU/SC4_BE/models"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct{}

func HashPassword(password string) (string, error) {
	//GenerateFromPassword : 두번째 인자는 cost값, 높을수록 보안 좋은대신 느려짐. 기본값 설정시 10
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *UserService) LoginUser(l *models.LoginData) (*models.User, error) {
	var user models.User
	var password string
	err := db.DB.QueryRow(`
         SELECT u.*, ur.rest_id
        FROM users u
        LEFT JOIN gdgdb.user_rest ur on u.id = ur.user_id
        WHERE email = ?
		`, l.Email).Scan(
		&user.ID,
		&user.Email,
		&password,
		&user.PhoneNum,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.RestID,
	)
	if err != nil {
		return nil, fmt.Errorf("존재하지 않는 회원입니다.")
	}
	// 비밀번호 검증
	if !CheckPasswordHash(l.Password, password) {
		return nil, fmt.Errorf("invalid password")
	}

	return &user, nil
}

func (s *UserService) CreateUser(u *models.User) error {
	var existingID int64
	err := db.DB.QueryRow("SELECT id FROM users WHERE email = ?", u.Email).Scan(&existingID)
	if err == nil || existingID != 0 {
		return fmt.Errorf("email already exists")
	}

	hashedPassword, err := HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword

	result, err := db.DB.Exec(`
        INSERT INTO users (email, password, phone_num, name)  
        VALUES (?, ?, ?, ?)`,
		u.Email, u.Password, u.PhoneNum, u.Name)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	u.ID = id

	err = db.DB.QueryRow("SELECT id, email, name, phone_num,created_at, updated_at FROM users WHERE id = ?", id).
		Scan(
			&u.ID,
			&u.Email,
			&u.Name,
			&u.PhoneNum,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) GetUser(id int64) (*models.User, error) {
	var user models.User
	err := db.DB.QueryRow("SELECT id, email, name, phone_num,created_at, updated_at FROM users WHERE id = ?", id).
		Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.PhoneNum,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UpdateUser(u *models.User) error {
	_, err := db.DB.Exec(`
		UPDATE users 
		SET name = ?, email = ?, phone_num = ?
		WHERE id = ?`,
		u.Name, u.Email, u.PhoneNum, u.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) DeleteUser(id int64) error {
	_, err := db.DB.Exec(`
		DELETE FROM users 
		WHERE id = ?`,
		id)
	if err != nil {
		return err
	}
	return nil
}
