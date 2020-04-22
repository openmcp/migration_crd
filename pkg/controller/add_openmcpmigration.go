package controller

import (
	"nanum.co.kr/openmcp/migration/pkg/controller/openmcpmigration"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, openmcpmigration.Add)
}
