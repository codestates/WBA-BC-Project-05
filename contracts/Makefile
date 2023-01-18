all:
	rm -rf *.go ./build
	solc --abi --bin TotoroVote.sol -o build
	abigen --bin=build/TotoroVote.bin --abi=build/TotoroVote.abi --pkg=contracts --out=TotoroVote.go
	go run ./deploy/deploy.go ../config/config.toml

clean:
	rm -rf *.go ./build

