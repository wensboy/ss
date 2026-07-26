package config

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/spf13/cast"
	"github.com/urfave/cli/v3"
)

var (
	command_separator   = "."
	_global_command_map map[string]*cli.Command
)

func SetGlobalCommandMap(cmdMap map[string]*cli.Command) {
	_global_command_map = cmdMap
}

func GFlagSource(path ...string) Source {
	return FlagSource(_global_command_map, path...)
}

func FlagSource(cmdMap map[string]*cli.Command, path ...string) Source {
	return func() (any, bool) {
		if len(path) == 0 {
			return nil, false
		}
		if len(path) == 1 {
			idx := strings.LastIndex(path[0], command_separator)
			if idx == -1 {
				return nil, false
			}
			path = append([]string{path[0][:idx]}, path[0][idx+1:])
		}
		cmdKey := strings.Join(path[:len(path)-1], command_separator)
		flagKey := path[len(path)-1]
		if cmd, ok := cmdMap[cmdKey]; ok {
			if cmd.IsSet(flagKey) {
				return cmd.Value(flagKey), true
			}
		}
		return nil, false
	}
}

type FlagValue struct {
	Type    string      `json:"type"`
	Default interface{} `json:"default"`
	Env     []string    `json:"env"`
	Config  []string    `json:"config"`
}

type Flag struct {
	Name        string    `json:"name"`
	Usage       string    `json:"usage"`
	Description string    `json:"description"`
	Local       bool      `json:"local"`
	Required    bool      `json:"required"`
	Value       FlagValue `json:"value"`
	Short       string    `json:"short"`
}

type Command struct {
	Name        string    `json:"name"`
	Label       string    `json:"label"`
	Usage       string    `json:"usage"`
	Description string    `json:"description"`
	Flags       []Flag    `json:"flags"`
	Sub         []Command `json:"sub"`
}

type CommandConfig struct {
	Version string  `json:"version"`
	Entry   Command `json:"entry"`
}

func SetCommandSep(sep string) {
	command_separator = sep
}

func InitCommand(path string) *cli.Command {
	var cmdConfig CommandConfig
	_global_command_map = make(map[string]*cli.Command)
	fd, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	if err := json.NewDecoder(fd).Decode(&cmdConfig); err != nil {
		panic(err)
	}

	return buildCommand(cmdConfig.Entry, cmdConfig.Entry.Label)
}

func GetCommand(key string) *cli.Command {
	if _global_command_map == nil {
		panic("command map is nil, please call InitCommand first")
	}
	if cmd, ok := _global_command_map[key]; ok {
		return cmd
	}
	return nil
}

func buildCommand(cmd Command, key string) *cli.Command {
	root := &cli.Command{
		Name:        cmd.Name,
		Usage:       cmd.Usage,
		Description: cmd.Description,
	}

	for _, flag := range cmd.Flags {
		root.Flags = append(root.Flags, buildFlag(flag))
	}

	for _, subCmd := range cmd.Sub {
		root.Commands = append(root.Commands, buildCommand(subCmd, key+command_separator+subCmd.Label))
	}

	_global_command_map[key] = root
	return root
}

func buildFlag(flag Flag) cli.Flag {
	var rflag cli.Flag
	vsc := cli.NewValueSourceChain()
	vsc.Append(cli.EnvVars(flag.Value.Env...))
	vsc.Append(ConfigVars(flag.Value.Config...))
	switch flag.Value.Type {
	case "string":
		rflag = &cli.StringFlag{
			Name:     flag.Name,
			Usage:    flag.Usage,
			Local:    flag.Local,
			Required: flag.Required,
			Aliases:  []string{flag.Short},
			Sources:  vsc,
		}
	case "int":
		rflag = &cli.IntFlag{
			Name:     flag.Name,
			Usage:    flag.Usage,
			Local:    flag.Local,
			Required: flag.Required,
			Aliases:  []string{flag.Short},
			Sources:  vsc,
		}
	case "uint":
		rflag = &cli.UintFlag{
			Name:     flag.Name,
			Usage:    flag.Usage,
			Local:    flag.Local,
			Required: flag.Required,
			Aliases:  []string{flag.Short},
			Sources:  vsc,
		}
	case "bool":
		rflag = &cli.BoolFlag{
			Name:     flag.Name,
			Usage:    flag.Usage,
			Local:    flag.Local,
			Required: flag.Required,
			Aliases:  []string{flag.Short},
			Sources:  vsc,
		}
	case "float":
		rflag = &cli.FloatFlag{
			Name:     flag.Name,
			Usage:    flag.Usage,
			Local:    flag.Local,
			Required: flag.Required,
			Aliases:  []string{flag.Short},
			Sources:  vsc,
		}
	default:
		panic("unmatch type flag...\n")
	}
	setFlagValue(rflag, flag.Value.Type, flag.Value.Default)
	return rflag
}

func setFlagValue(flag cli.Flag, t string, value any) {
	if value == nil {
		return
	}
	switch t {
	case "string":
		strFlag, _ := flag.(*cli.StringFlag)
		strFlag.Value = cast.Must[string](cast.ToE[string](value))
	case "int":
		intFlag, _ := flag.(*cli.IntFlag)
		intFlag.Value = cast.Must[int](cast.ToE[int](value))
	case "uint":
		uintFlag, _ := flag.(*cli.UintFlag)
		uintFlag.Value = cast.Must[uint](cast.ToE[uint](value))
	case "float":
		floatFlag, _ := flag.(*cli.FloatFlag)
		floatFlag.Value = cast.Must[float64](cast.ToE[float64](value))
	default:
		panic("unmatch type flag...\n")
	}
}
