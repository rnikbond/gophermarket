# Накопительная система лояльности «Гофермаркет»

В проекте реализовано HTTP API для регистрации и аутентификации пользователей, осуществления заказов и взаимодействия с внешней системой лояльности.

## Конфигурация
Настройка конфигурации происходит через переменные окруженис ОС или через аргументы командной строки.\
Для просмотра переменных окружения используйте параметр **-h**.\
Для обработки аргументов командной строки используется пакет
```
flag
```
Для считывания переменных окружения используется пакет
```
github.com/caarlos0/env
```

## Авторизация
Для авторизации ипользуется JWT-токен, который помещается в Cookie.\
В JWT-токен помещается login пользователя.\
Проверка авторизации для запросов, доступных только авторизованным пользователям, осуществляется в middleware.
Для работы с JWT токеном используется пакет:
```
github.com/dgrijalva/jwt-go
```

## СУБД
PostgrSQL.
Файлы миграций находятся в 
```
/db/migrations/*
```
Работа с миграциями реализована с использованием пакета
```
github.com/golang-migrate/migrate
```

## Хранение данных регистрации
Вместо пароля, в базу сохраняется хеш пароля с добавлением *соли*. 

## Unit-тесты
Для теситрования используются mock-файлы.\
Используемые пакеты:
```
github.com/golang/mock/gomock
github.com/stretchr/testify
```

## Docker
Конфигурация образа находится в *Dockerfile*.\
Для запуска контейнера используется docker-compose, конфигурация находится в *docker-compose.yml*.

# Дальнейшее развитие проекта
- Добавить кеширование данных
  - Кеширование ID пользователя, чтоб не ходить за ним каждый раз в базу
  - Кеширование списка заказов
- Пагинация при загрузке списков заказов из базы
