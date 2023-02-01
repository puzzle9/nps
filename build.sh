export GOPROXY=direct

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -extldflags -static -extldflags -static" ./cmd/nps/nps.go
tar -czvf linux_amd64_server.tar.gz conf/nps.conf conf/tasks.json conf/clients.json conf/hosts.json web/dist nps

#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -extldflags -static -extldflags -static" ./cmd/npc/npc.go
#tar -czvf linux_amd64_client.tar.gz npc

#CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -extldflags -static -extldflags -static" ./cmd/npc/npc.go
#tar -czvf windows_amd64_client.tar.gz npc.exe
