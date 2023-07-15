go build client.go

for i in {1..30}
do
   ./client &
done