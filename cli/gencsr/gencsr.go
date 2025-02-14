// Package gencsr implements the gencsr command.
package gencsr

import (
	"encoding/json"
	"errors"

	"github.com/sgorbunov/cfssl/cli"
	"github.com/sgorbunov/cfssl/csr"
	"github.com/sgorbunov/cfssl/helpers"
	"github.com/sgorbunov/cfssl/signer"
)

var gencsrUsageText = `cfssl gencsr -- generate a csr from a private key with existing CSR json specification or certificate

Usage of gencsr:
        cfssl gencsr -key private_key_file [-host hostname_override] CSRJSON
        cfssl gencsr -key private_key_file [-host hostname_override] -cert certificate_file

Arguments:
        CSRJSON:    JSON file containing the request, use '-' for reading JSON from stdin

Flags:
`

var gencsrFlags = []string{"key", "cert"}

func gencsrMain(args []string, c cli.Config) (err error) {
	if c.KeyFile == "" {
		return errors.New("private key file is required through '-key', please check with usage")
	}

	keyBytes, err := helpers.ReadBytes(c.KeyFile)
	if err != nil {
		return err
	}

	key, err := helpers.ParsePrivateKeyPEM(keyBytes)
	if err != nil {
		return err
	}

	// prepare a stub CertificateRequest
	req := &csr.CertificateRequest{
		KeyRequest: csr.NewKeyRequest(),
	}

	if c.CertFile != "" {
		if len(args) > 0 {
			return errors.New("no argument is accepted with '-cert', please check with usage")
		}

		certBytes, err := helpers.ReadBytes(c.CertFile)
		if err != nil {
			return err
		}

		cert, err := helpers.ParseCertificatePEM(certBytes)
		if err != nil {
			return err
		}

		req = csr.ExtractCertificateRequest(cert)
	} else {
		csrFile, args, err := cli.PopFirstArgument(args)
		if err != nil {
			return err
		}

		if len(args) > 0 {
			return errors.New("only one argument is accepted, please check with usage")
		}

		csrFileBytes, err := cli.ReadStdin(csrFile)
		if err != nil {
			return err
		}

		err = json.Unmarshal(csrFileBytes, req)
		if err != nil {
			return err
		}
	}

	if c.Hostname != "" {
		req.Hosts = signer.SplitHosts(c.Hostname)
	}

	csrBytes, err := csr.Generate(key, req)
	if err != nil {
		return err
	}

	cli.PrintCert(keyBytes, csrBytes, nil)
	return nil
}

// Command assembles the definition of Command 'gencsr'
var Command = &cli.Command{UsageText: gencsrUsageText, Flags: gencsrFlags, Main: gencsrMain}
