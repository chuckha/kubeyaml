kube-validate:
	GOOS=linux go build -o kube-validate ./cmd/server

.PHONY: clean

clean:
	rm -f kube-validate
