test-publish-endpoint:
	curl -d '{"key":"val"}' http://localhost:8080/publish?t=test1