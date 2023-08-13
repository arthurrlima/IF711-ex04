go build client.go

# for j in {1}
# do
#    ./client &
# done

for i in {1..5}
do
   ./client &
done
read
