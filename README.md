# go-todo-app
REST API пет проект на GO с использованием gin-gonic.

<p>Перед стартом приложения необходимо выполнить миграцию в БД.</p>
<p>В проекте используется postgres.</p>

<b>Запуск миграции:</b><br />
migrate -path ./schema -database 'postgres://postgres:postgres@localhost:5432/go?sslmode=disable' up<br />
<b>Откат:</b><br />
migrate -path ./schema -database 'postgres://postgres:postgres@localhost:5432/go?sslmode=disable' down


