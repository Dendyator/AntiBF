## Сборка приложения

Из корневого каталога программы необходимо запустить команду `make build`

## Инициализация базы данных

При первом запуске приложения через `make run`, инициализация базы данных PostgreSQL происходит автоматически
благодаря файлу `init-db.sql`.

Если вы хотите изменить структуру базы данных вручную или открыть порты PostgreSQL для внешнего доступа,
выполните следующие шаги:

1. Используйте `docker-compose.yaml`, чтобы указать порты в сервисной секции `db`:
   ```yaml
   ports:
     - "5432:5432"
   ```
2. Подключитесь к базе данных с использованием SQL-клиента.
3. Внесите необходимые изменения в структуру базы данных.


## Swagger

после запуска сервера по адресу http://localhost:8080/swagger/index.html можно потестировать аутентификацию

## Особенности
1). IP допускается только в формате IPv4, в противном случае приложение выдаст ошибку "Invalid IP format"
2). При нахождении IP одновременно в white и black листах приоритет имеет white - авторизация будет одобрена.

## Команды для gRPC

Добавление IP в белый список
grpcurl -plaintext -d '{"subnet": "192.168.1.1."}' localhost:50051 antibruteforce.AntiBruteForce/AddToWhitelist
grpcurl -plaintext -d '{\"subnet\":\"192.168.1.1/25\"}' localhost:50051 antibruteforce.AntiBruteForce/AddToWhitelist (для Powershell)

Удаление IP из белого списка
grpcurl -plaintext -d '{"subnet": "192.168.1.100/32"}' localhost:50051 antibruteforce.AntiBruteForce/RemoveFromWhitelist

Добавление IP в черный список
grpcurl -plaintext -d '{"subnet": "192.168.1.200/32"}' localhost:50051 antibruteforce.AntiBruteForce/AddToBlacklist
grpcurl -plaintext -d '{\"subnet\": \"192.168.1.200/32\"}' localhost:50051 antibruteforce.AntiBruteForce/AddToBlacklist

Удаление IP из черного списка
grpcurl -plaintext -d '{"subnet": "192.168.1.200/32"}' localhost:50051 antibruteforce.AntiBruteForce/RemoveFromBlacklist

Проверка IP в белом списке
grpcurl -plaintext -d '{\"subnet\":\"192.168.1.1\"}' localhost:50051 antibruteforce.AntiBruteForce/CheckWhitelist

Проверка IP в черном списке
grpcurl -plaintext -d '{\"subnet\":\"192.168.1.1\"}' localhost:50051 antibruteforce.AntiBruteForce/CheckBlacklist

Сброс бакета
grpcurl -plaintext -d '{\"login\":\"testuser\", \"ip\":\"192.168.1.1\"}' localhost:50051 antibruteforce.AntiBruteForce/ResetBucket

Проверка авторизации
grpcurl -plaintext -d '{\"login\":\"testuser\", \"password\":\"testpass\", \"ip\":\"192.168.1.1/25\"}' localhost:50051 antibruteforce.AntiBruteForce/CheckAuthorization
