echo "=========== Redis Benchmark Test ==========="
rm -rf ./results/*

echo "=============== Start Redis ================"
docker run --name redis -p 6379 -d redis

echo "================ Performance of 10 bytes ================" >> ./results/10_bytes
docker exec -it redis redis-benchmark -t set,get -d 10 >> ./results/10_bytes
echo "================ Performance of 20 bytes ================" >> ./results/20_bytes
docker exec -it redis redis-benchmark -p 6379 -t set,get -d 20 >> ./results/20_bytes
echo "================ Performance of 50 bytes ================" >> ./results/50_bytes
docker exec -it redis redis-benchmark -p 6379 -t set,get -d 50 >> ./results/50_bytes
echo "================ Performance of 100 bytes ================" >> ./results/100_bytes
docker exec -it redis redis-benchmark -p 6379 -t set,get -d 100 >> ./results/100_bytes
echo "================ Performance of 200 bytes ================" >> ./results/200_bytes
docker exec -it redis redis-benchmark -p 6379 -t set,get -d 200 >> ./results/200_bytes
echo "================ Performance of 1000 bytes ================" >> ./results/1000_bytes
docker exec -it redis redis-benchmark -p 6379 -t set,get -d 1000 >> ./results/1000_bytes
echo "================ Performance of 5000 bytes ================" >> ./results/5000_bytes
docker exec -it redis redis-benchmark -p 6379 -t set,get -d 5000 >> ./results/5000_bytes

echo "=============== Stop Redis ================"
docker stop redis &> /dev/null
docker rm -f redis &> /dev/null