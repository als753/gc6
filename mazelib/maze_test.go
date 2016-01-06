package mazelib

import (
	"testing"
	"github.com/stretchr/testify/require"
);

func TestAddWall(t *testing.T){
	var r Room
	r.AddWall(TOP)
	require.Equal(t, true, r.Walls.Top)
}
