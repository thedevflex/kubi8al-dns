build 
go build ./cmd/... -o bin/k8dns

run (for dev)
DEV_MODE=true ./bin/k8dns 

for prod
./bin/k8dns 
