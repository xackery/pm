SERVER = z420

build:
	GOOS=linux go build -o bin/pm -trimpath -ldflags "-s -w" main.go
copy: build
	scp bin/pm xackery@${SERVER}:~/pm
# execute a command on the remote server
run: copy
	ssh xackery@${SERVER} "cd ~/; ./pm"