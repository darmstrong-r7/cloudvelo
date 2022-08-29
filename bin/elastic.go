package main

import (
	"fmt"

	"www.velocidex.com/golang/cloudvelo/schema"
	"www.velocidex.com/golang/cloudvelo/services"
)

var (
	elastic_command = app.Command(
		"elastic", "Manipulate the elastic datastore")

	elastic_command_reset = elastic_command.Command(
		"reset", "Drop all the indexes and recreate them")

	elastic_command_reset_org_id = elastic_command.Flag(
		"org_id", "An OrgID to initialize").String()

	elastic_command_reset_filter = elastic_command_reset.Arg(
		"index_filter", "If specified only re-create these indexes").String()
)

func doResetElastic() error {
	config_obj, err := makeDefaultConfigLoader().
		LoadAndValidate()
	if err != nil {
		return fmt.Errorf("loading config file: %w", err)
	}

	ctx, cancel := install_sig_handler()
	defer cancel()

	err = services.StartElasticSearchService(config_obj, *elastic_config)
	if err != nil {
		return err
	}

	return schema.Initialize(ctx,
		*elastic_command_reset_org_id,
		*elastic_command_reset_filter, true /* reset */)
}

func init() {
	command_handlers = append(command_handlers, func(command string) bool {
		if command == elastic_command_reset.FullCommand() {
			FatalIfError(elastic_command_reset, doResetElastic)
			return true
		}
		return false
	})
}
