docker run -p 5432:5432 \
    -v pgsql_data:/var/lib/postgresql/data \
    -e POSTGRES_PASSWORD=mypassword \
    -e POSTGRES_USER=myusername \
    -e POSTGRES_DB=myusername \
    --name local_postgres \
    -d postgres
