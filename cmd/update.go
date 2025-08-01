package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Azure/draft/pkg/cmdhelpers"
	"github.com/Azure/draft/pkg/config"
	dryrunpkg "github.com/Azure/draft/pkg/dryrun"
	"github.com/Azure/draft/pkg/handlers"
	"github.com/Azure/draft/pkg/templatewriter"
	"github.com/Azure/draft/pkg/templatewriter/writers"
)

type updateCmd struct {
	dest                     string
	provider                 string
	addon                    string
	flagVariables            []string
	templateWriter           templatewriter.TemplateWriter
	templateVariableRecorder config.TemplateVariableRecorder
}

var dryRunRecorder *dryrunpkg.DryRunRecorder

func newUpdateCmd() *cobra.Command {
	uc := &updateCmd{}
	// updateCmd represents the update command
	var cmd = &cobra.Command{
		Use:   "update",
		Short: "Updates your application to be internet accessible",
		Long: `This command automatically updates your yaml files as necessary so that your application
		will be able to receive external requests.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := uc.run(); err != nil {
				return err
			}
			log.Info("Draft has successfully updated your yaml files so that your application will be able to receive external requests 😃")
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVarP(&uc.dest, "destination", "d", ".", "specify the path to the project directory")
	f.StringVarP(&uc.provider, "provider", "p", "azure", "cloud provider")
	f.StringVarP(&uc.addon, "addon", "a", "", "addon name")
	f.StringArrayVarP(&uc.flagVariables, "variable", "", []string{}, "pass template variables (e.g. --variable ingress-tls-cert-keyvault-uri=test.uri ingress-host=host)")

	uc.templateWriter = &writers.LocalFSWriter{}

	return cmd
}

func (uc *updateCmd) run() error {
	flagVariablesMap = flagVariablesToMap(uc.flagVariables)

	if dryRun {
		dryRunRecorder = dryrunpkg.NewDryRunRecorder()
		uc.templateVariableRecorder = dryRunRecorder
		uc.templateWriter = dryRunRecorder
	}

	updatedDest, err := cmdhelpers.GetAddonDestPath(uc.dest)
	if err != nil {
		log.Errorf("error getting addon destination path: %s", err.Error())
		return err
	}

	// Use the specified addon template, default to app-routing-ingress for backward compatibility
	templateName := "app-routing-ingress"
	if uc.addon != "" {
		templateName = uc.addon
	}

	addonTemplate, err := handlers.GetTemplate(templateName, "", updatedDest, uc.templateWriter)
	if err != nil {
		log.Errorf("error getting addon template: %s", err.Error())
		return err
	}
	if addonTemplate == nil {
		return errors.New("DraftConfig is nil")
	}

	addonTemplate.Config.VariableMapToDraftConfig(flagVariablesMap)

	err = cmdhelpers.PromptAddonValues(uc.dest, addonTemplate.Config)
	if err != nil {
		return err
	}

	if dryRun {
		for _, variable := range addonTemplate.Config.Variables {
			uc.templateVariableRecorder.Record(variable.Name, variable.Value)
		}
	}

	err = addonTemplate.Generate()
	if err != nil {
		log.Errorf("error generating addon template: %s", err.Error())
		return err
	}

	if dryRun {
		dryRunText, err := json.MarshalIndent(dryRunRecorder.DryRunInfo, "", TWO_SPACES)
		if err != nil {
			return err
		}
		fmt.Println(string(dryRunText))
		if dryRunFile != "" {
			log.Printf("writing dry run info to file %s", dryRunFile)
			err = os.WriteFile(dryRunFile, dryRunText, 0644)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func init() {
	rootCmd.AddCommand(newUpdateCmd())
}
