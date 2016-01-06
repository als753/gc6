package commands


import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/als753/gc6/mazelib"
);

func TestGetRoomShouldNotReturnInValidRoom(t *testing.T){
	maze := emptyMaze()
	_ , err := maze.GetRoom(20,10)
	require.Error(t, err, "Error was not thrown")
}

func TestGetRoomShouldReturnValidRoom(t *testing.T){
	maze := emptyMaze()
	_ , err := maze.GetRoom(5,5)
	require.NoError(t, err, "Error was thrown")
}

func TestDiscoverRoomShouldReturnSurveyOfThatRoom(t *testing.T){
	maze := emptyMaze()
	maze.rooms[5][6].Walls.Top = true
	walls , _ := maze.Discover(6,5)
	require.Equal(t, walls, mazelib.Survey{true, false, false, false})
}

func TestLookAroundShouldReturnErrVictoryIfIcarusIsAtTheTreasure(t *testing.T){
	maze := emptyMaze()
	maze.SetTreasure(1,1)
	maze.SetStartPoint(2,1)
	maze.MoveLeft()
	_, err := maze.LookAround()
	require.Equal(t, err, mazelib.ErrVictory)
}

func TestSetStartPointCannotStartAtTreasure(t *testing.T) {
	maze := emptyMaze()
	maze.SetTreasure(1, 1)
	err := maze.SetStartPoint(1, 1)
	require.Error(t, err)
}

func TestCreateMazeHasPlacedIcarusAndTreasure(t *testing.T){
	maze := createMaze()
	require.NotNil(t, maze.icarus.X)
	require.NotNil(t, maze.icarus.Y)
	require.NotNil(t, maze.end.X)
	require.NotNil(t, maze.end.Y)
	require.NotEqual(t, maze.icarus, maze.end)
}

func TestCreateMazeHasWallAtTopCorner(t *testing.T){
	maze := createMaze()
	require.True(t, maze.rooms[0][0].Walls.Top)
	require.True(t, maze.rooms[0][0].Walls.Left)
}

func TestRemoveEdgeRemovesFromBothSpacesLeft(t *testing.T){
	currentMaze = fullMaze()
	removeEdge(4,5, mazelib.LEFT)
	actualEdge, _ := currentMaze.Discover(4,5)
	edgeToLeft, _ := currentMaze.Discover(3,5)
	require.Equal(t, actualEdge, mazelib.Survey{true, true, true, false})
	require.Equal(t, edgeToLeft, mazelib.Survey{true, false, true, true})
}

func TestRemoveEdgeRemovesFromBothSpacesRight(t *testing.T){
	currentMaze = fullMaze()
	removeEdge(4,5, mazelib.RIGHT)
	actualEdge, _ := currentMaze.Discover(4,5)
	edgeToLeft, _ := currentMaze.Discover(5,5)
	require.Equal(t, actualEdge, mazelib.Survey{true, false, true, true})
	require.Equal(t, edgeToLeft, mazelib.Survey{true, true, true, false})
}

func TestRemoveEdgeRemovesFromBothSpacesTop(t *testing.T){
	currentMaze = fullMaze()
	removeEdge(4,5, mazelib.TOP)
	actualEdge, _ := currentMaze.Discover(4,5)
	edgeToLeft, _ := currentMaze.Discover(4,4)
	require.Equal(t, actualEdge, mazelib.Survey{false, true, true, true})
	require.Equal(t, edgeToLeft, mazelib.Survey{true, true, false, true})
}

func TestRemoveEdgeRemovesFromBothSpacesBottom(t *testing.T){
	currentMaze = fullMaze()
	removeEdge(4,5, mazelib.BOTTOM)
	actualEdge, _ := currentMaze.Discover(4,5)
	edgeToLeft, _ := currentMaze.Discover(4,6)
	require.Equal(t, actualEdge, mazelib.Survey{true, true, false, true})
	require.Equal(t, edgeToLeft, mazelib.Survey{false, true, true, true})
}


func TestProcessRowShouldHaveAtLeastOneOpeningOnBottom(t *testing.T){
	currentMaze = fullMaze()
	processRow(0)
	wallOnBottom := true
	for i:= 0; i < currentMaze.Width(); i ++ {
		survey, _ := currentMaze.Discover(i, 0)
		wallOnBottom = wallOnBottom && survey.Bottom
	}
	require.False(t, wallOnBottom)
}

func TestProcessRowShouldHaveEqualOpeningOnBottomAsSectionsInRow(t *testing.T){
	currentMaze = fullMaze()
	processRow(0)
	countOfOpeningsOnBottom := 0
	countOfSections := 0
	for i:= 0; i < currentMaze.Width(); i ++ {
		survey, _ := currentMaze.Discover(i, 0)
		if !survey.Bottom { countOfOpeningsOnBottom ++}
		if survey.Right { countOfSections ++}
	}
	require.Equal(t, countOfOpeningsOnBottom, countOfSections)
}

func TestProcessLastRowShouldHaveOpeningToRightIfRowAboveDoesNot(t *testing.T){
	currentMaze = fullMaze()
	lastRow := currentMaze.Height()-1
	processRow(lastRow-1)
	processLastRow(lastRow)
	mazelib.PrintMaze(currentMaze)
	for i:= 0; i < currentMaze.Width()-1; i ++ {
		currentSpace, _ := currentMaze.Discover(i, lastRow)
		spaceAbove, _ := currentMaze.Discover(i, lastRow -1)
		require.NotEqual(t, spaceAbove.Right, currentSpace.Right)
	}
}

func TestCreateMaze(t *testing.T){
	maze := createMaze()
	mazelib.PrintMaze(maze)
}

func TestMoveDirection(t *testing.T){
	createMaze()
//	c := gin.Context
//	MoveDirection(c)
}