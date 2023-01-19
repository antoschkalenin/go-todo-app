package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	auth "github.com/korylprince/go-ad-auth/v3"
	"github.com/zhashkevych/todo/model"
	"github.com/zhashkevych/todo/pkg/repository"
	"os"
	"strconv"
	"time"
)

const (
	// соль для пароля со случайными символами
	salt = "sajdiu5asdyasf"
	// ключ для расшифровки токена
	signingKey = "asdjiOAjf93fsd@$1!sfQQQ"
	tokenTTL   = 12 * time.Hour
)

// расширяем стандартный ответ StandardClaims своим tokenClaims с id пользователем
type tokenClaims struct {
	jwt.StandardClaims
	UserId string `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

// ParseToken метод парсит токен и возвращает id пользователя
// ParseWithClaims из библиотеки jwt принимает 3 параметра:
// 1) токен для парсинга
// 2) структуру нашего token claims который мы используем для ответа
// 3) функцию которая возвращает ключ подпись или ошибку, в этой функции проверяем метод подписи токена
func (s *AuthService) ParseToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(accessToken *jwt.Token) (interface{}, error) {
		// проверяем метод подписи токена если это не HMAC
		// если все ок то возаращаем ключ подписи иначе ошибку
		if _, ok := accessToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неверный метод подписи")
		}
		return []byte(signingKey), nil
	})
	// проверяем на наличие ошибок при вызове метода ParseWithClaims
	if err != nil {
		return "", err
	}
	// функция ParseWithClaims возвращает объект токена в котором
	// есть Claims типа interface,
	// приводим полученный токен к структуре tokenClaims
	claims, ok := token.Claims.(*tokenClaims)
	// проверяем результат приведения к структуре Claims
	if !ok {
		return "", errors.New("claims не относится к типу Claims")
	}
	return claims.UserId, nil
}

// NewAuthService конструктор
// иницилизируем в конструкторе репозиторий
func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

// CreateUser создае пользователя
func (s *AuthService) CreateUser(user model.User) (int, error) {
	// перед записью в БД хэшируем пароль
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

// AuthTokenDB входим под пользователем из БД
func (s *AuthService) AuthTokenDB(username, password string) (string, error) {
	// получаем пользователя из БД
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}
	return generateToken(strconv.Itoa(user.Id))
}

// AuthTokenAD входим под пользователем из AD
func (s *AuthService) AuthTokenAD(username, password string) (string, error) {
	port, _ := strconv.Atoi(os.Getenv("AD_PORT"))
	config := &auth.Config{
		Server:   os.Getenv("AD_HOST"),
		Port:     port,
		BaseDN:   os.Getenv("AD_BASE_DN"),
		Security: auth.SecurityStartTLS,
	}
	status, err := auth.Authenticate(config, username, password)
	if err != nil {
		return "", err
	}
	if !status {
		return "", err
	}
	return generateToken(username)
}

func generateToken(userId string) (string, error) {
	// если такой пользователь существует то сегенрируем токен JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(), //  установим на 12ч больше текущего времени
			IssuedAt:  time.Now().Unix(),               // время когда токен был сгенерирован
		},
		userId,
	})
	return token.SignedString([]byte(signingKey))
}

// функция хэширования пароля
func generatePasswordHash(password string) string {
	// используем алгоритм sha1
	hash := sha1.New()
	hash.Write([]byte(password))
	// вовзращаем строку с хэшом примешав к паролю соль
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
