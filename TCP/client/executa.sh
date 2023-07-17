go build client.go
go build client2.go

for j in {1}
do
   ./client &
done

for i in {1..4}
do
   ./client2 &
done
read
