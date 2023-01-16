package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/zhashkevych/todo"
	"github.com/zhashkevych/todo/pkg/repository"
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
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

// ParseToken метод парсит токен и возвращает id пользователя
// ParseWithClaims из библиотеки jwt принимает 3 параметра:
// 1) токен для парсинга
// 2) структуру нашего token claims который мы используем для ответа
// 3) функцию которая возвращает ключ подпись или ошибку, в этой функции проверяем метод подписи токена
func (s *AuthService) ParseToken(accessToken string) (int, error) {
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
		return 0, err
	}
	// функция ParseWithClaims возвращает объект токена в котором
	// есть Claims типа interface,
	// приводим полученный токен к структуре tokenClaims
	claims, ok := token.Claims.(*tokenClaims)
	// проверяем результат приведения к структуре Claims
	if !ok {
		return 0, errors.New("claims не относится к типу Claims")
	}
	return claims.UserId, nil
}

// NewAuthService конструктор
// иницилизируем в конструкторе репозиторий
func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

// CreateUser создае пользователя
func (s *AuthService) CreateUser(user todo.User) (int, error) {
	// перед записью в БД хэшируем пароль
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

// GenerateToken генерирует JWT токен и возвращает его
func (s *AuthService) GenerateToken(username, password string) (string, error) {
	// получаем пользователя из БД
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}
	// если такой пользователь существует то сегенрируем токен JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(), //  установим на 12ч больше текущего времени
			IssuedAt:  time.Now().Unix(),               // время когда токен был сгенерирован
		},
		user.Id,
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
