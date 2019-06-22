local:
	go build -o rpi.darwin
arm:
	GOOS=linux GOARCH=arm GOARM=6 go build -o rpi.linux.arm

deploy: arm
	ssh pi@192.168.0.24 -- 'sudo killall -9 rpi.linux.arm || exit 0'
	scp rpi.linux.arm pi@192.168.0.24:.
	ssh pi@192.168.0.24 -- 'sudo killall -9 rpi.linux.arm; sleep 1; sudo ./rpi.linux.arm'