package model

import (
	"fmt"
	"sort"
)

// ValidateManifestPrices 校验盲盒 manifest 的价格表是否覆盖了该容器支持的所有货币
func ValidateManifestPrices(container Container) error {
	required := make(map[string]struct{}, len(container.Currencies))
	for _, currency := range container.Currencies {
		required[currency.ID] = struct{}{}
	}

	if len(required) == 0 {
		return nil
	}

	covered := make(map[string]struct{})
	for _, price := range container.Manifest.Prices {
		for currencyID := range price.Cost {
			covered[currencyID] = struct{}{}
		}
	}

	var missing []string
	for currencyID := range required {
		if _, ok := covered[currencyID]; !ok {
			missing = append(missing, currencyID)
		}
	}

	if len(missing) > 0 {
		sort.Strings(missing)
		return fmt.Errorf("manifest %s (%s) price table missing currencies: %v",
			container.Manifest.ID, container.Manifest.Name, missing)
	}

	return nil
}

// ValidateManifestDraws 校验盲盒 manifest 的价格表条目数严格等于 draws
func ValidateManifestDraws(container Container) error {
	actual := len(container.Manifest.Prices)
	expected := container.Manifest.Draws

	if actual != expected {
		return fmt.Errorf("manifest %s (%s) price table has %d entries, want %d (draws)",
			container.Manifest.ID, container.Manifest.Name, actual, expected)
	}

	return nil
}
