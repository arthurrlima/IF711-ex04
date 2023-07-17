go build client.go

for i in {1..10}
do
   ./client &
done
read