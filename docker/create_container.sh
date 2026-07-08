TAG=$(cat VERSION)
echo $TAG
IMAGE_NAME="dbwrapper"
CWD=$(pwd)
echo $CWD

sudo docker rm -fv $IMAGE_NAME || true
sudo docker run -d \
  --name $IMAGE_NAME \
  --link postgresql:postgresql \
  -p "8080:8080" \
  -v "/var/log/dbwrapper:/tmp/dbwrapper" \
  -v "${CWD}/config.json:/opt/config.json:ro" \
  "$IMAGE_NAME:$TAG"

