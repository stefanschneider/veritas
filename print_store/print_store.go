package print_store

import (
	"encoding/json"
	"io"

	"github.com/cloudfoundry-incubator/veritas/veritas_models"
)

func PrintStore(verbose bool, tasks bool, lrps bool, services bool, f io.Reader) error {
	decoder := json.NewDecoder(f)
	var dump veritas_models.StoreDump
	err := decoder.Decode(&dump)
	if err != nil {
		return err
	}

	if tasks {
		printTasks(verbose, dump.Tasks)
	}

	if lrps {
		printLRPS(verbose, dump.LRPS)
	}

	if services {
		printServices(verbose, dump.Services)
	}

	return nil
}
