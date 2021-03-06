// Copyright 2018 The Kubeflow Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	kftypes "github.com/kubeflow/kubeflow/bootstrap/v3/pkg/apis/apps"
	"github.com/kubeflow/kubeflow/bootstrap/v3/pkg/kfapp/coordinator"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deleteCfg = viper.New()

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Args:  cobra.NoArgs,
	Use:   "delete",
	Short: "Delete a kubeflow application.",
	Long:  `Delete a kubeflow application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.SetLevel(log.InfoLevel)
		if deleteCfg.GetBool(string(kftypes.VERBOSE)) != true {
			log.SetLevel(log.WarnLevel)
		}
		// Load config from exisiting app.yaml
		if configFilePath == "" {
			return fmt.Errorf("Must pass in -f configFile")
		}
		kfApp, err = coordinator.BuildKfAppFromURI(configFilePath)
		if err != nil || kfApp == nil {
			return fmt.Errorf("error loading kfapp: %v", err)
		}
		// TODO(lunkai): do we need set delete storage here?
		// kfGetter, ok := kfApp.(coordinator.KfDefGetter)
		// if !ok {
		// 	return errors.New("internal error: coordinator does not implement KfDefGetter")
		// }
		// kfGetter.GetKfDef().Spec.DeleteStorage = deleteCfg.GetBool(string(kftypes.DELETE_STORAGE))
		deleteErr := kfApp.Delete(kftypes.ALL)
		if deleteErr != nil {
			return fmt.Errorf("couldn't delete KfApp: %v", deleteErr)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCfg.SetConfigName("app")
	deleteCfg.SetConfigType("yaml")

	deleteCmd.PersistentFlags().StringVarP(&configFilePath, string(kftypes.FILE), "f", "",
		"The local config file of KfDef.")

	// verbose output
	deleteCmd.Flags().BoolP(string(kftypes.VERBOSE), "V", false,
		string(kftypes.VERBOSE)+" output default is false")
	bindErr := deleteCfg.BindPFlag(string(kftypes.VERBOSE), deleteCmd.Flags().Lookup(string(kftypes.VERBOSE)))
	if bindErr != nil {
		log.Errorf("couldn't set flag --%v: %v", string(kftypes.VERBOSE), bindErr)
		return
	}

	deleteCmd.Flags().Bool(string(kftypes.DELETE_STORAGE), false,
		"Set if you want to delete app's storage cluster used for mlpipeline.")
	bindErr = deleteCfg.BindPFlag(string(kftypes.DELETE_STORAGE), deleteCmd.Flags().Lookup(string(kftypes.DELETE_STORAGE)))
	if bindErr != nil {
		log.Errorf("couldn't set flag --%v: %v", string(kftypes.DELETE_STORAGE), bindErr)
		return
	}
}
