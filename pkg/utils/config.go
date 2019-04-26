/*
 __  __       _ _      _    _____
|  \/  | __ _(_) |    / \  |  ___|
| |\/| |/ _` | | |   / _ \ | |_
| |  | | (_| | | |_ / ___ \|  _|
|_|  |_|\__,_|_|_(_)_/   \_\_|

Send mails as fuck!
Author : Kunal Dawn (kunal.dawn@gmail.com)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>
*/
package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"sync"
)

type Config struct {
	viper *viper.Viper
}

var configInstance *Config // package private singleton instance of the configuration
var singleton sync.Once    // package private singleton helper utility

func GetConfig() *viper.Viper {
	// create an instance if not available
	singleton.Do(func() {
		configInstance = &Config{viper.New()}
	})

	return configInstance.viper
}

func Start() {
	if configInstance == nil {
		GetConfig()
	}

	// Find and read the config file
	err := configInstance.viper.ReadInConfig()
	if err != nil {
		// Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s \n", err))
	}
}
