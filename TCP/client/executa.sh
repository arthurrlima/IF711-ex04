go build client.go
go build client2.go

for j in {1}
do
   ./client 5 &
done

for i in {1..4}
do
   ./client2 $i &
done
read
