// Package revoke implements the revoke command.
package revoke

import (
	"errors"

	"github.com/sgorbunov/cfssl/certdb/dbconf"
	"github.com/sgorbunov/cfssl/certdb/sql"
	"github.com/sgorbunov/cfssl/cli"
	"github.com/sgorbunov/cfssl/log"
	"github.com/sgorbunov/cfssl/ocsp"
)

var revokeUsageTxt = `cfssl revoke -- revoke a certificate in the certificate store

Usage:

Revoke a certificate:
	   cfssl revoke -db-config config_file -serial serial -aki authority_key_id [-reason reason]

Reason can be an integer code or a string in ReasonFlags in RFC 5280

Flags:
`

var revokeFlags = []string{"serial", "reason"}

func revokeMain(args []string, c cli.Config) error {
	if len(args) > 0 {
		return errors.New("argument is provided but not defined; please refer to the usage by flag -h")
	}

	if len(c.Serial) == 0 {
		return errors.New("serial number is required but not provided")
	}

	if len(c.AKI) == 0 {
		return errors.New("authority key id is required but not provided")
	}

	if c.DBConfigFile == "" {
		return errors.New("need DB config file (provide with -db-config)")
	}

	db, err := dbconf.DBFromConfig(c.DBConfigFile)
	if err != nil {
		return err
	}

	dbAccessor := sql.NewAccessor(db)

	reasonCode, err := ocsp.ReasonStringToCode(c.Reason)
	if err != nil {
		log.Error("Invalid reason code: ", err)
		return err
	}

	return dbAccessor.RevokeCertificate(c.Serial, c.AKI, reasonCode)
}

// Command assembles the definition of Command 'revoke'
var Command = &cli.Command{
	UsageText: revokeUsageTxt,
	Flags:     revokeFlags,
	Main:      revokeMain,
}
