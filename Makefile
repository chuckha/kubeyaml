kubeyaml:
	GOOS=linux go build -o kubeyaml ./cmd/kubeyaml

kubeyaml-server:
	GOOS=linux go build -o kubeyaml-server ./cmd/server

.PHONY: clean

clean:
	rm -f kubeyaml kubeyaml-server
