docker rm $(docker stop $(docker ps -a -q --filter ancestor=golang-test --format="{{.ID}}"))
docker images -a | grep "golang-test" | awk '{print $3}' | xargs docker rmi
docker build -t golang-test .
docker run -d -p 8080:8081 golang-test