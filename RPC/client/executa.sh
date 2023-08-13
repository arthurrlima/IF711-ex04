go build client.go
go build client2.go

for j in {1}
do
   ./client2 20 &
done

for i in {1..19}
do
   ./client $i &
done
wait
