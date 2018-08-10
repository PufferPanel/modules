package mojangdl

import (
	"testing"
	"github.com/pufferpanel/pufferd/programs/operations/ops"
	"github.com/pufferpanel/modules/operations"
)

func TestMojangDlOperationFactory_Create(t *testing.T) {
	var factory ops.OperationFactory

	factory = operations.MojangDlOperationFactory{}

	if factory.Key() != "mojangdl" {
		t.Fail()
		return
	}

	version := "1.13"
	filename := "server.jar"

	createCmd := ops.CreateOperation{
		DataMap: make(map[string]interface{}),
	}

	createCmd.DataMap["version"] = version
	createCmd.DataMap["target"] = filename


	op := factory.Create(createCmd)

	err := op.Run(nil)

	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
}