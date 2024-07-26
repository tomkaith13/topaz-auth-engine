clean:
	docker-compose down --rmi all -v
clean-db:
	rm -rf ./db/directory.db
run:
	docker-compose up --build -d
stop:
	docker-compose down
load-db:
	docker compose exec topaz ./topaz directory set manifest --no-check -i /data/manifest.yaml & docker compose exec topaz ./topaz directory import --no-check -i -H localhost:9292 -d /data


