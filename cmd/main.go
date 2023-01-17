package main

import (
	"context"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/zhashkevych/todo"
	"github.com/zhashkevych/todo/pkg/controller"
	"github.com/zhashkevych/todo/pkg/repository"
	"github.com/zhashkevych/todo/pkg/service"
	"os"
	"os/signal"
	"syscall"
)

// основной метод запуска приложения
func main() {
	// устанавливаем формат JSON для логера logrus для систем сбора логов
	logrus.SetFormatter(new(logrus.JSONFormatter))

	// иницилизируем конфиг файл
	if err := initConfig(); err != nil {
		logrus.Fatalf("ошибка иницилизации конфиг файла: %s", err.Error())
	}

	// загружаем переменные окружения из файла .env с помощью метода Load
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("ошибка загрузки переменных окружения: %s", err.Error())
	}

	// иницилизируем repository.Config настройками БД из .env файла
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_DBNAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	})

	if err != nil {
		logrus.Fatalf("Ошибка иницилизации БД: %s", err.Error())
	}

	// dependency injection (внутренний слой не зависит от внешнего, внешний зависит от внутреннего).
	// https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
	// Наши обработчики (handlers) будут вызывать сервисы поэтому
	// добавим в структуру указатели на сервис.
	// Создадим в нужном порядке наши зависимости:
	// 1)репо 2)сервисы 3)обработчики
	repo := repository.NewRepository(db)
	services := service.NewService(repo)
	handlers := controller.NewRouters(services)

	srv := new(todo.Server)
	// запускаем сервер
	// добавляем плавное заверщение работы приложения а именно:
	// - завершение работы http запросов
	// - завершение работы SQL запросов к БД
	// для этого запустим наш сервис в goroutine (синтаксис go func())
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("ошибка запуска приложения: %s", err.Error())
		}
	}()
	logrus.Println("Старт приложения")
	// что бы избежать выхода из метода main и завершения
	// приложения добавим блокировку функции main с помощью канала os.Signal
	quit := make(chan os.Signal, 1)
	// запись в канал будет происходить когда процесс в котором выполняется наше приложение
	// получит сигнал от системы типа syscall.SIGTERM или syscall.SIGINT, это сигналы в UNIX системах
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	// строка для чтения из канала которая будет блокировать выполнение главной goroutine main
	<-quit

	// после выхода из приложения плавно завершим работу, это гарантирует нам что мы закончим
	// выполнение всех текущий операций перед выходом из приложения
	logrus.Println("Остановка приложения")
	// остановка сервера
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("Ошибка остановки сервера: %s", err.Error())
	}
	// закрытие соедниения с БД
	if err := db.Close(); err != nil {
		logrus.Errorf("Ошибка закрытия соедниения с БД: %s", err.Error())
	}
}

// метод описывает иницилизацию конфиг файла
func initConfig() error {
	// путь к файлу
	viper.AddConfigPath("configs")
	// название конфиг файла
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
