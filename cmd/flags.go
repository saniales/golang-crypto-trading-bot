// Copyright © 2017 Alessandro Sanino <saninoale@gmail.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package bot

//GlobalFlags provides flag definitions valid for the whole system.
var GlobalFlags struct {
	Verbose    int    //Tells the program to print everything to screen (used multiple times for better verbosity).
	ConfigFile string //Config file path (assumed ./.gobot if not specified)
}

//rootFlags provides flag definitions valid for root command.
var rootFlags struct {
	Version bool
}

// initFlags provdes flag definition for init command.
var initFlags struct {
	ConfigFile string
	Exchange   string
	Strategies []struct {
		Market   string
		Strategy string
	}
	BTCAddress string
}

// startFlags provdes flag definition for start command.
var startFlags struct {
	Simulate bool
}
