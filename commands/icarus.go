// Copyright © 2015 Steve Francia <spf@spf13.com>.
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
//

package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/als753/gc6/mazelib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// Defining the icarus command.
// This will be called as 'laybrinth icarus'
var icarusCmd = &cobra.Command{
	Use:     "icarus",
	Aliases: []string{"client"},
	Short:   "Start the laybrinth solver",
	Long: `Icarus wakes up to find himself in the middle of a Labyrinth.
  Due to the darkness of the Labyrinth he can only see his immediate cell and if
  there is a wall or not to the top, right, bottom and left. He takes one step
  and then can discover if his new cell has walls on each of the four sides.

  Icarus can connect to a Daedalus and solve many laybrinths at a time.`,
	Run: func(cmd *cobra.Command, args []string) {
		RunIcarus()
	},
}

const (
	UP = "up"
	DOWN = "down"
	RIGHT = "right"
	LEFT = "left"
)

func init() {
	RootCmd.AddCommand(icarusCmd)
}

func RunIcarus() {
	// Run the solver as many times as the user desires.
	fmt.Println("Solving", viper.GetInt("times"), "times")
	for x := 0; x < viper.GetInt("times"); x++ {

		solveMaze()
	}

	// Once we have solved the maze the required times, tell daedalus we are done
	makeRequest("http://127.0.0.1:" + viper.GetString("port") + "/done")
}

// Make a call to the laybrinth server (daedalus) that icarus is ready to wake up
func awake() mazelib.Survey {
	contents, err := makeRequest("http://127.0.0.1:" + viper.GetString("port") + "/awake")
	if err != nil {
		fmt.Println(err)
	}
	r := ToReply(contents)
	return r.Survey
}

// Make a call to the laybrinth server (daedalus)
// to move Icarus a given direction
// Will be used heavily by solveMaze
func Move(direction string) (mazelib.Survey, error) {
	if direction == LEFT || direction == RIGHT || direction == UP || direction == DOWN {

		contents, err := makeRequest("http://127.0.0.1:" + viper.GetString("port") + "/move/" + direction)
		if err != nil {
			return mazelib.Survey{}, err
		}

		rep := ToReply(contents)
		if rep.Victory == true {
			fmt.Println(rep.Message)
			 os.Exit(1)
			return rep.Survey, mazelib.ErrVictory
		} else {
			return rep.Survey, errors.New(rep.Message)
		}
	}

	return mazelib.Survey{}, errors.New("invalid direction")
}

// utility function to wrap making requests to the daedalus server
func makeRequest(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return contents, nil
}

// Handling a JSON response and unmarshalling it into a reply struct
func ToReply(in []byte) mazelib.Reply {
	res := &mazelib.Reply{}
	json.Unmarshal(in, &res)
	return *res
}

func solveMaze() {
	currentLocation := awake()
	currentLocation, err := moveThroughMaze(currentLocation, "")
	fmt.Println(err)

}

func moveThroughMaze(currentLocation mazelib.Survey, lastMove string) (mazelib.Survey, error) {
	var err error
	if !currentLocation.Left && lastMove != RIGHT {
		currentLocation, err = Move(LEFT)
		currentLocation, err = moveThroughMaze(currentLocation, LEFT)
		currentLocation, err = Move(RIGHT)
	}
	if !currentLocation.Top && lastMove != DOWN {
		currentLocation, err = Move(UP)
		currentLocation, err = moveThroughMaze(currentLocation, UP)
		currentLocation, err = Move(DOWN)
	}
	if !currentLocation.Right && lastMove != LEFT {
		currentLocation, err = Move(RIGHT)
		currentLocation, err = moveThroughMaze(currentLocation, RIGHT)
		currentLocation, err = Move(LEFT)
	}
	if !currentLocation.Bottom && lastMove != UP {
		currentLocation, err = Move(DOWN)
		currentLocation, err = moveThroughMaze(currentLocation, DOWN)
		currentLocation, err = Move(UP)
	}
	return currentLocation, err
}
