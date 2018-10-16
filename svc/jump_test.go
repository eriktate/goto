package svc_test

import (
	"testing"

	"github.com/eriktate/jump/svc"
)

func Test_NewJumpSvc(t *testing.T) {
	// SETUP
	paths := []string{"../"}

	// RUN
	j := svc.NewJumpSvc(paths)

	// ASSERT
	j.Jump("cmd")
}
