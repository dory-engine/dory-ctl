package cmd

import (
	"fmt"
	"github.com/dory-engine/dory-ctl/pkg"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type OptionsDefReplace struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	FileNames      []string `yaml:"fileNames" json:"fileNames" bson:"fileNames" validate:""`
	Recursive      bool     `yaml:"recursive" json:"recursive" bson:"recursive" validate:""`
	Output         string   `yaml:"output" json:"output" bson:"output" validate:""`
	Param          struct {
		FileNames []string      `yaml:"fileNames" json:"fileNames" bson:"fileNames" validate:""`
		Defs      []pkg.DefKind `yaml:"defs" json:"defs" bson:"defs" validate:""`
	}
}

func NewOptionsDefReplace() *OptionsDefReplace {
	var o OptionsDefReplace
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdDefReplace() *cobra.Command {
	o := NewOptionsDefReplace()

	msgUse := fmt.Sprintf(`replace -f [filename]`)
	msgShort := fmt.Sprintf("replace project definition")
	msgLong := fmt.Sprintf(`replace project definition in dory-core server by file name or stdin.
# JSON and YAML formats are accepted, the complete definition must be provided.
# YAML format support replace multiple project definitions at the same time.
# if [filename] is a directory, it will read all *.json and *.yaml and *.yml files in this directory.`)
	msgExample := fmt.Sprintf(`  # replace project definition from file or directory
  doryctl def replace -f def1.yaml -f def2.json

  # replace project definition from stdin
  cat def1.yaml | doryctl def replace -f -`)

	cmd := &cobra.Command{
		Use:                   msgUse,
		DisableFlagsInUseLine: true,
		Short:                 msgShort,
		Long:                  msgLong,
		Example:               msgExample,
		Run: func(cmd *cobra.Command, args []string) {
			CheckError(o.Complete(cmd))
			CheckError(o.Validate(args))
			CheckError(o.Run(args))
		},
	}
	cmd.Flags().StringVarP(&o.Output, "output", "o", "", "output format (options: yaml / json)")
	cmd.Flags().BoolVarP(&o.Recursive, "recursive", "r", false, "process the directory used in -f, --file recursively.")
	cmd.Flags().StringSliceVarP(&o.FileNames, "file", "f", []string{}, "project definition file name or directory, support *.json and *.yaml and *.yml files")
	return cmd
}

func (o *OptionsDefReplace) Complete(cmd *cobra.Command) error {
	var err error
	err = o.GetOptionsCommon()
	return err
}

func (o *OptionsDefReplace) Validate(args []string) error {
	var err error

	if len(o.FileNames) == 0 {
		err = fmt.Errorf("--file required")
		return err
	}
	var fileNames []string
	for _, name := range o.FileNames {
		fileNames = append(fileNames, strings.Trim(name, " "))
	}
	var isStdin bool
	for _, name := range fileNames {
		if name == "-" {
			isStdin = true
			break
		}
	}
	if isStdin && len(fileNames) > 1 {
		err = fmt.Errorf(`"--file -" found, can not use multiple --file options`)
		return err
	}

	if isStdin {
		bs, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		if len(bs) == 0 {
			err = fmt.Errorf("--file - required os.stdin\n example: echo 'xxx' | %s def replace -f -", pkg.BaseCmdName)
			return err
		}
		defs, err := GetDefKinds("", bs)
		if err != nil {
			return err
		}
		o.Param.Defs = append(o.Param.Defs, defs...)
	} else {
		for _, fileName := range fileNames {
			fi, err := os.Stat(fileName)
			if err != nil {
				return err
			}
			if fi.IsDir() {
				if o.Recursive {
					err = filepath.Walk(fileName, func(path string, info os.FileInfo, err error) error {
						if err != nil {
							return err
						}
						ext := filepath.Ext(path)
						if !info.IsDir() && (ext == ".json" || ext == ".yaml" || ext == ".yml") {
							o.Param.FileNames = append(o.Param.FileNames, path)
						}
						return nil
					})
				} else {
					infos, err := ioutil.ReadDir(fileName)
					if err != nil {
						return err
					}
					for _, info := range infos {
						ext := filepath.Ext(info.Name())
						if !info.IsDir() && (ext == ".json" || ext == ".yaml" || ext == ".yml") {
							if strings.HasSuffix(fileName, "/") {
								fileName = strings.TrimSuffix(fileName, "/")
							}
							o.Param.FileNames = append(o.Param.FileNames, fmt.Sprintf("%s/%s", fileName, info.Name()))
						}
					}
				}
			} else {
				ext := filepath.Ext(fileName)
				if ext != ".json" && ext != ".yaml" && ext != ".yml" {
					err = fmt.Errorf("file %s error: file extension name not json, yaml or yml", fileName)
					return err
				}
				o.Param.FileNames = append(o.Param.FileNames, fileName)
			}
		}

		fileNames = []string{}
		m := map[string]bool{}
		for _, fileName := range o.Param.FileNames {
			m[fileName] = true
		}
		for fileName, _ := range m {
			fileNames = append(fileNames, fileName)
		}
		sort.Strings(fileNames)
		o.Param.FileNames = fileNames

		for _, fileName := range o.Param.FileNames {
			bs, err := os.ReadFile(fileName)
			if err != nil {
				err = fmt.Errorf("read file %s error: %s", fileName, err.Error())
				return err
			}

			defs, err := GetDefKinds(fileName, bs)
			if err != nil {
				return err
			}
			o.Param.Defs = append(o.Param.Defs, defs...)
		}
	}

	if o.Output != "" {
		if o.Output != "yaml" && o.Output != "json" {
			err = fmt.Errorf("--output must be yaml or json")
			return err
		}
	}
	return err
}

func (o *OptionsDefReplace) Run(args []string) error {
	var err error

	bs, _ := pkg.YamlIndent(o)
	log.Debug(fmt.Sprintf("command options:\n%s", string(bs)))

	//param := map[string]interface{}{}
	//result, _, err := o.QueryAPI(fmt.Sprintf("api/cicd/projectDef/%s", o.Param.ProjectName), http.MethodGet, "", param, false)
	//if err != nil {
	//	return err
	//}
	//project := pkg.ProjectOutput{}
	//err = json.Unmarshal([]byte(result.Get("data.project").Raw), &project)
	//if err != nil {
	//	return err
	//}

	return err
}
