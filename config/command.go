package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/spf13/cast"
	"github.com/urfave/cli/v3"
)

var (
	command_separator = "."
	_globalCommandMap = make(map[string]*cli.Command)
)

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
	Required    bool      `json:"required"`
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
	if _globalCommandMap == nil {
		panic("command map is nil, please call InitCommand first")
	}
	if cmd, ok := _globalCommandMap[key]; ok {
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

	_globalCommandMap[key] = root
	return root
}

func buildFlag(flag Flag) cli.Flag {
	vsc := cli.NewValueSourceChain()
	vsc.Append(cli.EnvVars(flag.Value.Env...))
	vsc.Append(ConfigVars(flag.Value.Config...))
	switch flag.Value.Type {
	case "string":
		return &cli.StringFlag{
			Name:     flag.Name,
			Usage:    flag.Usage,
			Local:    flag.Local,
			Required: flag.Required,
			Value:    cast.Must[string](cast.ToE[string](flag.Value.Default)),
			Aliases:  []string{flag.Short},
			Sources:  vsc,
		}
	case "int":
		return &cli.IntFlag{
			Name:     flag.Name,
			Usage:    flag.Usage,
			Local:    flag.Local,
			Required: flag.Required,
			Value:    cast.Must[int](cast.ToE[int](flag.Value.Default)),
			Aliases:  []string{flag.Short},
			Sources:  vsc,
		}
	case "uint":
		return &cli.UintFlag{
			Name:     flag.Name,
			Usage:    flag.Usage,
			Local:    flag.Local,
			Required: flag.Required,
			Value:    cast.Must[uint](cast.ToE[uint](flag.Value.Default)),
			Aliases:  []string{flag.Short},
			Sources:  vsc,
		}
	case "bool":
		return &cli.BoolFlag{
			Name:     flag.Name,
			Usage:    flag.Usage,
			Local:    flag.Local,
			Required: flag.Required,
			Value:    cast.Must[bool](cast.ToE[bool](flag.Value.Default)),
			Aliases:  []string{flag.Short},
			Sources:  vsc,
		}
	case "float64":
		return &cli.Float64Flag{
			Name:     flag.Name,
			Usage:    flag.Usage,
			Local:    flag.Local,
			Required: flag.Required,
			Value:    cast.Must[float64](cast.ToE[float64](flag.Value.Default)),
			Aliases:  []string{flag.Short},
			Sources:  vsc,
		}
	default:
		log.Fatalf("unmatch type flag...\n")
		return nil
	}
}
