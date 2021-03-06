package upl

import "github.com/phoebetron/trades/typ/trades"

func (r *run) tra() []*trades.Trade {
	var err error

	var all []*trades.Trade
	for _, h := range r.misfra {
		var tra *trades.Trades
		{
			tra, err = r.stotra.Search(h.Sta)
			if err != nil {
				panic(err)
			}
		}

		{
			all = append(all, tra.TR...)
		}
	}

	return all
}
