package view

import (
	"blind-tools/model"
	"maps"
	"sort"
)

// planStep describes how a single draw is paid.
type planStep struct {
	Draw         int
	CurrencyID   string
	CurrencyName string
	Cost         int
	Remaining    int // remaining amount of the spent currency after this draw
}

// planResult is the outcome of calculating a spending plan.
type planResult struct {
	Steps        []planStep
	Final        map[string]int // final balance for every currency
	Insufficient bool
	FailAtDraw   int
}

// calculatePlan computes, for every draw in the inclusive range [start, end],
// which currency to spend so that the preferred currency (preferKeep) is used
// as little as possible. Non-preferred currencies are spent first (in the order
// they appear in the container), and the preferred currency is only spent when
// no other currency can cover the draw.
func calculatePlan(container model.Container, start, end int, balances map[string]int, preferKeep string) planResult {
	costs := buildCostLookup(container)
	order := buildSpendOrder(container, preferKeep)

	res := planResult{
		Final: make(map[string]int, len(balances)),
	}

	maps.Copy(res.Final, balances)

	for draw := start; draw <= end; draw++ {
		cost, ok := costs[draw]
		if !ok {
			res.Insufficient = true
			res.FailAtDraw = draw
			return res
		}

		chosen := ""
		for _, id := range order {
			amount, ok := cost[id]
			if !ok {
				continue // this draw does not accept this currency
			}
			if res.Final[id] >= amount {
				chosen = id
				break
			}
		}
		if chosen == "" {
			res.Insufficient = true
			res.FailAtDraw = draw
			return res
		}

		res.Final[chosen] -= cost[chosen]
		res.Steps = append(res.Steps, planStep{
			Draw:         draw,
			CurrencyID:   chosen,
			CurrencyName: currencyName(container, chosen),
			Cost:         cost[chosen],
			Remaining:    res.Final[chosen],
		})
	}

	return res
}

// buildCostLookup maps a draw number (1 based) to its cost table.
func buildCostLookup(container model.Container) map[int]map[string]int {
	costs := make(map[int]map[string]int, len(container.Manifest.Prices))
	for _, price := range container.Manifest.Prices {
		// Copy so the caller can never mutate the source data.
		table := make(map[string]int, len(price.Cost))
		maps.Copy(table, price.Cost)
		costs[price.Draw] = table
	}
	return costs
}

// buildSpendOrder returns currency IDs ordered by spending priority: every
// non-preferred currency first (stable container order), then the preferred one
// last. Cost table keys not present in the container currencies are included so
// they are never ignored.
func buildSpendOrder(container model.Container, preferKeep string) []string {
	order := make([]string, 0, len(container.Currencies)+1)
	seen := make(map[string]bool, len(container.Currencies)+1)

	for _, currency := range container.Currencies {
		if currency.ID != preferKeep {
			order = append(order, currency.ID)
			seen[currency.ID] = true
		}
	}

	// Include any cost-table keys that are not declared currencies.
	for _, price := range container.Manifest.Prices {
		for id := range price.Cost {
			if !seen[id] && id != preferKeep {
				order = append(order, id)
				seen[id] = true
			}
		}
	}

	if preferKeep != "" {
		order = append(order, preferKeep)
	}

	return order
}

// currencyName resolves a currency ID to its display name.
func currencyName(container model.Container, id string) string {
	for _, currency := range container.Currencies {
		if currency.ID == id {
			return currency.Name
		}
	}
	return id
}

// sortedCurrencyIDs returns the container currency IDs in a stable order, used
// to present final balances consistently.
func sortedCurrencyIDs(container model.Container) []string {
	ids := make([]string, 0, len(container.Currencies))
	for _, currency := range container.Currencies {
		ids = append(ids, currency.ID)
	}
	sort.Strings(ids)
	return ids
}
