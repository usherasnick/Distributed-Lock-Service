filer_server_upload:
	curl -F file=@/path/to/man.jpg "http://127.0.0.1:8888/faces/man.jpg?collection=faces"

filer_server_download:
	curl -o man.jpg "http://127.0.0.1:8888/faces/man.jpg?collection=faces"

master_server_assgin:
	curl -H "Accept: application/json" "http://localhost:9333/dir/assign?pretty=y&collection=faces&ttl=90d"

volume_server_upload:
	curl -F file=@/path/to/man.jpg "http://127.0.0.1:8080/3,0430e17f9c?ttl=90d"

volume_server_download:
	curl -o man.jpg "http://127.0.0.1:8080/3,0430e17f9c"
