package dow

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/phoebetron/backup/pkg/cli/apicliaws"
	"github.com/phoebetron/trades/sto/tradesredis"
	"github.com/phoebetron/trades/typ/trades"
	"github.com/spf13/cobra"
	"github.com/xh3b4sd/framer"
	"google.golang.org/protobuf/proto"
)

type run struct {
	cliaws *apicliaws.AWS
	cmdfla *fla
	misfra framer.Frames
	stotra trades.Storage
}

func (r *run) run(cmd *cobra.Command, args []string) {
	var err error

	{
		r.cmdfla.Verify()
	}

	// --------------------------------------------------------------------- //

	{
		r.cliaws = apicliaws.Default()
	}

	{
		r.stotra = tradesredis.Default()
	}

	{
		r.misfra = r.franew()
	}

	// --------------------------------------------------------------------- //

	var sta time.Time
	var end time.Time
	{
		sta = r.misfra.Min().Sta
		end = r.misfra.Max().End
	}

	// --------------------------------------------------------------------- //

	{
		fmt.Printf("checking backup between %s and %s\n", scrfmt(sta), scrfmt(end))
	}

	var buc string
	{
		buc = "xh3b4sd-phoebe-backup"
	}

	var pre string
	{
		pre = fmt.Sprintf("tra-raw.exc-%s.ass-%s", r.stotra.Market().Exc(), r.stotra.Market().Ass())
	}

	var suf string
	{
		suf = fmt.Sprintf("%s.pb.raw", bacfmt(sta))
	}

	var byt []byte
	{
		byt, err = r.cliaws.Download(buc, filepath.Join(pre, suf))
		if err != nil {
			panic(err)
		}
	}

	tra := &trades.Trades{}
	{
		err := proto.Unmarshal(byt, tra)
		if err != nil {
			panic(err)
		}
	}

	{
		err := r.stotra.Create(tra.ST.AsTime(), tra)
		if tradesredis.IsAlreadyExists(err) {
			err = r.stotra.Update(tra.ST.AsTime(), tra)
			if err != nil {
				panic(err)
			}
		} else if err != nil {
			panic(err)
		}
	}
}
