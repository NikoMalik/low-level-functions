


bench:
	@go test  -benchmem    -benchtime=10s  -bench ./   > benchmark_results.txt
