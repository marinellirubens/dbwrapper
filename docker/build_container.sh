TAG=$(cat VERSION)
echo $TAG
IMAGE_NAME="dbwrapper"

sudo docker build --build-arg VERSION_STR=$TAG -f docker/Dockerfile --tag "$IMAGE_NAME:$TAG" .

