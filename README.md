# go-todo-app
REST API пет проект на GO с использованием gin-gonic

<b>Команда инициализирует миграцию (создает папку schema и 2 sql файла для миграций)</b><br />
 migrate create -ext sql -dir ./schema -seq init<br />
Детальный разбор команды:<br />
migrate create  - создание миграции через migrate утилиту<br />
-ext sql  - расширение файлов указываем sql<br />
-dir ./schema - создаем папку куда будем сохранять миграции<br />
-seq init - название файлов миграции при создании<br />

<b>Запуск миграции</b><br />
migrate -path ./schema -database 'postgres://postgres:postgres@localhost:5432/go?sslmode=disable' up<br />
<b>Откат</b><br />
migrate -path ./schema -database 'postgres://postgres:postgres@localhost:5432/go?sslmode=disable' down
