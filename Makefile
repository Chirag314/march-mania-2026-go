.PHONY: all download features train predict clean

DATA_DIR ?= data
ART_DIR  ?= artifacts
SUB_DIR  ?= submissions

all: download features train predict

download:
	bash scripts/download_data.sh

features:
	go run ./cmd/build_features --data_dir $(DATA_DIR) --out_dir $(ART_DIR)

train:
	go run ./cmd/train --art_dir $(ART_DIR) --out_dir $(ART_DIR)

predict:
	go run ./cmd/predict --art_dir $(ART_DIR) --out_dir $(SUB_DIR)

clean:
	rm -rf $(ART_DIR) $(SUB_DIR)
