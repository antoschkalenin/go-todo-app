# go-todo-app
REST API пет проект на GO с использованием gin-gonic.

<p>В проекте используется БД postgres.</p>

Для подключения к БД postgres перед стартом приложения необходимо добавить коннекты:<br />
1. Добавить в корень проекта файл .env<br /><br />
2. добавить в файл .env переменные:<br />
DB_HOST=хост БД<br />
DB_PORT=порт БД<br />
DB_USERNAME=uername БД<br />
DB_PASSWORD=пароль БД<br />
DB_DBNAME=название БД<br />
DB_SSLMODE=disable<br /><br />
Пример:<br />
DB_HOST=localhost<br />
DB_PORT=5432<br />
DB_USERNAME=postgres<br />
DB_PASSWORD=postgres<br />
DB_DBNAME=go<br />
DB_SSLMODE=disable<br /><br />
3. Необходимо выполнить миграцию в БД.<br />
<b>Запуск миграции:</b><br />
migrate -path ./schema -database 'postgres://postgres:postgres@localhost:5432/go?sslmode=disable' up<br />
<b>Откат:</b><br />
migrate -path ./schema -database 'postgres://postgres:postgres@localhost:5432/go?sslmode=disable' down
4. Запустить проект через консоль: <b>go run cmd/main.go</b>

В проекте по пути configs/config.yml установлен порт приложения по умолчанию 8001.

Так же в проекте есть папка postman для тестирования приложения