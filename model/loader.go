package model

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	GlobalDirectory = "data"
	// ManifestFile 盲盒数据
	// 包含盲盒价格
	ManifestFile = "manifest.json"
	// CurrencyFile 货币信息
	CurrencyFile = "currency.json"
	// ProbabilityFile 综合概率
	// 这个目前没什么用
	ProbabilityFile = "probability.json"
)

func LoadData(directory string) ([]Container, error) {
	// 加载全局货币信息
	globalCurrency, err := LoadCurrency(filepath.Join(directory, CurrencyFile))

	if err != nil {
		log.Printf("load global currency %s failed: %v", directory, err)
	}

	// 读取目录下所有文件
	files, err := os.ReadDir(directory)

	if err != nil {
		return nil, fmt.Errorf("error read directory %s: %w", directory, err)
	}

	var containers []Container

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		containerDirectory := filepath.Join(directory, file.Name())
		manifestFile := filepath.Join(containerDirectory, ManifestFile)

		content, err := os.ReadFile(manifestFile)

		if err != nil {
			// 没有 manifest 的文件夹跳过
			if os.IsNotExist(err) {
				continue
			}
			log.Printf("skip container %s: read manifest failed: %v", containerDirectory, err)
			continue
		}

		var manifest Manifest

		if err := json.Unmarshal(content, &manifest); err != nil {
			log.Printf("skip container %s: unmarshal manifest failed: %v", containerDirectory, err)
			continue
		}

		// 加载本地货币信息
		localCurrency, err := LoadCurrency(filepath.Join(containerDirectory, CurrencyFile))

		if err != nil {
			log.Printf("container %s: load local currency failed: %v", containerDirectory, err)
		}

		if globalCurrency == nil && localCurrency == nil {
			log.Printf("skip container %s: global currency and local currency are nil", containerDirectory)
			continue
		}

		var container Container

		container.Manifest = manifest

		if localCurrency != nil {
			container.Currencies = localCurrency
		} else {
			container.Currencies = globalCurrency
		}

		if err := ValidateManifest(container); err != nil {
			log.Printf("skip container %s: %v", containerDirectory, err)
			continue
		}

		containers = append(containers, container)
	}

	return containers, nil
}

func LoadDataDefault() ([]Container, error) {
	return LoadData(GlobalDirectory)
}

func LoadCurrency(file string) ([]Currency, error) {
	text, err := os.ReadFile(file)

	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("error read currency file %s : %w", file, err)
	}

	var currencies []Currency

	if err := json.Unmarshal(text, &currencies); err != nil {
		return nil, fmt.Errorf("error unmarshal currency file %s : %w", file, err)
	}

	return currencies, nil
}
