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
package cmd

import (
	"github.com/spf13/cobra"
	//"github.com/spf13/viper"
	"melato.org/lxdtool/op"
)

func ServerFlags(cmd *cobra.Command, server *op.Server) {
	cmd.PersistentFlags().StringVarP(&server.Socket, "socket", "s", "/var/snap/lxd/common/lxd/unix.socket", "path to unix socket")
	cmd.PersistentFlags().StringVarP(&server.Remote, "remote", "r", "", "LXD remote")
	cmd.PersistentFlags().StringVarP(&server.ConfigDir, "config", "c", "", "config dir (with client.crt, client.key)")
}
