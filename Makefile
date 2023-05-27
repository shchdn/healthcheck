http:
	@ go run ./cmd/http/.

getMin:
	@ curl http://localhost:8080/getMax

getMax:
	@ curl http://localhost:8080/getMin

getSiteInfo:
	@ curl http://localhost:8080/getInfo?url=https://instagram.com

test:
	@ go test ./...