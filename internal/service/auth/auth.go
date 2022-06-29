package auth

import (
	"fmt"
	"log"
	"time"

	"social-network/internal/errapp"
	"social-network/internal/models"
	"social-network/internal/password"

	"github.com/golang-jwt/jwt"

	"github.com/google/uuid"
)

type SignInData struct {
	Login    string `json:"login" binding:"required,min=3,max=20"`
	Password string `json:"password,required" binding:"required,min=8,max=50"`
}

type AccessData struct {
	UserID uuid.UUID `json:"id"`
	Token  string    `json:"token,required"`
}

type Service interface {
	SignIn(data SignInData) (accessData AccessData, err error)
	CheckToken(token string) (allow bool, err error)
}

type Storage interface {
	CheckAccessExist(login string, password string) (accessData models.UserAccess, exist bool, err error)
}

type service struct {
	storage        Storage
	hashSalt       string
	signingKey     []byte
	expireDuration time.Duration
}

func NewService(store Storage) Service {
	return &service{store, "", []byte{}, 1 * time.Hour}
	// viper.GetString("auth.hash_salt"),
	// 	[]byte(viper.GetString("auth.signing_key")),
	// 	viper.GetDuration("auth.token_ttl")*time.Second,
}

func (s *service) SignIn(data SignInData) (accessData AccessData, err error) {
	encryptedPassword := password.Encrypt(data.Password)
	accessDataModel, exist, err := s.storage.CheckAccessExist(data.Login, encryptedPassword)
	if err != nil {
		return AccessData{}, err
	}

	if !exist {
		return AccessData{}, fmt.Errorf("%w", errapp.AccessDataNotFound)
	}

	token, err := s.generateToken(data.Login)
	if err != nil {
		return AccessData{}, err
	}

	log.Println("token", token)
	return AccessData{
		UserID: accessDataModel.UserID,
		Token:  token,
	}, nil
}

func (s *service) CheckToken(token string) (allow bool, err error) {
	allow, err = s.parseToken(token)
	if err != nil {
		return false, err
	}

	if !allow {
		return false, nil
	}

	return true, nil
}

func (s *service) generateToken(login string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.expireDuration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Username: login,
	})

	tokenStr, err := token.SignedString(s.signingKey)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func (s *service) parseToken(accessToken string) (allow bool, err error) {
	token, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected metoho %v", token.Header["alg"])
		}
		return s.signingKey, nil
	})
	if err != nil {
		return false, err
	}

	if _, ok := token.Claims.(*Claims); ok && token.Valid {
		return true, nil
	}
	return false, err
}

type Claims struct {
	jwt.StandardClaims
	Username string `json:"username"`
}
