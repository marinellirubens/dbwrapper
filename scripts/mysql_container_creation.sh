sudo docker run --name mysql \
  -e MYSQL_ROOT_PASSWORD=strong_password \
  -e MYSQL_DATABASE=example_db \
  -e MYSQL_USER=user \
  -e MYSQL_PASSWORD=user_password \
  -p "3306:3306" \
  -v db_data:/var/lib/mysql \
  -d mysql:8.3.0
