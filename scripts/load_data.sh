sudo docker cp ./sql/dummy_data.sql local_postgres:/tmp
sudo docker exec -ti local_postgres psql -U myusername -f /tmp/dummy_data.sql
