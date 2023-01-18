all:
	rm -rf *.go ./build
	solc --abi --bin TotoroVerify.sol -o build
	abigen --bin=build/TotoroVerify.bin --abi=build/TotoroVerify.abi --pkg=contracts --out=TotoroVerify.go
	go run ./deploy/deploy.go ../config/config.toml

clean:
	rm -rf *.go ./build

