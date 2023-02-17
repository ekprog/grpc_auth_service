evans \
./api/auth_service.proto \
-p 8086

show package \
show message \
show service \
show header

package pb
service AuthService
call Register
call --repeat Register


protoc -I ./proto \
--go_out ./pkg/pb \
--go_opt paths=source_relative \
--go-grpc_out ./pkg/pb \
--go-grpc_opt paths=source_relative \
--grpc-gateway_out ./pkg/pb \
--grpc-gateway_opt paths=source_relative \
--openapiv2_out ./docs \
./proto/api/**/*.proto

**Make certs**
openssl genrsa -out ca.key 4096
openssl req -new -x509 -key ca.key -sha256 -subj "/C=US/ST=NJ/O=CA, Inc." -days 365 -out ca.cert
openssl genrsa -out service.key 4096
openssl req -new -key service.key -out service.csr -config certificate.conf
openssl x509 -req -in service.csr -CA ca.cert -CAkey ca.key -CAcreateserial -out service.pem -days 365 -sha256 -extfile certificate.conf -extensions req_ext
