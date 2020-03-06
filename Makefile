docker-server:
	GOOS=linux go build -o kubeyaml-server ./cmd/server

kubeyaml:
	go build -o kubeyaml ./cmd/kubeyaml

kubeyaml-server:
	go build -o kubeyaml-server ./cmd/server

docker:
	GOOS=linux go build -o kubeyaml ./cmd/kubeyaml
	docker build -f Dockerfile --tag=atsuio/k8slynter:latest .

.PHONY: clean

clean:
	rm -f kubeyaml kubeyaml-server docker
