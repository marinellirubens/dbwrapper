sudo docker cp ./sql/dummy_data.sql postgresql:/tmp
sudo docker exec -ti postgresql psql -U myusername -f /tmp/dummy_data.sql
