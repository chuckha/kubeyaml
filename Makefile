kubeyaml:
	GOOS=linux go build -o kubeyaml ./cmd/server

.PHONY: clean

clean:
	rm -f kubeyaml
