package main

import (
	"EtherDrop/core"
	"EtherDrop/tools"
	"flag"
	"strconv"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

func main() {
	defer tools.ExitsRecover()

	config.AddDriver(yaml.Driver)

	err := config.LoadFiles("configs/config.yml")
	if err != nil {
		panic(err)
	}

	tools.PrintLogo()

	var selectedTools int

	flagArg := flag.Int("action", 0, "Input Choice With Flag -action, 1 = Auto Farming (Unlimited Loop), 2 = Get Ref Code 3 = Register Account With Ref Code")
	tools.Logger("1", "Auto Farming (Unlimited Loop)")
	tools.Logger("2", "Get Ref Code")
	tools.Logger("3", "Register Account With Ref Code")

	flag.Parse()

	if *flagArg > 3 {
		tools.Logger("error", "Invalid Flag Choice")
	} else if *flagArg != 0 {
		selectedTools = *flagArg
	} else {
		selectedTools, _ = strconv.Atoi(tools.InputTerminal("Select Tools: "))
		if selectedTools <= 0 || selectedTools > 3 {
			tools.Logger("error", "Invalid Choice")
			return
		}
	}

	core.LaunchBot(selectedTools)
}
