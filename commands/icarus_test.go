package commands

import (
	"testing"
	"github.com/stretchr/testify/require"
);

func TestMoveInvalidDirrection(t *testing.T){
	_ , err := Move("leftish")
	require.Error(t, err)
}

func TestMoveLeft(t *testing.T){
	_ , err := Move("left")
	require.Error(t, err)
}
