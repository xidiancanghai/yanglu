build:
	cd app && go build -o api_server

clean:
	rm -rf app/api_server

cp_dev_config:
	cp config/env.dev.toml config/env.toml

cp_online_config:
	cp config/env.online.toml config/env.toml

api_dev: cp_dev_config clean build
	cd app; bash restart.sh

api: cp_online_config clean build
	cd app;bash restart.sh

.PHONY: build clean cp_dev_config cp_online_config api_dev api