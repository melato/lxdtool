/* Copyright 2018 Alex Athanasopoulos
*
*  Licensed under the Apache License, Version 2.0 (the "License");
*  you may not use this file except in compliance with the License.
*  You may obtain a copy of the License at
*
*   http://www.apache.org/licenses/LICENSE-2.0
*
*  Unless required by applicable law or agreed to in writing, software
*  distributed under the License is distributed on an "AS IS" BASIS,
*  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*  See the License for the specific language governing permissions and
*  limitations under the License.
*/
package op

import (
	"github.com/lxc/lxd/client"
)

type Tool struct {
	Server Server
}

func (t *Tool) GetServer() (lxd.ContainerServer, error) {
	return t.Server.GetServer()
}
