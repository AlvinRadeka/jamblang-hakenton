.PHONY: build
build:
	@go build -o ./build/app ./cmd

.PHONY: run
run:
	@./build/app

.PHONY: clean
clean:
	@rm -rf ./build/*

.PHONY: build-run
build-run: build run

.PHONY: clean-build-run
clean-build-run: clean build run

.PHONY: publish-ecr
	@docker build -t warehouse:latest .
	@docker tag warehouse:latest 689895324334.dkr.ecr.us-east-1.amazonaws.com/warehouse
	@docker push 689895324334.dkr.ecr.us-east-1.amazonaws.com/warehouse