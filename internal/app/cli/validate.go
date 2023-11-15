package cli

import (
	"fmt"

	"github.com/web-of-things-open-source/tm-catalog-cli/internal"
	"github.com/web-of-things-open-source/tm-catalog-cli/internal/commands/validate"
)

func ValidateFile(filename string) error {

	_, raw, err := internal.ReadRequiredFile(filename)
	if err != nil {
		fmt.Printf("could not read file: %v\n", err)
		return err
	}

	_, err = validate.ValidateThingModel(raw)
	if err != nil {
		fmt.Printf("validation error: %v\n", err)
		return err
	}
	fmt.Printf("validated successfully: %s\n", filename)
	return nil
}
