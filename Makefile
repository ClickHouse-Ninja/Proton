proto:
	@echo "Make Pinba proto"
	@protoc --gofast_out=. proto/pinba/*.proto

.PHONY: proto