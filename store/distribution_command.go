package store

import (
	"flag"
	"io"
	"time"

	"github.com/cloudfoundry-incubator/runtime-schema/bbs"
	"github.com/cloudfoundry/storeadapter/etcdstoreadapter"
	"github.com/onsi/say"
	"github.com/pivotal-cf-experimental/veritas/common"
	"github.com/pivotal-cf-experimental/veritas/config_finder"
	"github.com/pivotal-cf-experimental/veritas/store/fetch_store"
	"github.com/pivotal-cf-experimental/veritas/store/print_store"
)

func DistributionCommand() common.Command {
	var (
		etcdClusterFlag   string
		consulClusterFlag string
		tasks             bool
		lrps              bool
		rate              time.Duration
	)

	flagSet := flag.NewFlagSet("distribution", flag.ExitOnError)
	flagSet.StringVar(&etcdClusterFlag, "etcdCluster", "", "comma-separated etcd cluster urls")
	flagSet.StringVar(&consulClusterFlag, "consulCluster", "", "comma-separated consul cluster urls")
	flagSet.BoolVar(&tasks, "tasks", true, "print tasks")
	flagSet.BoolVar(&lrps, "lrps", true, "print lrps")
	flagSet.DurationVar(&rate, "rate", time.Duration(0), "rate at which to poll the store")

	return common.Command{
		Name:        "distribution",
		Description: "- Fetch and print distribution of Tasks and LRPs",
		FlagSet:     flagSet,
		Run: func(args []string) {
			veritasBBS, etcdStore, err := config_finder.ConstructBBS(etcdClusterFlag, consulClusterFlag)
			common.ExitIfError("Could not construct BBS", err)

			if rate == 0 {
				err = distribution(veritasBBS, etcdStore, tasks, lrps, false)
				common.ExitIfError("Failed to print distribution", err)
				return
			}

			ticker := time.NewTicker(rate)
			for {
				<-ticker.C
				err = distribution(veritasBBS, etcdStore, tasks, lrps, true)
				if err != nil {
					say.Println(0, say.Red("Failed to print distribution: %s", err.Error()))
				}
			}
		},
	}
}

func distribution(veritasBBS bbs.VeritasBBS, etcdStore *etcdstoreadapter.ETCDStoreAdapter, tasks bool, lrps bool, clear bool) error {
	reader, writer := io.Pipe()

	errs := make(chan error)
	go func() {
		err := fetch_store.Fetch(veritasBBS, etcdStore, false, writer)
		errs <- err
	}()
	go func() {
		err := print_store.PrintDistribution(tasks, lrps, clear, reader)
		errs <- err
	}()

	err1 := <-errs
	if err1 != nil {
		return err1
	}
	err2 := <-errs
	if err2 != nil {
		return err2
	}
	return nil
}
