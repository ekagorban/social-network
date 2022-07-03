# social-network

docker network create sn-network
docker run -tid -p 3306:3306 --name mysql-db --network sn-network -e MYSQL_ROOT_PASSWORD=secret -e MYSQL_DATABASE=test -d mysql:8.0
docker build --tag social-network .
docker run -p 3004:3004 --name social-network --network sn-network social-network