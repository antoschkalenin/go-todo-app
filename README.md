# go-todo-app
REST API пет проект на GO с использованием gin-gonic

Команда инициализирует миграцию (создает папку schema и 2 sql файла для миграций)<br />
 migrate create -ext sql -dir ./schema -seq init
Детальный разбор команды:
migrate create  - создание миграции через migrate утилиту
-ext sql  - расширение файлов указываем sql
-dir ./schema - создаем папку куда будем сохранять миграции
-seq init - название файлов миграции при создании

Запуск миграции
migrate -path ./schema -database 'postgres://postgres:postgres@localhost:5432/go?sslmode=disable' up
Откат
migrate -path ./schema -database 'postgres://postgres:postgres@localhost:5432/go?sslmode=disable' down
