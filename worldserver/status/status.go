/*
   Copyright 2014 Franc[e]sco (lolisamurai@tfwno.gf)
   This file is part of kagami.
   kagami is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   kagami is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
   GNU General Public License for more details.
   You should have received a copy of the GNU General Public License
   along with kagami. If not, see <http://www.gnu.org/licenses/>.
*/

// Package status contains various information about the current status of the world server
// such as world config, port and connections that are shared globally within the package
package status

import (
	"github.com/Francesco149/kagami/common"
	"github.com/Francesco149/kagami/common/config"
)

// TODO: move this stuff to channels.Connection?

var worldconf *config.WorldConf = nil
var worldport int16 = -1
var loginconn *common.InterserverClient = nil // connection to the loginserver

func Conf() *config.WorldConf                  { return worldconf }
func SetConf(c *config.WorldConf)              { worldconf = c }
func Port() int16                              { return worldport }
func SetPort(port int16)                       { worldport = port }
func LoginConn() *common.InterserverClient     { return loginconn }
func SetLoginConn(c *common.InterserverClient) { loginconn = c }
