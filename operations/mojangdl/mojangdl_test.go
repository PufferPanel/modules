package mojangdl

import (
	"testing"
	"github.com/pufferpanel/pufferd/programs/operations/ops"
)

func TestMojangDlOperationFactory_Create(t *testing.T) {
	var factory ops.OperationFactory

	factory = MojangDlOperationFactory{}

	if factory.Key() != "mojangdl" {
		t.Fail()
		return
	}

	version := "1.13"
	filename := "server.jar"

	createCmd := ops.CreateOperation{
		OperationArgs: make(map[string]interface{}),
		DataMap: make(map[string]interface{}),
	}

	createCmd.OperationArgs["version"] = version
	createCmd.OperationArgs["target"] = filename


	op := factory.Create(createCmd)

	err := op.Run(nil)

	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
}