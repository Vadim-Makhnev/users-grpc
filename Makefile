gen:
	@protoc \
	--proto_path=proto \
	--proto_path=/usr/local/include \
	--go_out=proto --go_opt=paths=source_relative \
	--go-grpc_out=proto --go-grpc_opt=paths=source_relative \
	proto/user.proto

docker-run:
	docker compose up --build -d

docker-down:
	docker compose down -v

evans:
	@evans proto/user.proto --port 4000