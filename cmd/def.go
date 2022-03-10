package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dory-engine/dory-ctl/pkg"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

func NewCmdDef() *cobra.Command {
	msgUse := fmt.Sprintf("def")
	msgShort := fmt.Sprintf("manage project definition")
	msgLong := fmt.Sprintf(`manage project definition in dory-core server`)
	msgExample := fmt.Sprintf(`  # get project build modules definition
  doryctl def get test-project1 build`)

	cmd := &cobra.Command{
		Use:                   msgUse,
		DisableFlagsInUseLine: true,
		Short:                 msgShort,
		Long:                  msgLong,
		Example:               msgExample,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}

	cmd.AddCommand(NewCmdDefGet())
	cmd.AddCommand(NewCmdDefReplace())
	return cmd
}

func GetDefKinds(fileName string, bs []byte) ([]pkg.DefKind, error) {
	var err error
	defs := []pkg.DefKind{}
	ext := filepath.Ext(fileName)
	if ext == ".json" {
		var def pkg.DefKind
		err = json.Unmarshal(bs, &def)
		if err != nil {
			err = fmt.Errorf("parse file %s error: %s", fileName, err.Error())
			return defs, err
		}
		defs = append(defs, def)
	} else if ext == ".yaml" || ext == ".yml" {
		dec := yaml.NewDecoder(bytes.NewReader(bs))
		var def pkg.DefKind
		for dec.Decode(&def) == nil {
			defs = append(defs, def)
		}
	} else if fileName == "" {
		var def pkg.DefKind
		err = json.Unmarshal(bs, &def)
		if err == nil {
			defs = append(defs, def)
		} else {
			dec := yaml.NewDecoder(bytes.NewReader(bs))
			for dec.Decode(&def) == nil {
				defs = append(defs, def)
			}
		}
	} else {
		err = fmt.Errorf("file extension name not json, yaml or yml")
		return defs, err
	}

	for _, def := range defs {
		if def.Kind == "" {
			err = fmt.Errorf("parse file %s error: kind is empty", fileName)
			return defs, err
		}
		if def.Metadata.ProjectName == "" {
			err = fmt.Errorf("parse file %s error: metadata.projectName is empty", fileName)
			return defs, err
		}
		var found bool
		for _, d := range pkg.DefKinds {
			if def.Kind == d {
				found = true
				break
			}
		}
		if !found {
			err = fmt.Errorf("parse file %s error: kind %s not correct", fileName, def.Kind)
			return defs, err
		}
	}

	return defs, err
}
