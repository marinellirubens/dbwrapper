docker run -p 5432:5432 \
    -v /tmp/database:/var/lib/postgresql/data \
    -e POSTGRES_PASSWORD=mypassword \
    -e POSTGRES_USER=myusername \
    -e POSTGRES_DB=postgres \
    --name local_postgres \
    -d postgres
