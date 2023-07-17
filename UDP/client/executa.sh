go build client.go
go build client2.go



for i in {1..5}
do
   ./client &
done
read