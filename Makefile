SERVER_IP=""

deploy: build upload restart

install: user build upload data service

build:
	go build

clean-data:
	rm -f data/sessions/*
	rm -f data/transactions/*

test: build clean-data
	-pkill -f smtp_honeypot
	./smtp_honeypot &
	cd tests && python3 runner.py

upload:
	scp -i smtp_honeypot.pem smtp_honeypot "ubuntu@$(SERVER_IP):~/"
	ssh -i smtp_honeypot.pem "ubuntu@$(SERVER_IP)" sudo mv smtp_honeypot /usr/local/bin

data:
	ssh -i smtp_honeypot.pem "ubuntu@$(SERVER_IP)" sudo mkdir -p /etc/smtp_honeypot
	ssh -i smtp_honeypot.pem "ubuntu@$(SERVER_IP)" sudo chown smtp_honeypot /etc/smtp_honeypot
	scp -i smtp_honeypot.pem server.crt "ubuntu@$(SERVER_IP):~/"
	scp -i smtp_honeypot.pem server.key "ubuntu@$(SERVER_IP):~/"
	ssh -i smtp_honeypot.pem "ubuntu@$(SERVER_IP)" sudo mv server.* /etc/smtp_honeypot
	ssh -i smtp_honeypot.pem "ubuntu@$(SERVER_IP)" sudo chown smtp_honeypot /etc/smtp_honeypot/*
	ssh -i smtp_honeypot.pem "ubuntu@$(SERVER_IP)" sudo mkdir -p /etc/smtp_honeypot/data/sessions
	ssh -i smtp_honeypot.pem "ubuntu@$(SERVER_IP)" sudo mkdir -p /etc/smtp_honeypot/data/transactions
	ssh -i smtp_honeypot.pem "ubuntu@$(SERVER_IP)" sudo chown smtp_honeypot /etc/smtp_honeypot/data/sessions
	ssh -i smtp_honeypot.pem "ubuntu@$(SERVER_IP)" sudo chown smtp_honeypot /etc/smtp_honeypot/data/transactions

service:
	scp -i smtp_honeypot.pem smtp_honeypot.service "ubuntu@$(SERVER_IP):~/"
	ssh -i smtp_honeypot.pem "ubuntu@$(SERVER_IP)" sudo mv smtp_honeypot.service /etc/systemd/system/smtp_honeypot.service
	ssh -i smtp_honeypot.pem "ubuntu@$(SERVER_IP)" sudo systemctl daemon-reload
	ssh -i smtp_honeypot.pem "ubuntu@$(SERVER_IP)" sudo systemctl enable --now smtp_honeypot

restart:
	ssh -i smtp_honeypot.pem "ubuntu@$(SERVER_IP)" sudo systemctl restart smtp_honeypot.service

start:
	ssh -i smtp_honeypot.pem "ubuntu@$(SERVER_IP)" sudo systemctl start smtp_honeypot.service

stop:
	ssh -i smtp_honeypot.pem "ubuntu@$(SERVER_IP)" sudo systemctl stop smtp_honeypot.service

status:
	ssh -i smtp_honeypot.pem "ubuntu@$(SERVER_IP)" sudo systemctl status smtp_honeypot.service

user:
	ssh -i smtp_honeypot.pem "ubuntu@$(SERVER_IP)" sudo useradd --system --no-create-home --shell /usr/sbin/nologin smtp_honeypot

logs:
	ssh -i smtp_honeypot.pem "ubuntu@$(SERVER_IP)" sudo journalctl -u smtp_honeypot -n 20