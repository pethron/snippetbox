DB=snippetboxdb

init:
	docker run --name $(DB) -e MYSQL_ROOT_PASSWORD=pwd -p 3306:3306 -d mysql:latest

start:
	docker start $(DB)
	go run ./cmd/web

stop:
	docker stop $(DB)