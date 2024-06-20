package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
)

var usage = `
usage: vault-k8s-auth-backup [-quiet|--quiet] [-file|--file <vault-k8s-auth-backup-json-file-path>] [<k8s-auth-method-mount-path>]

Usage of vault-k8s-auth-backup:

Flags:

  -file / --file string (Optional)
      vault k8s auth backup json file path (default "vault_k8s_auth_backup.json")

  -quiet / --quiet (Optional)
      quiet progress (default false).
      By default vault-k8s-auth-backup CLI will show a lot of details
      about the backup process and detailed progress during the
      backup process

  -h / -help / --help (Optional)
      show help

Arguments:

  k8s-auth-method-mount-path string (Optional)
      vault k8s auth method mount path.
      If none is given, as it's optional, by default vault-k8s-auth-backup CLI will
      backup all k8s auth methods at different mount paths

examples:

# show help
vault-k8s-auth-backup -h

# show help
vault-k8s-auth-backup --help

# backs up all vault k8s auth methods
vault-k8s-auth-backup

# backs up vault k8s auth method mounted
# at "production/" mount path.
# it will throw an error if it does not exist
vault-k8s-auth-backup production

# quietly backup all vault k8s auth methods.
# this will just show dots (.) for progress
vault-k8s-auth-backup -quiet

# OR you can use --quiet too instead of -quiet

vault-k8s-auth-backup --quiet
`

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "%s", usage)
		os.Exit(0)
	}
	quietProgress := flag.Bool("quiet", false, "quiet progress")
	vaultK8sAuthBackupJsonFileName := flag.String("file", "vault_k8s_auth_backup.json", "vault k8s auth backup json file path")
	flag.Parse()

	if !(flag.NArg() == 1 || flag.NArg() == 0) {
		fmt.Fprintf(os.Stderr, "invalid number of arguments: %d. expected 0 or 1 arguments.\n\n", flag.NArg())
		flag.Usage()
	}

	config := api.DefaultConfig()
	client, err := api.NewClient(config)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating vault client: %s\n", err)
		os.Exit(1)
	}

	k8sAuthMethodMountPaths := []string{}

	if flag.NArg() == 0 {
		allK8sAuthMethodMountPaths, err := getAllK8sAuthMethodMountPaths(client)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error listing all vault k8s auth methods: %s\n", err)
			os.Exit(1)
		}
		k8sAuthMethodMountPaths = append(k8sAuthMethodMountPaths, allK8sAuthMethodMountPaths...)
	} else {
		k8sAuthMethodMountPath := flag.Args()[0]
		k8sAuthMethodMountPaths = append(k8sAuthMethodMountPaths, k8sAuthMethodMountPath)
	}

	vaultK8sAuthBackup := backupK8sAuthMethods(client, k8sAuthMethodMountPaths, *quietProgress)

	vaultK8sAuthBackupJSON, err := convertVaultK8sAuthBackupToJSON(vaultK8sAuthBackup)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error converting vault k8s auth backup to json: %s\n", err)
		os.Exit(1)
	}
	err = writeToFile(vaultK8sAuthBackupJSON, *vaultK8sAuthBackupJsonFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error writing vault k8s auth backup to json file: %s\n", err)
		os.Exit(1)
	}
}
