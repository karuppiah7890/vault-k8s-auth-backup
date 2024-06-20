# vault-k8s-auth-backup

Using this CLI tool, you can backup Vault Kubernetes (K8s) Auth Method(s) from a Vault instance to your local machine! :D

Note: The tool is written in Golang and uses the Vault Official Golang Client API. The Official Vault Golang Client API documentation is accessible here - https://pkg.go.dev/github.com/hashicorp/vault/api or directly in the source code of Vault as comments. Source code - https://github.com/hashicorp/vault. Look for the API directory `api` in the source code which has all the Client API code and docs / documentation, for example, latest Golang Client API docs can be found here - https://github.com/hashicorp/vault/tree/main/api

Note: The tool needs Vault credentials of a user/account that has access to Vault, to read and list (if needed) the Vault Kubernetes (K8s) Auth Method(s) that you want to backup. Look at [Authorization Details for the Vault Credentials](#authorization-details-for-the-vault-credentials) for more details

Note: We have tested this only with some versions of Vault (like v1.17.x). So beware to test this in a testing environment with whatever version of Vault you are using, before using this in critical environments like production! Also, ensure that the testing environment is as close to your production environment as possible so that your testing makes sense

## Building

```bash
CGO_ENABLED=0 go build -v
```

or

```bash
make
```

## Authorization Details for the Vault Credentials

As mentioned before in a note, the tool needs Vault credentials of a user/account that has access to Vault, to read and list the Vault Kubernetes (K8s) Auth Method(s) that you want to backup

Access to list the Vault Kubernetes (K8s) Auth Methods is required ONLY if you want to backup all the Vault Kubernetes (K8s) Auth Methods in one go, in one command, without specifying the name of any Kubernetes (K8s) Auth Method since you want to backup all of them. By the way, you actually need access to list all the Vault Auth Methods in general, to be able to list the Vault Kubernetes (K8s) Auth Methods too. The tool lists all the Vault Auth Methods first and then finds which Vault Auth Methods are of type `kubernetes` and then understands the list of Vault Kubernetes (K8s) Auth Methods present in the Vault instance. So, that's how the tool works

An example liberal / relaxed Vault Policy that could be used to provide access to backup all the config and roles in a Vault Kubernetes (K8s) Auth Method is -

```hcl
# Vault Kubernetes (K8s) Auth Method's mount path is "kubernetes"
path "auth/kubernetes/*" {
  capabilities = ["read", "list"]
}

path "sys/auth" {
  capabilities = ["read"]
}
```

OR

```hcl
# Vault Kubernetes (K8s) Auth Method's mount path is "kubernetes"
path "/auth/kubernetes/*" {
  capabilities = ["read", "list"]
}

path "/sys/auth" {
  capabilities = ["read"]
}
```

You can use a similar Vault Policy based on the mount path of the Vault Kubernetes (K8s) Auth Method(s) that you are using and want to backup. You can create a Vault Token that has this Vault Policy attached to it and use that token to backup the Vault Kubernetes (K8s) Auth Method(s) using the `vault-k8s-auth-backup` tool :)

The above policy is required because the tool will access the following paths with the given operation

`Read` - `<vault-http-api-url>/v1/sys/auth` - To list all the Vault Auth Methods

`Read` - `<vault-http-api-url>/v1/auth/<kubernetes-auth-method-mount-path>/config` - To get the given Vault Kubernetes (K8s) Auth Method's Configuration / Config

`List` - `<vault-http-api-url>/v1/auth/<kubernetes-auth-method-mount-path>/role` - To list the given Vault Kubernetes (K8s) Auth Method's Roles, that is, to list the Names of all the Roles available in the given Vault Kubernetes (K8s) Auth Method

`Read` - `<vault-http-api-url>/v1/auth/<kubernetes-auth-method-mount-path>/role/<role-name>` - To get the given Vault Kubernetes (K8s) Auth Method's particular Role's Configuration / Config.

Here, `Read` is the same as `GET` in HTTP API / Web API. In Vault CLI, it's basically `vault read` CLI command

Here, `List` is the same as `GET` in HTTP API / Web API with `list=true` query parameter. In Vault CLI, it's basically `vault list` CLI command

As you can see, the tool accesses a lot of paths to get different information as each information is available at a specific path. For the tool to work well and perfectly, it needs all of these access. You can ignore the Vault Auth Methods List access alone in case you don't want to backup all Auth Methods in one go in one command

## Usage

```bash
$ ./vault-k8s-auth-backup --help

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
```

# Demo

I created a new dummy local Vault instance in developer mode for this demo. I ran the Vault server like this -

```bash
vault server -dev -dev-root-token-id root -dev-listen-address 127.0.0.1:8200
```

I'm going to create and configure some dummy Vault Kubernetes Auth Methods, which is mounted at `perf/`, `staging/`, `production/` and `dummy-test/`. I'll be using the Vault CLI to do this but you can do it in any way you want. I'll be using the `root` user's Vault API token to create the Kubernetes Auth Methods, but it's not necessary to use the root token, you can use any token with less privileges too, following the Principle of Least Privilege. Just ensure that the token has enough access to create and configure Vault Kubernetes Auth Methods. Also ensure that your token is safe and secure, regardless of it being `root` user's token or not

Note: These are just dummy auth methods / test auth methods. These are just for demo purposes.

Initially the Vault looks like this -

```bash
$ export VAULT_ADDR='http://127.0.0.1:8200'
$ export VAULT_TOKEN="root"

$ vault status
Key             Value
---             -----
Seal Type       shamir
Initialized     true
Sealed          false
Total Shares    1
Threshold       1
Version         1.17.0
Build Date      2024-06-10T10:11:34Z
Storage Type    inmem
Cluster Name    vault-cluster-5058b87d
Cluster ID      df0054c9-5b96-2d5e-10c4-4e43eae55e35
HA Enabled      false

$ vault auth list
Path      Type     Accessor               Description                Version
----      ----     --------               -----------                -------
token/    token    auth_token_91fc8379    token based credentials    n/a
```

Let's put in some data ;) :D

First, let's enable some Kubernetes Auth Methods at different Mount Paths

```bash
export VAULT_ADDR="http://127.0.0.1:8200"
export VAULT_TOKEN="root"

# For perf testing / performance testing environment kubernetes cluster
vault auth enable -path=perf kubernetes

# For staging environment kubernetes cluster
vault auth enable -path=staging kubernetes

# For production environment kubernetes cluster
vault auth enable -path=production kubernetes

# For a dummy test kubernetes cluster
vault auth enable -path=dummy-test kubernetes
```

When you run the above, you will find output like this -

```bash
$ export VAULT_ADDR="http://127.0.0.1:8200"
$ export VAULT_TOKEN="root"

$ # For perf testing / performance testing environment kubernetes cluster
$ vault auth enable -path=perf kubernetes
Success! Enabled kubernetes auth method at: perf/

$ # For staging environment kubernetes cluster
$ vault auth enable -path=staging kubernetes
Success! Enabled kubernetes auth method at: staging/

$ # For production environment kubernetes cluster
$ vault auth enable -path=production kubernetes
Success! Enabled kubernetes auth method at: production/

$ # For a dummy test kubernetes cluster
$ vault auth enable -path=dummy-test kubernetes
Success! Enabled kubernetes auth method at: dummy-test/
```

Now, let's create some Kubernetes Auth Method Config and Roles at these different Mount Paths. Only for one of these Kubernetes Auth Method's, I won't configure Config and Roles for it - it will be the Kubernetes Auth Method mounted at mount path `dummy-test/`

```bash
export VAULT_ADDR="http://127.0.0.1:8200"
export VAULT_TOKEN="root"

cat /Users/karuppiah.n/every-day-log/example-ca-cert.cert

vault write -non-interactive \
   /auth/perf/config \
   disable_local_ca_jwt=true \
   kubernetes_host=https://dummy.my-cluster.com \
   kubernetes_ca_cert=@/Users/karuppiah.n/every-day-log/example-ca-cert.cert \
   token_reviewer_jwt=dummy_jwt_token

vault read -non-interactive \
  -format json \
  auth/perf/config

vault read -non-interactive \
  -format json \
  /auth/perf/config | jq

vault write -non-interactive \
  auth/perf/role/dummy \
  alias_name_source=serviceaccount_uid \
  audience=dummy \
  bound_service_account_names=dummy-service-account \
  bound_service_account_namespace_selector='{"matchLabels":{"environment":"perf","org":"company1"}}' \
  bound_service_account_namespaces=dummy-namespace \
  token_bound_cidrs=192.168.1.0/24,10.0.1.0/24 \
  token_explicit_max_ttl=7200s \
  token_max_ttl=7200s \
  token_no_default_policy=true \
  token_num_uses=10 \
  token_period=3600.1s \
  token_policies=dummy-policy \
  token_ttl=1h \
  token_type=default

# yes, 3600.1s is a valid value :O :D It will be taken as 3600s by Vault though

# Note that here using `/auth/perf` is the same as `auth/perf`.
# Both end up using `<vault-http-api-url>/v1/auth/perf` as the complete URL

vault list -non-interactive \
  auth/perf/role

vault list -non-interactive \
  -format json \
  /auth/perf/role

vault list -non-interactive \
  -format json \
  /auth/perf/role | jq

vault read -non-interactive \
  -format json \
  auth/perf/role/dummy

vault read -non-interactive \
  -format json \
  /auth/perf/role/dummy | jq

####

vault write -non-interactive \
   /auth/staging/config \
   disable_local_ca_jwt=true \
   kubernetes_host=https://dummy.my-staging-cluster.com \
   kubernetes_ca_cert=@/Users/karuppiah.n/every-day-log/example-ca-cert.cert \
   token_reviewer_jwt=dummy_jwt_token_2

vault read -non-interactive \
  -format json \
  auth/staging/config

vault read -non-interactive \
  -format json \
  /auth/staging/config | jq

####

vault write -non-interactive \
   /auth/production/config \
   disable_local_ca_jwt=true \
   kubernetes_host=https://dummy.my-prod-cluster.com \
   kubernetes_ca_cert=@/Users/karuppiah.n/every-day-log/example-ca-cert.cert \
   token_reviewer_jwt=dummy_jwt_token_3

vault read -non-interactive \
  -format json \
  auth/production/config

vault read -non-interactive \
  -format json \
  /auth/production/config | jq

vault write -non-interactive \
  /auth/production/role/dummy \
  alias_name_source=serviceaccount_uid \
  audience=dummy \
  bound_service_account_names=dummy-service-account \
  bound_service_account_namespace_selector='{"matchLabels":{"environment":"prod","org":"company-1"}}' \
  bound_service_account_namespaces=dummy-namespace \
  token_bound_cidrs=192.168.1.0/24 \
  token_bound_cidrs=10.0.1.0/24 \
  token_bound_cidrs=10.0.2.0/24 \
  token_explicit_max_ttl=1200s \
  token_max_ttl=1200s \
  token_no_default_policy=true \
  token_num_uses=10 \
  token_period=600.9s \
  token_policies=dummy-policy \
  token_ttl=10m \
  token_type=default

# yes, 600.9s is a valid value :O :D It will be taken as 600s by Vault though

vault list -non-interactive \
  auth/production/role

vault list -non-interactive \
  -format json \
  /auth/production/role

vault list -non-interactive \
  -format json \
  /auth/production/role | jq

vault read -non-interactive \
  -format json \
  auth/production/role/dummy

vault read -non-interactive \
  -format json \
  /auth/production/role/dummy | jq

vault write -non-interactive \
  /auth/production/role/dummy-2 \
  alias_name_source=serviceaccount_uid \
  audience=dummy-2 \
  bound_service_account_names=dummy-2-service-account \
  bound_service_account_namespace_selector='{"matchLabels":{"environment":"prod","org":"company-1"}}' \
  bound_service_account_namespaces=dummy-2-namespace \
  token_bound_cidrs=192.168.1.0/24 \
  token_bound_cidrs=10.0.1.0/24 \
  token_bound_cidrs=10.0.2.0/24 \
  token_explicit_max_ttl=1200s \
  token_max_ttl=1200s \
  token_no_default_policy=true \
  token_num_uses=10 \
  token_period=600.9s \
  token_policies=dummy-2-policy \
  token_ttl=10m \
  token_type=default

vault list -non-interactive \
  auth/production/role

vault list -non-interactive \
  -format json \
  /auth/production/role

vault list -non-interactive \
  -format json \
  /auth/production/role | jq

vault read -non-interactive \
  -format json \
  auth/production/role/dummy-2

vault read -non-interactive \
  -format json \
  /auth/production/role/dummy-2 | jq
```

I'm using [`jq` tool](https://jqlang.github.io/jq/) here for colorful JSON output. [`jq` tool](https://jqlang.github.io/jq/) is pretty cool for JSON stuff - JSON prettifying, querying, modification etc and is very powerful. I'm just using it here for seeing the output in color. You can find the full [`jq` manual](https://jqlang.github.io/jq/manual/) here - https://jqlang.github.io/jq/manual/

When you run the above, you will find output like this -

```bash
$ export VAULT_ADDR="http://127.0.0.1:8200"
$ export VAULT_TOKEN="root"


$ cat /Users/karuppiah.n/every-day-log/example-ca-cert.cert
-----BEGIN CERTIFICATE-----
MIIDBTCCAe2gAwIBAgIIGkf4cKDZKbgwDQYJKoZIhvcNAQELBQAwFTETMBEGA1UE
AxMKa3ViZXJuZXRlczAeFw0yNDA0MDMwNzQ5MjhaFw0zNDA0MDEwNzU0MjhaMBUx
EzARBgNVBAMTCmt1YmVybmV0ZXMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK
AoIBAQDf8V5Mxa5OAtZckVZo4RslyD1zgtpnkLCJYIPqD23U1cmmF75VV04httOE
sPHv610WhkOje2jNEZO4SY0wJi6A9QVOJyyCfXAzehY4IYZCWWbfFL99dg28WH7N
tEpU648GUB9M8Sd/sngof3/CRfi0OELKejmn3xmEYV74Vj3hB57KC8dvNpU0Zgs1
62oF/ZXXMWLOugM8WonekIwpjy71b3VfRatgBCcqr5yQvyR3r9MJjxgQG/eJAZHI
UCiIF4GKFsRdCl7hSl+MRf4beg5N2Qc/FeomxxFD8Mc7guYaA5errLlapdHYd0Kx
VHRwLQ/+hnmnP5FINV+kj7k782c1AgMBAAGjWTBXMA4GA1UdDwEB/wQEAwICpDAP
BgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBSCGsrJ79HjbShkkJQrjmOtl9BX4DAV
BgNVHREEDjAMggprdWJlcm5ldGVzMA0GCSqGSIb3DQEBCwUAA4IBAQAbC0IJLDVZ
sYK5sOZi4Z0WwdLKHn5jTF0cS+6E4gWO/3qZMH1lQELlvUa6B3rCOvJ12/a4MoWL
9JGzRNqi1G9Nox93OW9MrJsfRN+a6HB7cq0qMPhCRv+h7KBurN+MRZu0AZuSWJ5G
BJ3eIrIFQbBpHtho2Ueu4JYlifJIEmn5yWNvIHYCumEevXPB5dEASXE7djywteE+
w3Pi64gYnj3Tb3T8ZIFyWsBqdWZzeFPDUasVi/IuFY/7plDuIOY27BDhhvX2TirH
OzwEMN9nZ9PWaSRyeHLSslFTjCndoVO90Y95UbBTjz2YO/nueG+4UN08ApqSffzk
1kpyPWcIU0/8
-----END CERTIFICATE-----


$ vault write -non-interactive \
    /auth/perf/config \
    disable_local_ca_jwt=true \
    kubernetes_host=https://dummy.my-cluster.com \
    kubernetes_ca_cert=@/Users/karuppiah.n/every-day-log/example-ca-cert.cert \
    token_reviewer_jwt=dummy_jwt_token
Success! Data written to: auth/perf/config


$ vault read -non-interactive \
   -format json \
   auth/perf/config
{
  "request_id": "3e8f2ae5-8046-398c-e608-9c76b7d26efd",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "disable_iss_validation": true,
    "disable_local_ca_jwt": true,
    "issuer": "",
    "kubernetes_ca_cert": "-----BEGIN CERTIFICATE-----\nMIIDBTCCAe2gAwIBAgIIGkf4cKDZKbgwDQYJKoZIhvcNAQELBQAwFTETMBEGA1UE\nAxMKa3ViZXJuZXRlczAeFw0yNDA0MDMwNzQ5MjhaFw0zNDA0MDEwNzU0MjhaMBUx\nEzARBgNVBAMTCmt1YmVybmV0ZXMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK\nAoIBAQDf8V5Mxa5OAtZckVZo4RslyD1zgtpnkLCJYIPqD23U1cmmF75VV04httOE\nsPHv610WhkOje2jNEZO4SY0wJi6A9QVOJyyCfXAzehY4IYZCWWbfFL99dg28WH7N\ntEpU648GUB9M8Sd/sngof3/CRfi0OELKejmn3xmEYV74Vj3hB57KC8dvNpU0Zgs1\n62oF/ZXXMWLOugM8WonekIwpjy71b3VfRatgBCcqr5yQvyR3r9MJjxgQG/eJAZHI\nUCiIF4GKFsRdCl7hSl+MRf4beg5N2Qc/FeomxxFD8Mc7guYaA5errLlapdHYd0Kx\nVHRwLQ/+hnmnP5FINV+kj7k782c1AgMBAAGjWTBXMA4GA1UdDwEB/wQEAwICpDAP\nBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBSCGsrJ79HjbShkkJQrjmOtl9BX4DAV\nBgNVHREEDjAMggprdWJlcm5ldGVzMA0GCSqGSIb3DQEBCwUAA4IBAQAbC0IJLDVZ\nsYK5sOZi4Z0WwdLKHn5jTF0cS+6E4gWO/3qZMH1lQELlvUa6B3rCOvJ12/a4MoWL\n9JGzRNqi1G9Nox93OW9MrJsfRN+a6HB7cq0qMPhCRv+h7KBurN+MRZu0AZuSWJ5G\nBJ3eIrIFQbBpHtho2Ueu4JYlifJIEmn5yWNvIHYCumEevXPB5dEASXE7djywteE+\nw3Pi64gYnj3Tb3T8ZIFyWsBqdWZzeFPDUasVi/IuFY/7plDuIOY27BDhhvX2TirH\nOzwEMN9nZ9PWaSRyeHLSslFTjCndoVO90Y95UbBTjz2YO/nueG+4UN08ApqSffzk\n1kpyPWcIU0/8\n-----END CERTIFICATE-----\n",
    "kubernetes_host": "https://dummy.my-cluster.com",
    "pem_keys": [],
    "token_reviewer_jwt_set": true,
    "use_annotations_as_alias_metadata": false
  },
  "warnings": null,
  "mount_type": "kubernetes"
}


$ vault read -non-interactive \
   -format json \
   /auth/perf/config | jq
{
  "request_id": "b8bcaea7-8507-c9be-8403-d2d961457d3a",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "disable_iss_validation": true,
    "disable_local_ca_jwt": true,
    "issuer": "",
    "kubernetes_ca_cert": "-----BEGIN CERTIFICATE-----\nMIIDBTCCAe2gAwIBAgIIGkf4cKDZKbgwDQYJKoZIhvcNAQELBQAwFTETMBEGA1UE\nAxMKa3ViZXJuZXRlczAeFw0yNDA0MDMwNzQ5MjhaFw0zNDA0MDEwNzU0MjhaMBUx\nEzARBgNVBAMTCmt1YmVybmV0ZXMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK\nAoIBAQDf8V5Mxa5OAtZckVZo4RslyD1zgtpnkLCJYIPqD23U1cmmF75VV04httOE\nsPHv610WhkOje2jNEZO4SY0wJi6A9QVOJyyCfXAzehY4IYZCWWbfFL99dg28WH7N\ntEpU648GUB9M8Sd/sngof3/CRfi0OELKejmn3xmEYV74Vj3hB57KC8dvNpU0Zgs1\n62oF/ZXXMWLOugM8WonekIwpjy71b3VfRatgBCcqr5yQvyR3r9MJjxgQG/eJAZHI\nUCiIF4GKFsRdCl7hSl+MRf4beg5N2Qc/FeomxxFD8Mc7guYaA5errLlapdHYd0Kx\nVHRwLQ/+hnmnP5FINV+kj7k782c1AgMBAAGjWTBXMA4GA1UdDwEB/wQEAwICpDAP\nBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBSCGsrJ79HjbShkkJQrjmOtl9BX4DAV\nBgNVHREEDjAMggprdWJlcm5ldGVzMA0GCSqGSIb3DQEBCwUAA4IBAQAbC0IJLDVZ\nsYK5sOZi4Z0WwdLKHn5jTF0cS+6E4gWO/3qZMH1lQELlvUa6B3rCOvJ12/a4MoWL\n9JGzRNqi1G9Nox93OW9MrJsfRN+a6HB7cq0qMPhCRv+h7KBurN+MRZu0AZuSWJ5G\nBJ3eIrIFQbBpHtho2Ueu4JYlifJIEmn5yWNvIHYCumEevXPB5dEASXE7djywteE+\nw3Pi64gYnj3Tb3T8ZIFyWsBqdWZzeFPDUasVi/IuFY/7plDuIOY27BDhhvX2TirH\nOzwEMN9nZ9PWaSRyeHLSslFTjCndoVO90Y95UbBTjz2YO/nueG+4UN08ApqSffzk\n1kpyPWcIU0/8\n-----END CERTIFICATE-----\n",
    "kubernetes_host": "https://dummy.my-cluster.com",
    "pem_keys": [],
    "token_reviewer_jwt_set": true,
    "use_annotations_as_alias_metadata": false
  },
  "warnings": null,
  "mount_type": "kubernetes"
}


$ vault write -non-interactive \
   auth/perf/role/dummy \
   alias_name_source=serviceaccount_uid \
   audience=dummy \
   bound_service_account_names=dummy-service-account \
   bound_service_account_namespace_selector='{"matchLabels":{"environment":"perf","org":"company1"}}' \
   bound_service_account_namespaces=dummy-namespace \
   token_bound_cidrs=192.168.1.0/24,10.0.1.0/24 \
   token_explicit_max_ttl=7200s \
   token_max_ttl=7200s \
   token_no_default_policy=true \
   token_num_uses=10 \
   token_period=3600.1s \
   token_policies=dummy-policy \
   token_ttl=1h \
   token_type=default
Success! Data written to: auth/perf/role/dummy


$ # yes, 3600.1s is a valid value :O :D It will be taken as 3600s by Vault though


$ # Note that here using `/auth/perf` is the same as `auth/perf`.
$ # Both end up using `<vault-http-api-url>/v1/auth/perf` as the complete URL

$ vault list -non-interactive \
   auth/perf/role
Keys
----
dummy


$ vault list -non-interactive \
   -format json \
   /auth/perf/role
[
  "dummy"
]


$ vault list -non-interactive \
   -format json \
   /auth/perf/role | jq
[
  "dummy"
]

$ vault read -non-interactive \
   -format json \
   auth/perf/role/dummy
{
  "request_id": "e36d22bc-331e-2c0b-f9e3-10787f0baa66",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "alias_name_source": "serviceaccount_uid",
    "audience": "dummy",
    "bound_service_account_names": [
      "dummy-service-account"
    ],
    "bound_service_account_namespace_selector": "{\"matchLabels\":{\"environment\":\"perf\",\"org\":\"company1\"}}",
    "bound_service_account_namespaces": [
      "dummy-namespace"
    ],
    "token_bound_cidrs": [
      "192.168.1.0/24",
      "10.0.1.0/24"
    ],
    "token_explicit_max_ttl": 7200,
    "token_max_ttl": 7200,
    "token_no_default_policy": true,
    "token_num_uses": 10,
    "token_period": 3600,
    "token_policies": [
      "dummy-policy"
    ],
    "token_ttl": 3600,
    "token_type": "default"
  },
  "warnings": null,
  "mount_type": "kubernetes"
}


$ vault read -non-interactive \
   -format json \
   /auth/perf/role/dummy | jq
{
  "request_id": "1a254867-7401-2a86-7c5c-c0df30e2db60",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "alias_name_source": "serviceaccount_uid",
    "audience": "dummy",
    "bound_service_account_names": [
      "dummy-service-account"
    ],
    "bound_service_account_namespace_selector": "{\"matchLabels\":{\"environment\":\"perf\",\"org\":\"company1\"}}",
    "bound_service_account_namespaces": [
      "dummy-namespace"
    ],
    "token_bound_cidrs": [
      "192.168.1.0/24",
      "10.0.1.0/24"
    ],
    "token_explicit_max_ttl": 7200,
    "token_max_ttl": 7200,
    "token_no_default_policy": true,
    "token_num_uses": 10,
    "token_period": 3600,
    "token_policies": [
      "dummy-policy"
    ],
    "token_ttl": 3600,
    "token_type": "default"
  },
  "warnings": null,
  "mount_type": "kubernetes"
}


$ vault write -non-interactive \
    /auth/staging/config \
    disable_local_ca_jwt=true \
    kubernetes_host=https://dummy.my-staging-cluster.com \
    kubernetes_ca_cert=@/Users/karuppiah.n/every-day-log/example-ca-cert.cert \
    token_reviewer_jwt=dummy_jwt_token_2
Success! Data written to: auth/staging/config


$ vault read -non-interactive \
   -format json \
   auth/staging/config
{
  "request_id": "c8018ee5-6c5d-fac3-2a54-ca767af521cd",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "disable_iss_validation": true,
    "disable_local_ca_jwt": true,
    "issuer": "",
    "kubernetes_ca_cert": "-----BEGIN CERTIFICATE-----\nMIIDBTCCAe2gAwIBAgIIGkf4cKDZKbgwDQYJKoZIhvcNAQELBQAwFTETMBEGA1UE\nAxMKa3ViZXJuZXRlczAeFw0yNDA0MDMwNzQ5MjhaFw0zNDA0MDEwNzU0MjhaMBUx\nEzARBgNVBAMTCmt1YmVybmV0ZXMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK\nAoIBAQDf8V5Mxa5OAtZckVZo4RslyD1zgtpnkLCJYIPqD23U1cmmF75VV04httOE\nsPHv610WhkOje2jNEZO4SY0wJi6A9QVOJyyCfXAzehY4IYZCWWbfFL99dg28WH7N\ntEpU648GUB9M8Sd/sngof3/CRfi0OELKejmn3xmEYV74Vj3hB57KC8dvNpU0Zgs1\n62oF/ZXXMWLOugM8WonekIwpjy71b3VfRatgBCcqr5yQvyR3r9MJjxgQG/eJAZHI\nUCiIF4GKFsRdCl7hSl+MRf4beg5N2Qc/FeomxxFD8Mc7guYaA5errLlapdHYd0Kx\nVHRwLQ/+hnmnP5FINV+kj7k782c1AgMBAAGjWTBXMA4GA1UdDwEB/wQEAwICpDAP\nBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBSCGsrJ79HjbShkkJQrjmOtl9BX4DAV\nBgNVHREEDjAMggprdWJlcm5ldGVzMA0GCSqGSIb3DQEBCwUAA4IBAQAbC0IJLDVZ\nsYK5sOZi4Z0WwdLKHn5jTF0cS+6E4gWO/3qZMH1lQELlvUa6B3rCOvJ12/a4MoWL\n9JGzRNqi1G9Nox93OW9MrJsfRN+a6HB7cq0qMPhCRv+h7KBurN+MRZu0AZuSWJ5G\nBJ3eIrIFQbBpHtho2Ueu4JYlifJIEmn5yWNvIHYCumEevXPB5dEASXE7djywteE+\nw3Pi64gYnj3Tb3T8ZIFyWsBqdWZzeFPDUasVi/IuFY/7plDuIOY27BDhhvX2TirH\nOzwEMN9nZ9PWaSRyeHLSslFTjCndoVO90Y95UbBTjz2YO/nueG+4UN08ApqSffzk\n1kpyPWcIU0/8\n-----END CERTIFICATE-----\n",
    "kubernetes_host": "https://dummy.my-staging-cluster.com",
    "pem_keys": [],
    "token_reviewer_jwt_set": true,
    "use_annotations_as_alias_metadata": false
  },
  "warnings": null,
  "mount_type": "kubernetes"
}


$ vault read -non-interactive \
   -format json \
   /auth/staging/config | jq
{
  "request_id": "8b7f10ef-186b-8c92-19f8-427896080514",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "disable_iss_validation": true,
    "disable_local_ca_jwt": true,
    "issuer": "",
    "kubernetes_ca_cert": "-----BEGIN CERTIFICATE-----\nMIIDBTCCAe2gAwIBAgIIGkf4cKDZKbgwDQYJKoZIhvcNAQELBQAwFTETMBEGA1UE\nAxMKa3ViZXJuZXRlczAeFw0yNDA0MDMwNzQ5MjhaFw0zNDA0MDEwNzU0MjhaMBUx\nEzARBgNVBAMTCmt1YmVybmV0ZXMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK\nAoIBAQDf8V5Mxa5OAtZckVZo4RslyD1zgtpnkLCJYIPqD23U1cmmF75VV04httOE\nsPHv610WhkOje2jNEZO4SY0wJi6A9QVOJyyCfXAzehY4IYZCWWbfFL99dg28WH7N\ntEpU648GUB9M8Sd/sngof3/CRfi0OELKejmn3xmEYV74Vj3hB57KC8dvNpU0Zgs1\n62oF/ZXXMWLOugM8WonekIwpjy71b3VfRatgBCcqr5yQvyR3r9MJjxgQG/eJAZHI\nUCiIF4GKFsRdCl7hSl+MRf4beg5N2Qc/FeomxxFD8Mc7guYaA5errLlapdHYd0Kx\nVHRwLQ/+hnmnP5FINV+kj7k782c1AgMBAAGjWTBXMA4GA1UdDwEB/wQEAwICpDAP\nBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBSCGsrJ79HjbShkkJQrjmOtl9BX4DAV\nBgNVHREEDjAMggprdWJlcm5ldGVzMA0GCSqGSIb3DQEBCwUAA4IBAQAbC0IJLDVZ\nsYK5sOZi4Z0WwdLKHn5jTF0cS+6E4gWO/3qZMH1lQELlvUa6B3rCOvJ12/a4MoWL\n9JGzRNqi1G9Nox93OW9MrJsfRN+a6HB7cq0qMPhCRv+h7KBurN+MRZu0AZuSWJ5G\nBJ3eIrIFQbBpHtho2Ueu4JYlifJIEmn5yWNvIHYCumEevXPB5dEASXE7djywteE+\nw3Pi64gYnj3Tb3T8ZIFyWsBqdWZzeFPDUasVi/IuFY/7plDuIOY27BDhhvX2TirH\nOzwEMN9nZ9PWaSRyeHLSslFTjCndoVO90Y95UbBTjz2YO/nueG+4UN08ApqSffzk\n1kpyPWcIU0/8\n-----END CERTIFICATE-----\n",
    "kubernetes_host": "https://dummy.my-staging-cluster.com",
    "pem_keys": [],
    "token_reviewer_jwt_set": true,
    "use_annotations_as_alias_metadata": false
  },
  "warnings": null,
  "mount_type": "kubernetes"
}


$ vault write -non-interactive \
    /auth/production/config \
    disable_local_ca_jwt=true \
    kubernetes_host=https://dummy.my-prod-cluster.com \
    kubernetes_ca_cert=@/Users/karuppiah.n/every-day-log/example-ca-cert.cert \
    token_reviewer_jwt=dummy_jwt_token_3
Success! Data written to: auth/production/config


$ vault read -non-interactive \
   -format json \
   auth/production/config
{
  "request_id": "b59f6c28-9403-1025-b180-9e730c28d5d6",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "disable_iss_validation": true,
    "disable_local_ca_jwt": true,
    "issuer": "",
    "kubernetes_ca_cert": "-----BEGIN CERTIFICATE-----\nMIIDBTCCAe2gAwIBAgIIGkf4cKDZKbgwDQYJKoZIhvcNAQELBQAwFTETMBEGA1UE\nAxMKa3ViZXJuZXRlczAeFw0yNDA0MDMwNzQ5MjhaFw0zNDA0MDEwNzU0MjhaMBUx\nEzARBgNVBAMTCmt1YmVybmV0ZXMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK\nAoIBAQDf8V5Mxa5OAtZckVZo4RslyD1zgtpnkLCJYIPqD23U1cmmF75VV04httOE\nsPHv610WhkOje2jNEZO4SY0wJi6A9QVOJyyCfXAzehY4IYZCWWbfFL99dg28WH7N\ntEpU648GUB9M8Sd/sngof3/CRfi0OELKejmn3xmEYV74Vj3hB57KC8dvNpU0Zgs1\n62oF/ZXXMWLOugM8WonekIwpjy71b3VfRatgBCcqr5yQvyR3r9MJjxgQG/eJAZHI\nUCiIF4GKFsRdCl7hSl+MRf4beg5N2Qc/FeomxxFD8Mc7guYaA5errLlapdHYd0Kx\nVHRwLQ/+hnmnP5FINV+kj7k782c1AgMBAAGjWTBXMA4GA1UdDwEB/wQEAwICpDAP\nBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBSCGsrJ79HjbShkkJQrjmOtl9BX4DAV\nBgNVHREEDjAMggprdWJlcm5ldGVzMA0GCSqGSIb3DQEBCwUAA4IBAQAbC0IJLDVZ\nsYK5sOZi4Z0WwdLKHn5jTF0cS+6E4gWO/3qZMH1lQELlvUa6B3rCOvJ12/a4MoWL\n9JGzRNqi1G9Nox93OW9MrJsfRN+a6HB7cq0qMPhCRv+h7KBurN+MRZu0AZuSWJ5G\nBJ3eIrIFQbBpHtho2Ueu4JYlifJIEmn5yWNvIHYCumEevXPB5dEASXE7djywteE+\nw3Pi64gYnj3Tb3T8ZIFyWsBqdWZzeFPDUasVi/IuFY/7plDuIOY27BDhhvX2TirH\nOzwEMN9nZ9PWaSRyeHLSslFTjCndoVO90Y95UbBTjz2YO/nueG+4UN08ApqSffzk\n1kpyPWcIU0/8\n-----END CERTIFICATE-----\n",
    "kubernetes_host": "https://dummy.my-prod-cluster.com",
    "pem_keys": [],
    "token_reviewer_jwt_set": true,
    "use_annotations_as_alias_metadata": false
  },
  "warnings": null,
  "mount_type": "kubernetes"
}


$ vault read -non-interactive \
   -format json \
   /auth/production/config | jq
{
  "request_id": "9c87ca62-71f0-63e2-3b67-b37a7800849b",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "disable_iss_validation": true,
    "disable_local_ca_jwt": true,
    "issuer": "",
    "kubernetes_ca_cert": "-----BEGIN CERTIFICATE-----\nMIIDBTCCAe2gAwIBAgIIGkf4cKDZKbgwDQYJKoZIhvcNAQELBQAwFTETMBEGA1UE\nAxMKa3ViZXJuZXRlczAeFw0yNDA0MDMwNzQ5MjhaFw0zNDA0MDEwNzU0MjhaMBUx\nEzARBgNVBAMTCmt1YmVybmV0ZXMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK\nAoIBAQDf8V5Mxa5OAtZckVZo4RslyD1zgtpnkLCJYIPqD23U1cmmF75VV04httOE\nsPHv610WhkOje2jNEZO4SY0wJi6A9QVOJyyCfXAzehY4IYZCWWbfFL99dg28WH7N\ntEpU648GUB9M8Sd/sngof3/CRfi0OELKejmn3xmEYV74Vj3hB57KC8dvNpU0Zgs1\n62oF/ZXXMWLOugM8WonekIwpjy71b3VfRatgBCcqr5yQvyR3r9MJjxgQG/eJAZHI\nUCiIF4GKFsRdCl7hSl+MRf4beg5N2Qc/FeomxxFD8Mc7guYaA5errLlapdHYd0Kx\nVHRwLQ/+hnmnP5FINV+kj7k782c1AgMBAAGjWTBXMA4GA1UdDwEB/wQEAwICpDAP\nBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBSCGsrJ79HjbShkkJQrjmOtl9BX4DAV\nBgNVHREEDjAMggprdWJlcm5ldGVzMA0GCSqGSIb3DQEBCwUAA4IBAQAbC0IJLDVZ\nsYK5sOZi4Z0WwdLKHn5jTF0cS+6E4gWO/3qZMH1lQELlvUa6B3rCOvJ12/a4MoWL\n9JGzRNqi1G9Nox93OW9MrJsfRN+a6HB7cq0qMPhCRv+h7KBurN+MRZu0AZuSWJ5G\nBJ3eIrIFQbBpHtho2Ueu4JYlifJIEmn5yWNvIHYCumEevXPB5dEASXE7djywteE+\nw3Pi64gYnj3Tb3T8ZIFyWsBqdWZzeFPDUasVi/IuFY/7plDuIOY27BDhhvX2TirH\nOzwEMN9nZ9PWaSRyeHLSslFTjCndoVO90Y95UbBTjz2YO/nueG+4UN08ApqSffzk\n1kpyPWcIU0/8\n-----END CERTIFICATE-----\n",
    "kubernetes_host": "https://dummy.my-prod-cluster.com",
    "pem_keys": [],
    "token_reviewer_jwt_set": true,
    "use_annotations_as_alias_metadata": false
  },
  "warnings": null,
  "mount_type": "kubernetes"
}


$ vault write -non-interactive \
   /auth/production/role/dummy \
   alias_name_source=serviceaccount_uid \
   audience=dummy \
   bound_service_account_names=dummy-service-account \
   bound_service_account_namespace_selector='{"matchLabels":{"environment":"prod","org":"company-1"}}' \
   bound_service_account_namespaces=dummy-namespace \
   token_bound_cidrs=192.168.1.0/24 \
   token_bound_cidrs=10.0.1.0/24 \
   token_bound_cidrs=10.0.2.0/24 \
   token_explicit_max_ttl=1200s \
   token_max_ttl=1200s \
   token_no_default_policy=true \
   token_num_uses=10 \
   token_period=600.9s \
   token_policies=dummy-policy \
   token_ttl=10m \
   token_type=default
Success! Data written to: auth/production/role/dummy


$ # yes, 600.9s is a valid value :O :D It will be taken as 600s by Vault though


$ vault list -non-interactive \
   auth/production/role
Keys
----
dummy


$ vault list -non-interactive \
   -format json \
   /auth/production/role
[
  "dummy"
]


$ vault list -non-interactive \
   -format json \
   /auth/production/role | jq
[
  "dummy"
]


$ vault read -non-interactive \
   -format json \
   auth/production/role/dummy
{
  "request_id": "70083054-1935-4d19-fcf0-8193f7f063ae",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "alias_name_source": "serviceaccount_uid",
    "audience": "dummy",
    "bound_service_account_names": [
      "dummy-service-account"
    ],
    "bound_service_account_namespace_selector": "{\"matchLabels\":{\"environment\":\"prod\",\"org\":\"company-1\"}}",
    "bound_service_account_namespaces": [
      "dummy-namespace"
    ],
    "token_bound_cidrs": [
      "192.168.1.0/24",
      "10.0.1.0/24",
      "10.0.2.0/24"
    ],
    "token_explicit_max_ttl": 1200,
    "token_max_ttl": 1200,
    "token_no_default_policy": true,
    "token_num_uses": 10,
    "token_period": 600,
    "token_policies": [
      "dummy-policy"
    ],
    "token_ttl": 600,
    "token_type": "default"
  },
  "warnings": null,
  "mount_type": "kubernetes"
}


$ vault read -non-interactive \
   -format json \
   /auth/production/role/dummy | jq
{
  "request_id": "a41a08b2-2727-7984-a9b2-63ae2052cefc",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "alias_name_source": "serviceaccount_uid",
    "audience": "dummy",
    "bound_service_account_names": [
      "dummy-service-account"
    ],
    "bound_service_account_namespace_selector": "{\"matchLabels\":{\"environment\":\"prod\",\"org\":\"company-1\"}}",
    "bound_service_account_namespaces": [
      "dummy-namespace"
    ],
    "token_bound_cidrs": [
      "192.168.1.0/24",
      "10.0.1.0/24",
      "10.0.2.0/24"
    ],
    "token_explicit_max_ttl": 1200,
    "token_max_ttl": 1200,
    "token_no_default_policy": true,
    "token_num_uses": 10,
    "token_period": 600,
    "token_policies": [
      "dummy-policy"
    ],
    "token_ttl": 600,
    "token_type": "default"
  },
  "warnings": null,
  "mount_type": "kubernetes"
}

$ vault write -non-interactive \
   /auth/production/role/dummy-2 \
   alias_name_source=serviceaccount_uid \
   audience=dummy-2 \
   bound_service_account_names=dummy-2-service-account \
   bound_service_account_namespace_selector='{"matchLabels":{"environment":"prod","org":"company-1"}}' \
   bound_service_account_namespaces=dummy-2-namespace \
   token_bound_cidrs=192.168.1.0/24 \
   token_bound_cidrs=10.0.1.0/24 \
   token_bound_cidrs=10.0.2.0/24 \
   token_explicit_max_ttl=1200s \
   token_max_ttl=1200s \
   token_no_default_policy=true \
   token_num_uses=10 \
   token_period=600.9s \
   token_policies=dummy-2-policy \
   token_ttl=10m \
   token_type=default
Success! Data written to: auth/production/role/dummy-2


$ vault list -non-interactive \
   auth/production/role
Keys
----
dummy
dummy-2


$ vault list -non-interactive \
   -format json \
   /auth/production/role
[
  "dummy",
  "dummy-2"
]


$ vault list -non-interactive \
   -format json \
   /auth/production/role | jq
[
  "dummy",
  "dummy-2"
]


$ vault read -non-interactive \
   -format json \
   auth/production/role/dummy-2
{
  "request_id": "5dee97d0-df86-35d3-3d2c-6f8cbe852a1b",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "alias_name_source": "serviceaccount_uid",
    "audience": "dummy-2",
    "bound_service_account_names": [
      "dummy-2-service-account"
    ],
    "bound_service_account_namespace_selector": "{\"matchLabels\":{\"environment\":\"prod\",\"org\":\"company-1\"}}",
    "bound_service_account_namespaces": [
      "dummy-2-namespace"
    ],
    "token_bound_cidrs": [
      "192.168.1.0/24",
      "10.0.1.0/24",
      "10.0.2.0/24"
    ],
    "token_explicit_max_ttl": 1200,
    "token_max_ttl": 1200,
    "token_no_default_policy": true,
    "token_num_uses": 10,
    "token_period": 600,
    "token_policies": [
      "dummy-2-policy"
    ],
    "token_ttl": 600,
    "token_type": "default"
  },
  "warnings": null,
  "mount_type": "kubernetes"
}


$ vault read -non-interactive \
   -format json \
   /auth/production/role/dummy-2 | jq
{
  "request_id": "c0c1ecc9-0f42-5d8c-36b8-45e519f5a409",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "alias_name_source": "serviceaccount_uid",
    "audience": "dummy-2",
    "bound_service_account_names": [
      "dummy-2-service-account"
    ],
    "bound_service_account_namespace_selector": "{\"matchLabels\":{\"environment\":\"prod\",\"org\":\"company-1\"}}",
    "bound_service_account_namespaces": [
      "dummy-2-namespace"
    ],
    "token_bound_cidrs": [
      "192.168.1.0/24",
      "10.0.1.0/24",
      "10.0.2.0/24"
    ],
    "token_explicit_max_ttl": 1200,
    "token_max_ttl": 1200,
    "token_no_default_policy": true,
    "token_num_uses": 10,
    "token_period": 600,
    "token_policies": [
      "dummy-2-policy"
    ],
    "token_ttl": 600,
    "token_type": "default"
  },
  "warnings": null,
  "mount_type": "kubernetes"
}
```

By the way, again, all this is just for demonstrating (Demo) and teaching purposes only. I'm not an expert in all the nitty gritty details but I know enough to do backups :) :D You just need to know that Kubernetes Auth Methods have Configuration and Roles, and in Roles - for each Role there's a Configuration

So, now, we have all the dummy data we need. We have put the dummy data in and we have also seen the dummy data using `Read` operations

Note that we have different examples of Kubernetes Auth Methods here. One is a Kubernetes Auth Method with No Configuration and No Roles, another is a Kubernetes Auth Method with just Configuration but No Roles, and then another is a Kubernetes Auth Method with Configuration and one Role, and then finally a Kubernetes Auth Method with Configuration and two Roles

Also, you can do the above operations using `curl` too, if you don't have Vault CLI in your local machine. Ideally you would have Vault CLI to run the Vault Server, but you are working with a remote Vault Server, and don't have the Vault CLI to access the Vault Server and don't have any other Vault Clients (GUI clients like Official Vault Web Client etc), the ubiquitous HTTP protocol and HTTP API come to the rescue. Just `curl` is enough. Below is how you can get the `curl` command to any `vault` CLI command, just use `-output-curl-string`

```bash
$ vault auth enable -output-curl-string -path=perf kubernetes
curl -X POST -H "X-Vault-Request: true" -d '{"type":"kubernetes","description":"","config":{"options":null,"default_lease_ttl":"0s","max_lease_ttl":"0s","force_no_cache":false},"local":false,"seal_wrap":false,"external_entropy_access":false,"options":null}' https://127.0.0.1:8200/v1/sys/auth/perf
```

I've put the `curl` commands for all the above commands in [here](curl-commands-for-demo.md) - [curl-commands-for-demo.md](curl-commands-for-demo.md)

Now let's create a token which has the least privilege to read the Vault Kubernetes Auth Methods at `perf/`, `staging/`, `production/` and `dummy-test/` mount paths and also list all the Vault Auth Methods. We will be using this token with the least privileges to do the backup so that it's clear that you DON'T need a `root` user's token to do the backup :) :D

```bash
$ export VAULT_ADDR='http://127.0.0.1:8200'
$ export VAULT_TOKEN="root"

$ cat /Users/karuppiah.n/every-day-log/k8s-auth-method-backup-vault-policy.hcl
# For listing all the auth methods
path "sys/auth" {
  capabilities = ["read"]
}

###

# For Vault Kubernetes (K8s) Auth Method's mount path "perf/"

# For reading config
path "auth/perf/config" {
  capabilities = ["read"]
}

# For listing roles
path "auth/perf/role" {
  capabilities = ["list"]
}

# For reading role config for all roles / any role name
path "auth/perf/role/*" {
  capabilities = ["read"]
}

###

# For Vault Kubernetes (K8s) Auth Method's mount path "staging/"

# For reading config
path "auth/staging/config" {
  capabilities = ["read"]
}

# For listing roles
path "auth/staging/role" {
  capabilities = ["list"]
}

# For reading role config for all roles / any role name
path "auth/staging/role/*" {
  capabilities = ["read"]
}

###

# For Vault Kubernetes (K8s) Auth Method's mount path "production/"

# For reading config
path "auth/production/config" {
  capabilities = ["read"]
}

# For listing roles
path "auth/production/role" {
  capabilities = ["list"]
}

# For reading role config for all roles / any role name
path "auth/production/role/*" {
  capabilities = ["read"]
}

###

# For Vault Kubernetes (K8s) Auth Method's mount path "dummy-test/"

# For reading config
path "auth/dummy-test/config" {
  capabilities = ["read"]
}

# For listing roles
path "auth/dummy-test/role" {
  capabilities = ["list"]
}

# For reading role config for all roles / any role name
path "auth/dummy-test/role/*" {
  capabilities = ["read"]
}


$ vault policy write k8s-auth-method-backup /Users/karuppiah.n/every-day-log/k8s-auth-method-backup-vault-policy.hcl
Success! Uploaded policy: k8s-auth-method-backup


$ vault policy read k8s-auth-method-backup
# For listing all the auth methods
path "sys/auth" {
  capabilities = ["read"]
}

###

# For Vault Kubernetes (K8s) Auth Method's mount path "perf/"

# For reading config
path "auth/perf/config" {
  capabilities = ["read"]
}

# For listing roles
path "auth/perf/role" {
  capabilities = ["list"]
}

# For reading role config for all roles / any role name
path "auth/perf/role/*" {
  capabilities = ["read"]
}

###

# For Vault Kubernetes (K8s) Auth Method's mount path "staging/"

# For reading config
path "auth/staging/config" {
  capabilities = ["read"]
}

# For listing roles
path "auth/staging/role" {
  capabilities = ["list"]
}

# For reading role config for all roles / any role name
path "auth/staging/role/*" {
  capabilities = ["read"]
}

###

# For Vault Kubernetes (K8s) Auth Method's mount path "production/"

# For reading config
path "auth/production/config" {
  capabilities = ["read"]
}

# For listing roles
path "auth/production/role" {
  capabilities = ["list"]
}

# For reading role config for all roles / any role name
path "auth/production/role/*" {
  capabilities = ["read"]
}

###

# For Vault Kubernetes (K8s) Auth Method's mount path "dummy-test/"

# For reading config
path "auth/dummy-test/config" {
  capabilities = ["read"]
}

# For listing roles
path "auth/dummy-test/role" {
  capabilities = ["list"]
}

# For reading role config for all roles / any role name
path "auth/dummy-test/role/*" {
  capabilities = ["read"]
}


$ vault token create -policy k8s-auth-method-backup
Key                  Value
---                  -----
token                hvs.CAESIGq2XmozH9u_bLTFVBPVdgelfLDfIF1yodQin7j23nmmGh4KHGh2cy5RMUY4ekVDYVJ5Y2JoWFZmMDJiNXdwVWI
token_accessor       UW8IP5kuTA9cizDWtHNjdvqt
token_duration       768h
token_renewable      true
token_policies       ["default" "k8s-auth-method-backup"]
identity_policies    []
policies             ["default" "k8s-auth-method-backup"]
~ $ 
```

Note that the above token has two Vault policies attached to it - one is `default` policy and another is our custom policy `k8s-auth-method-backup`. You can choose to modify the `default` policy to ensure how much access you want to give to a token by default. In this case, I'm fine with whatever `default` policy Vault is giving by default

If you don't want the token to have the `default` policy attached to it, you can use `-no-default-policy` flag while creating the token. It will look something like this -

```bash
$ vault token create -no-default-policy -policy k8s-auth-method-backup
Key                  Value
---                  -----
token                hvs.CAESIK0JvNFkyhzOjm8wef0Itl_k7sTZziRitmcJKa9iygo_Gh4KHGh2cy4xZm9waFQ0UHRmaXFlZjMzV0lzeThCb3c
token_accessor       X0uPjbnvNTUJuvfKQNY5qXiJ
token_duration       768h
token_renewable      true
token_policies       ["k8s-auth-method-backup"]
identity_policies    []
policies             ["k8s-auth-method-backup"]
```

Now, let's use the second token we created to create a backup of all the Vault Kubernetes Auth Methods mounted at `perf/`, `staging/`, `production/` and `dummy-test/`

```bash
$ export VAULT_ADDR='http://127.0.0.1:8200'
$ export VAULT_TOKEN="hvs.CAESIK0JvNFkyhzOjm8wef0Itl_k7sTZziRitmcJKa9iygo_Gh4KHGh2cy4xZm9waFQ0UHRmaXFlZjMzV0lzeThCb3c"

$ ./vault-k8s-auth-backup --quiet --file my_vault_k8s_auth_backup.json

backing up the vault k8s auth methods at the following mount paths: [dummy-test/ perf/ production/ staging/]
....

$ cat my_vault_k8s_auth_backup.json
{"k8sAuthMethods":[{"mountPath":"dummy-test/"},{"mountPath":"perf/","config":{"kubernetesHost":"https://dummy.my-cluster.com","kubernetesCaCert":"-----BEGIN CERTIFICATE-----\nMIIDBTCCAe2gAwIBAgIIGkf4cKDZKbgwDQYJKoZIhvcNAQELBQAwFTETMBEGA1UE\nAxMKa3ViZXJuZXRlczAeFw0yNDA0MDMwNzQ5MjhaFw0zNDA0MDEwNzU0MjhaMBUx\nEzARBgNVBAMTCmt1YmVybmV0ZXMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK\nAoIBAQDf8V5Mxa5OAtZckVZo4RslyD1zgtpnkLCJYIPqD23U1cmmF75VV04httOE\nsPHv610WhkOje2jNEZO4SY0wJi6A9QVOJyyCfXAzehY4IYZCWWbfFL99dg28WH7N\ntEpU648GUB9M8Sd/sngof3/CRfi0OELKejmn3xmEYV74Vj3hB57KC8dvNpU0Zgs1\n62oF/ZXXMWLOugM8WonekIwpjy71b3VfRatgBCcqr5yQvyR3r9MJjxgQG/eJAZHI\nUCiIF4GKFsRdCl7hSl+MRf4beg5N2Qc/FeomxxFD8Mc7guYaA5errLlapdHYd0Kx\nVHRwLQ/+hnmnP5FINV+kj7k782c1AgMBAAGjWTBXMA4GA1UdDwEB/wQEAwICpDAP\nBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBSCGsrJ79HjbShkkJQrjmOtl9BX4DAV\nBgNVHREEDjAMggprdWJlcm5ldGVzMA0GCSqGSIb3DQEBCwUAA4IBAQAbC0IJLDVZ\nsYK5sOZi4Z0WwdLKHn5jTF0cS+6E4gWO/3qZMH1lQELlvUa6B3rCOvJ12/a4MoWL\n9JGzRNqi1G9Nox93OW9MrJsfRN+a6HB7cq0qMPhCRv+h7KBurN+MRZu0AZuSWJ5G\nBJ3eIrIFQbBpHtho2Ueu4JYlifJIEmn5yWNvIHYCumEevXPB5dEASXE7djywteE+\nw3Pi64gYnj3Tb3T8ZIFyWsBqdWZzeFPDUasVi/IuFY/7plDuIOY27BDhhvX2TirH\nOzwEMN9nZ9PWaSRyeHLSslFTjCndoVO90Y95UbBTjz2YO/nueG+4UN08ApqSffzk\n1kpyPWcIU0/8\n-----END CERTIFICATE-----\n","disableLocalCaJwt":true,"tokenReviewerJwtSet":true,"useAnnotationsAsAliasMetadata":false},"roles":[{"aliasNameSource":"serviceaccount_uid","audience":"dummy","boundServiceAccountNames":["dummy-service-account"],"boundServiceAccountNamespaceSelector":"{\"matchLabels\":{\"environment\":\"perf\",\"org\":\"company1\"}}","boundServiceAccountNamespaces":["dummy-namespace"],"name":"dummy","tokenBoundCidrs":["192.168.1.0/24","10.0.1.0/24"],"tokenExplicitMaxTtl":7200,"tokenMaxTtl":7200,"tokenNoDefaultPolicy":true,"tokenNumUses":10,"tokenPeriod":3600,"tokenPolicies":["dummy-policy"],"tokenTtl":3600,"tokenType":"default"}]},{"mountPath":"production/","config":{"kubernetesHost":"https://dummy.my-prod-cluster.com","kubernetesCaCert":"-----BEGIN CERTIFICATE-----\nMIIDBTCCAe2gAwIBAgIIGkf4cKDZKbgwDQYJKoZIhvcNAQELBQAwFTETMBEGA1UE\nAxMKa3ViZXJuZXRlczAeFw0yNDA0MDMwNzQ5MjhaFw0zNDA0MDEwNzU0MjhaMBUx\nEzARBgNVBAMTCmt1YmVybmV0ZXMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK\nAoIBAQDf8V5Mxa5OAtZckVZo4RslyD1zgtpnkLCJYIPqD23U1cmmF75VV04httOE\nsPHv610WhkOje2jNEZO4SY0wJi6A9QVOJyyCfXAzehY4IYZCWWbfFL99dg28WH7N\ntEpU648GUB9M8Sd/sngof3/CRfi0OELKejmn3xmEYV74Vj3hB57KC8dvNpU0Zgs1\n62oF/ZXXMWLOugM8WonekIwpjy71b3VfRatgBCcqr5yQvyR3r9MJjxgQG/eJAZHI\nUCiIF4GKFsRdCl7hSl+MRf4beg5N2Qc/FeomxxFD8Mc7guYaA5errLlapdHYd0Kx\nVHRwLQ/+hnmnP5FINV+kj7k782c1AgMBAAGjWTBXMA4GA1UdDwEB/wQEAwICpDAP\nBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBSCGsrJ79HjbShkkJQrjmOtl9BX4DAV\nBgNVHREEDjAMggprdWJlcm5ldGVzMA0GCSqGSIb3DQEBCwUAA4IBAQAbC0IJLDVZ\nsYK5sOZi4Z0WwdLKHn5jTF0cS+6E4gWO/3qZMH1lQELlvUa6B3rCOvJ12/a4MoWL\n9JGzRNqi1G9Nox93OW9MrJsfRN+a6HB7cq0qMPhCRv+h7KBurN+MRZu0AZuSWJ5G\nBJ3eIrIFQbBpHtho2Ueu4JYlifJIEmn5yWNvIHYCumEevXPB5dEASXE7djywteE+\nw3Pi64gYnj3Tb3T8ZIFyWsBqdWZzeFPDUasVi/IuFY/7plDuIOY27BDhhvX2TirH\nOzwEMN9nZ9PWaSRyeHLSslFTjCndoVO90Y95UbBTjz2YO/nueG+4UN08ApqSffzk\n1kpyPWcIU0/8\n-----END CERTIFICATE-----\n","disableLocalCaJwt":true,"tokenReviewerJwtSet":true,"useAnnotationsAsAliasMetadata":false},"roles":[{"aliasNameSource":"serviceaccount_uid","audience":"dummy","boundServiceAccountNames":["dummy-service-account"],"boundServiceAccountNamespaceSelector":"{\"matchLabels\":{\"environment\":\"prod\",\"org\":\"company-1\"}}","boundServiceAccountNamespaces":["dummy-namespace"],"name":"dummy","tokenBoundCidrs":["192.168.1.0/24","10.0.1.0/24","10.0.2.0/24"],"tokenExplicitMaxTtl":1200,"tokenMaxTtl":1200,"tokenNoDefaultPolicy":true,"tokenNumUses":10,"tokenPeriod":600,"tokenPolicies":["dummy-policy"],"tokenTtl":600,"tokenType":"default"},{"aliasNameSource":"serviceaccount_uid","audience":"dummy-2","boundServiceAccountNames":["dummy-2-service-account"],"boundServiceAccountNamespaceSelector":"{\"matchLabels\":{\"environment\":\"prod\",\"org\":\"company-1\"}}","boundServiceAccountNamespaces":["dummy-2-namespace"],"name":"dummy-2","tokenBoundCidrs":["192.168.1.0/24","10.0.1.0/24","10.0.2.0/24"],"tokenExplicitMaxTtl":1200,"tokenMaxTtl":1200,"tokenNoDefaultPolicy":true,"tokenNumUses":10,"tokenPeriod":600,"tokenPolicies":["dummy-2-policy"],"tokenTtl":600,"tokenType":"default"}]},{"mountPath":"staging/","config":{"kubernetesHost":"https://dummy.my-staging-cluster.com","kubernetesCaCert":"-----BEGIN CERTIFICATE-----\nMIIDBTCCAe2gAwIBAgIIGkf4cKDZKbgwDQYJKoZIhvcNAQELBQAwFTETMBEGA1UE\nAxMKa3ViZXJuZXRlczAeFw0yNDA0MDMwNzQ5MjhaFw0zNDA0MDEwNzU0MjhaMBUx\nEzARBgNVBAMTCmt1YmVybmV0ZXMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK\nAoIBAQDf8V5Mxa5OAtZckVZo4RslyD1zgtpnkLCJYIPqD23U1cmmF75VV04httOE\nsPHv610WhkOje2jNEZO4SY0wJi6A9QVOJyyCfXAzehY4IYZCWWbfFL99dg28WH7N\ntEpU648GUB9M8Sd/sngof3/CRfi0OELKejmn3xmEYV74Vj3hB57KC8dvNpU0Zgs1\n62oF/ZXXMWLOugM8WonekIwpjy71b3VfRatgBCcqr5yQvyR3r9MJjxgQG/eJAZHI\nUCiIF4GKFsRdCl7hSl+MRf4beg5N2Qc/FeomxxFD8Mc7guYaA5errLlapdHYd0Kx\nVHRwLQ/+hnmnP5FINV+kj7k782c1AgMBAAGjWTBXMA4GA1UdDwEB/wQEAwICpDAP\nBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBSCGsrJ79HjbShkkJQrjmOtl9BX4DAV\nBgNVHREEDjAMggprdWJlcm5ldGVzMA0GCSqGSIb3DQEBCwUAA4IBAQAbC0IJLDVZ\nsYK5sOZi4Z0WwdLKHn5jTF0cS+6E4gWO/3qZMH1lQELlvUa6B3rCOvJ12/a4MoWL\n9JGzRNqi1G9Nox93OW9MrJsfRN+a6HB7cq0qMPhCRv+h7KBurN+MRZu0AZuSWJ5G\nBJ3eIrIFQbBpHtho2Ueu4JYlifJIEmn5yWNvIHYCumEevXPB5dEASXE7djywteE+\nw3Pi64gYnj3Tb3T8ZIFyWsBqdWZzeFPDUasVi/IuFY/7plDuIOY27BDhhvX2TirH\nOzwEMN9nZ9PWaSRyeHLSslFTjCndoVO90Y95UbBTjz2YO/nueG+4UN08ApqSffzk\n1kpyPWcIU0/8\n-----END CERTIFICATE-----\n","disableLocalCaJwt":true,"tokenReviewerJwtSet":true,"useAnnotationsAsAliasMetadata":false}}]}

$ cat my_vault_k8s_auth_backup.json | jq
{
  "k8sAuthMethods": [
    {
      "mountPath": "dummy-test/"
    },
    {
      "mountPath": "perf/",
      "config": {
        "kubernetesHost": "https://dummy.my-cluster.com",
        "kubernetesCaCert": "-----BEGIN CERTIFICATE-----\nMIIDBTCCAe2gAwIBAgIIGkf4cKDZKbgwDQYJKoZIhvcNAQELBQAwFTETMBEGA1UE\nAxMKa3ViZXJuZXRlczAeFw0yNDA0MDMwNzQ5MjhaFw0zNDA0MDEwNzU0MjhaMBUx\nEzARBgNVBAMTCmt1YmVybmV0ZXMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK\nAoIBAQDf8V5Mxa5OAtZckVZo4RslyD1zgtpnkLCJYIPqD23U1cmmF75VV04httOE\nsPHv610WhkOje2jNEZO4SY0wJi6A9QVOJyyCfXAzehY4IYZCWWbfFL99dg28WH7N\ntEpU648GUB9M8Sd/sngof3/CRfi0OELKejmn3xmEYV74Vj3hB57KC8dvNpU0Zgs1\n62oF/ZXXMWLOugM8WonekIwpjy71b3VfRatgBCcqr5yQvyR3r9MJjxgQG/eJAZHI\nUCiIF4GKFsRdCl7hSl+MRf4beg5N2Qc/FeomxxFD8Mc7guYaA5errLlapdHYd0Kx\nVHRwLQ/+hnmnP5FINV+kj7k782c1AgMBAAGjWTBXMA4GA1UdDwEB/wQEAwICpDAP\nBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBSCGsrJ79HjbShkkJQrjmOtl9BX4DAV\nBgNVHREEDjAMggprdWJlcm5ldGVzMA0GCSqGSIb3DQEBCwUAA4IBAQAbC0IJLDVZ\nsYK5sOZi4Z0WwdLKHn5jTF0cS+6E4gWO/3qZMH1lQELlvUa6B3rCOvJ12/a4MoWL\n9JGzRNqi1G9Nox93OW9MrJsfRN+a6HB7cq0qMPhCRv+h7KBurN+MRZu0AZuSWJ5G\nBJ3eIrIFQbBpHtho2Ueu4JYlifJIEmn5yWNvIHYCumEevXPB5dEASXE7djywteE+\nw3Pi64gYnj3Tb3T8ZIFyWsBqdWZzeFPDUasVi/IuFY/7plDuIOY27BDhhvX2TirH\nOzwEMN9nZ9PWaSRyeHLSslFTjCndoVO90Y95UbBTjz2YO/nueG+4UN08ApqSffzk\n1kpyPWcIU0/8\n-----END CERTIFICATE-----\n",
        "disableLocalCaJwt": true,
        "tokenReviewerJwtSet": true,
        "useAnnotationsAsAliasMetadata": false
      },
      "roles": [
        {
          "aliasNameSource": "serviceaccount_uid",
          "audience": "dummy",
          "boundServiceAccountNames": [
            "dummy-service-account"
          ],
          "boundServiceAccountNamespaceSelector": "{\"matchLabels\":{\"environment\":\"perf\",\"org\":\"company1\"}}",
          "boundServiceAccountNamespaces": [
            "dummy-namespace"
          ],
          "name": "dummy",
          "tokenBoundCidrs": [
            "192.168.1.0/24",
            "10.0.1.0/24"
          ],
          "tokenExplicitMaxTtl": 7200,
          "tokenMaxTtl": 7200,
          "tokenNoDefaultPolicy": true,
          "tokenNumUses": 10,
          "tokenPeriod": 3600,
          "tokenPolicies": [
            "dummy-policy"
          ],
          "tokenTtl": 3600,
          "tokenType": "default"
        }
      ]
    },
    {
      "mountPath": "production/",
      "config": {
        "kubernetesHost": "https://dummy.my-prod-cluster.com",
        "kubernetesCaCert": "-----BEGIN CERTIFICATE-----\nMIIDBTCCAe2gAwIBAgIIGkf4cKDZKbgwDQYJKoZIhvcNAQELBQAwFTETMBEGA1UE\nAxMKa3ViZXJuZXRlczAeFw0yNDA0MDMwNzQ5MjhaFw0zNDA0MDEwNzU0MjhaMBUx\nEzARBgNVBAMTCmt1YmVybmV0ZXMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK\nAoIBAQDf8V5Mxa5OAtZckVZo4RslyD1zgtpnkLCJYIPqD23U1cmmF75VV04httOE\nsPHv610WhkOje2jNEZO4SY0wJi6A9QVOJyyCfXAzehY4IYZCWWbfFL99dg28WH7N\ntEpU648GUB9M8Sd/sngof3/CRfi0OELKejmn3xmEYV74Vj3hB57KC8dvNpU0Zgs1\n62oF/ZXXMWLOugM8WonekIwpjy71b3VfRatgBCcqr5yQvyR3r9MJjxgQG/eJAZHI\nUCiIF4GKFsRdCl7hSl+MRf4beg5N2Qc/FeomxxFD8Mc7guYaA5errLlapdHYd0Kx\nVHRwLQ/+hnmnP5FINV+kj7k782c1AgMBAAGjWTBXMA4GA1UdDwEB/wQEAwICpDAP\nBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBSCGsrJ79HjbShkkJQrjmOtl9BX4DAV\nBgNVHREEDjAMggprdWJlcm5ldGVzMA0GCSqGSIb3DQEBCwUAA4IBAQAbC0IJLDVZ\nsYK5sOZi4Z0WwdLKHn5jTF0cS+6E4gWO/3qZMH1lQELlvUa6B3rCOvJ12/a4MoWL\n9JGzRNqi1G9Nox93OW9MrJsfRN+a6HB7cq0qMPhCRv+h7KBurN+MRZu0AZuSWJ5G\nBJ3eIrIFQbBpHtho2Ueu4JYlifJIEmn5yWNvIHYCumEevXPB5dEASXE7djywteE+\nw3Pi64gYnj3Tb3T8ZIFyWsBqdWZzeFPDUasVi/IuFY/7plDuIOY27BDhhvX2TirH\nOzwEMN9nZ9PWaSRyeHLSslFTjCndoVO90Y95UbBTjz2YO/nueG+4UN08ApqSffzk\n1kpyPWcIU0/8\n-----END CERTIFICATE-----\n",
        "disableLocalCaJwt": true,
        "tokenReviewerJwtSet": true,
        "useAnnotationsAsAliasMetadata": false
      },
      "roles": [
        {
          "aliasNameSource": "serviceaccount_uid",
          "audience": "dummy",
          "boundServiceAccountNames": [
            "dummy-service-account"
          ],
          "boundServiceAccountNamespaceSelector": "{\"matchLabels\":{\"environment\":\"prod\",\"org\":\"company-1\"}}",
          "boundServiceAccountNamespaces": [
            "dummy-namespace"
          ],
          "name": "dummy",
          "tokenBoundCidrs": [
            "192.168.1.0/24",
            "10.0.1.0/24",
            "10.0.2.0/24"
          ],
          "tokenExplicitMaxTtl": 1200,
          "tokenMaxTtl": 1200,
          "tokenNoDefaultPolicy": true,
          "tokenNumUses": 10,
          "tokenPeriod": 600,
          "tokenPolicies": [
            "dummy-policy"
          ],
          "tokenTtl": 600,
          "tokenType": "default"
        },
        {
          "aliasNameSource": "serviceaccount_uid",
          "audience": "dummy-2",
          "boundServiceAccountNames": [
            "dummy-2-service-account"
          ],
          "boundServiceAccountNamespaceSelector": "{\"matchLabels\":{\"environment\":\"prod\",\"org\":\"company-1\"}}",
          "boundServiceAccountNamespaces": [
            "dummy-2-namespace"
          ],
          "name": "dummy-2",
          "tokenBoundCidrs": [
            "192.168.1.0/24",
            "10.0.1.0/24",
            "10.0.2.0/24"
          ],
          "tokenExplicitMaxTtl": 1200,
          "tokenMaxTtl": 1200,
          "tokenNoDefaultPolicy": true,
          "tokenNumUses": 10,
          "tokenPeriod": 600,
          "tokenPolicies": [
            "dummy-2-policy"
          ],
          "tokenTtl": 600,
          "tokenType": "default"
        }
      ]
    },
    {
      "mountPath": "staging/",
      "config": {
        "kubernetesHost": "https://dummy.my-staging-cluster.com",
        "kubernetesCaCert": "-----BEGIN CERTIFICATE-----\nMIIDBTCCAe2gAwIBAgIIGkf4cKDZKbgwDQYJKoZIhvcNAQELBQAwFTETMBEGA1UE\nAxMKa3ViZXJuZXRlczAeFw0yNDA0MDMwNzQ5MjhaFw0zNDA0MDEwNzU0MjhaMBUx\nEzARBgNVBAMTCmt1YmVybmV0ZXMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK\nAoIBAQDf8V5Mxa5OAtZckVZo4RslyD1zgtpnkLCJYIPqD23U1cmmF75VV04httOE\nsPHv610WhkOje2jNEZO4SY0wJi6A9QVOJyyCfXAzehY4IYZCWWbfFL99dg28WH7N\ntEpU648GUB9M8Sd/sngof3/CRfi0OELKejmn3xmEYV74Vj3hB57KC8dvNpU0Zgs1\n62oF/ZXXMWLOugM8WonekIwpjy71b3VfRatgBCcqr5yQvyR3r9MJjxgQG/eJAZHI\nUCiIF4GKFsRdCl7hSl+MRf4beg5N2Qc/FeomxxFD8Mc7guYaA5errLlapdHYd0Kx\nVHRwLQ/+hnmnP5FINV+kj7k782c1AgMBAAGjWTBXMA4GA1UdDwEB/wQEAwICpDAP\nBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBSCGsrJ79HjbShkkJQrjmOtl9BX4DAV\nBgNVHREEDjAMggprdWJlcm5ldGVzMA0GCSqGSIb3DQEBCwUAA4IBAQAbC0IJLDVZ\nsYK5sOZi4Z0WwdLKHn5jTF0cS+6E4gWO/3qZMH1lQELlvUa6B3rCOvJ12/a4MoWL\n9JGzRNqi1G9Nox93OW9MrJsfRN+a6HB7cq0qMPhCRv+h7KBurN+MRZu0AZuSWJ5G\nBJ3eIrIFQbBpHtho2Ueu4JYlifJIEmn5yWNvIHYCumEevXPB5dEASXE7djywteE+\nw3Pi64gYnj3Tb3T8ZIFyWsBqdWZzeFPDUasVi/IuFY/7plDuIOY27BDhhvX2TirH\nOzwEMN9nZ9PWaSRyeHLSslFTjCndoVO90Y95UbBTjz2YO/nueG+4UN08ApqSffzk\n1kpyPWcIU0/8\n-----END CERTIFICATE-----\n",
        "disableLocalCaJwt": true,
        "tokenReviewerJwtSet": true,
        "useAnnotationsAsAliasMetadata": false
      }
    }
  ]
}
```

As you can see, all the Vault Kubernetes Auth Methods in the Vault Instance / Vault Server have been backed up :) 

# Possible Errors

There are quite some possible errors you can face. Mostly relating to one of the following

- DNS Resolution issues. If you are accessing Vault using it's domain name (DNS record), and not an IP address, then ensure that the DNS resolution works well
- Connectivity issues with Vault. Ensure you have good network connectivity to the Vault system. Ensure the IP you are connecting to is right and belongs to the Vault API server, and also check the API server port too.
- Access / Authorization issues. Ensure you have enough access to list and read the Vault Kubernetes Auth Methods that you want to take backup of - which includes - Kubernetes Auth Method Configuration, Kubernetes Auth Method Roles and Roles Configuration

Example access errors / authorization errors / permission errors

- Add Examples Here! [TODO]

# Contributing

Please look at https://github.com/karuppiah7890/vault-tooling-contributions for some basic details on how you can contribute

# Future Ideas

- Any and all issues in the [GitHub Issues](https://github.com/karuppiah7890/vault-k8s-auth-backup/issues) section

- Allow user to say "It's okay if the tool cannot backup some kubernetes auth methods and/ some kubernetes auth method mount paths, due to permission issues. Just backup the kubernetes auth methods the tool can" and be able to skip intermittent errors here and there and ignore the errors, rather than abruptly stop at errors

- Support backing up multiple specific Kubernetes Auth Methods in a single backup at once by providing a file which contains the mount paths of the Kubernetes Auth Methods to be backed up, or by providing the mount paths of the Kubernetes Auth Methods as arguments to the CLI, or provide the ability to use either of the two or even both

- Allow for flags to come even after the arguments. This would require using better CLI libraries / frameworks like [`spf13/pflag`](https://github.com/spf13/pflag), [`spf13/cobra`](https://github.com/spf13/cobra), [`urfave/cli`](https://github.com/urfave/cli) etc than using the basic Golang's built-in `flag` standard library / package

- Backup Vault Policies as part of backing up K8s auth methods

- Consider creating an option to be able to backup or not backup Vault Policies. By default we can maybe backup the Vault Policies

- Consider creating an option to be able to backup Vault Policies in a separate file. Then maybe K8s auth method restore tool can use separate backup JSON files for K8s auth method backup and Vault Policies backup and restore Vault Policies based on some option / flag

- Decide on how the backup JSON file should look like. Currently it's minimal - we remove any data that's not required and/ is empty. This makes the JSON file pretty compact / concise / minimal instead of putting lot of empty values like this -

```json
{
  "k8sAuthMethods": [
    {
      "mountPath": "production/",
      "config": {},
      "roles": []
    }
  ]
}
```

or like these -

```json
{
  "k8sAuthMethods": [
    {
      "mountPath": "production/",
      "config": {
        "kubernetesHost": "",
        "kubernetesCaCert": "",
        "disableLocalCaJwt": false,
        "tokenReviewerJwtSet": false,
        "useAnnotationsAsAliasMetadata": false
      },
      "roles": []
    }
  ]
}
```

Actually it's bad if `kubernetesHost` field / `kubernetes_host` is not present as it's the only mandatory field in the kubernetes auth method `config`. Others are mandatory depending upon other parameters. For example, if `disable_local_ca_jwt` is true, then it's required to set the `token_reviewer_jwt`, `kubernetes_ca_cert` and `kubernetes_host` so that the Vault system can talk to the external Kubernetes API server using the host, CA certificate(s), and token for authentication and authorization to do token reviews of Service Account tokens that Vault receives

For now, if a value is not present in the JSON file, we will assume it's empty or not present when doing restore. As long as there's no problem with not knowing difference between empty and not present, I think we are good. Or else we might have to differentiate between empty value and value not present.

- Consider backing up the auth method mount config. You can find this mount config by using

```bash
vault read sys/auth

# OR

vault read -format json sys/auth
```

You can find the mount config in `config` and other metadata in the top level. It will look something like this -

```bash
$ vault read -format json sys/auth
```

```json
{
  "request_id": "aff84721-bca2-4213-0228-ab72cb02c3cc",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "kubernetes/": {
      "accessor": "auth_kubernetes_38a45db5",
      "config": {
        "default_lease_ttl": 0,
        "force_no_cache": false,
        "max_lease_ttl": 0,
        "token_type": "default-service"
      },
      "deprecation_status": "supported",
      "description": "",
      "external_entropy_access": false,
      "local": false,
      "options": null,
      "plugin_version": "",
      "running_plugin_version": "v0.19.0+builtin",
      "running_sha256": "",
      "seal_wrap": false,
      "type": "kubernetes",
      "uuid": "1d2b4e16-10d2-be19-ca80-2c8a90960589"
    },
    "production-copy/": {
      "accessor": "auth_kubernetes_ecd9b889",
      "config": {
        "default_lease_ttl": 0,
        "force_no_cache": false,
        "max_lease_ttl": 0,
        "token_type": "default-service"
      },
      "deprecation_status": "supported",
      "description": "",
      "external_entropy_access": false,
      "local": false,
      "options": null,
      "plugin_version": "",
      "running_plugin_version": "v0.19.0+builtin",
      "running_sha256": "",
      "seal_wrap": false,
      "type": "kubernetes",
      "uuid": "56872f4c-cab4-1b8c-e681-d7b2c36276dd"
    },
    "production/": {
      "accessor": "auth_kubernetes_22442d97",
      "config": {
        "default_lease_ttl": 0,
        "force_no_cache": false,
        "max_lease_ttl": 0,
        "token_type": "default-service"
      },
      "deprecation_status": "supported",
      "description": "",
      "external_entropy_access": false,
      "local": false,
      "options": null,
      "plugin_version": "",
      "running_plugin_version": "v0.19.0+builtin",
      "running_sha256": "",
      "seal_wrap": false,
      "type": "kubernetes",
      "uuid": "4770e9e9-171a-bd7b-e718-843cefbd56e4"
    },
    "token/": {
      "accessor": "auth_token_67ab9c96",
      "config": {
        "default_lease_ttl": 0,
        "force_no_cache": false,
        "max_lease_ttl": 0,
        "token_type": "default-service"
      },
      "description": "token based credentials",
      "external_entropy_access": false,
      "local": false,
      "options": null,
      "plugin_version": "",
      "running_plugin_version": "v1.17.0+builtin.vault",
      "running_sha256": "",
      "seal_wrap": false,
      "type": "token",
      "uuid": "4c0b0f63-ff85-35e1-7245-d25a20970b8b"
    }
  },
  "warnings": null,
  "mount_type": "system"
}
```

The top level metadata are options you can give while doing a `vault auth enable` command

Note that some metadata might be immutable, like `running_plugin_version`, I think `running_plugin_version` is immutable

```bash
$ vault version
Vault v1.17.0 (72850df1bc10581b74ba5f0f7b3736f31c034235), built 2024-06-10T10:11:34Z

$ vault auth enable --help
```

```bash
$ vault version
Vault v1.17.0 (72850df1bc10581b74ba5f0f7b3736f31c034235), built 2024-06-10T10:11:34Z

$ vault auth enable --help
Usage: vault auth enable [options] TYPE

  Enables a new auth method. An auth method is responsible for authenticating
  users or machines and assigning them policies with which they can access
  Vault.

  Enable the userpass auth method at userpass/:

      $ vault auth enable userpass

  Enable the LDAP auth method at auth-prod/:

      $ vault auth enable -path=auth-prod ldap

  Enable a custom auth plugin (after it's registered in the plugin registry):

      $ vault auth enable -path=my-auth -plugin-name=my-auth-plugin plugin

      OR (preferred way):

      $ vault auth enable -path=my-auth my-auth-plugin

HTTP Options:

  -address=<string>
      Address of the Vault server. The default is https://127.0.0.1:8200. This
      can also be specified via the VAULT_ADDR environment variable.

  -agent-address=<string>
      Address of the Agent. This can also be specified via the
      VAULT_AGENT_ADDR environment variable.

  -ca-cert=<string>
      Path on the local disk to a single PEM-encoded CA certificate to verify
      the Vault server's SSL certificate. This takes precedence over -ca-path.
      This can also be specified via the VAULT_CACERT environment variable.

  -ca-path=<string>
      Path on the local disk to a directory of PEM-encoded CA certificates to
      verify the Vault server's SSL certificate. This can also be specified
      via the VAULT_CAPATH environment variable.

  -client-cert=<string>
      Path on the local disk to a single PEM-encoded CA certificate to use
      for TLS authentication to the Vault server. If this flag is specified,
      -client-key is also required. This can also be specified via the
      VAULT_CLIENT_CERT environment variable.

  -client-key=<string>
      Path on the local disk to a single PEM-encoded private key matching the
      client certificate from -client-cert. This can also be specified via the
      VAULT_CLIENT_KEY environment variable.

  -disable-redirects
      Disable the default client behavior, which honors a single redirect
      response from a request The default is false. This can also be specified
      via the VAULT_DISABLE_REDIRECTS environment variable.

  -header=<key=value>
      Key-value pair provided as key=value to provide http header added to any
      request done by the CLI.Trying to add headers starting with 'X-Vault-'
      is forbidden and will make the command fail This can be specified
      multiple times.

  -mfa=<string>
      Supply MFA credentials as part of X-Vault-MFA header. This can also be
      specified via the VAULT_MFA environment variable.

  -namespace=<string>
      The namespace to use for the command. Setting this is not necessary
      but allows using relative paths. -ns can be used as shortcut. The
      default is (not set). This can also be specified via the VAULT_NAMESPACE
      environment variable.

  -non-interactive
      When set true, prevents asking the user for input via the terminal. The
      default is false.

  -output-curl-string
      Instead of executing the request, print an equivalent cURL command
      string and exit. The default is false.

  -output-policy
      Instead of executing the request, print an example HCL policy that would
      be required to run this command, and exit. The default is false.

  -policy-override
      Override a Sentinel policy that has a soft-mandatory enforcement_level
      specified The default is false.

  -tls-server-name=<string>
      Name to use as the SNI host when connecting to the Vault server via TLS.
      This can also be specified via the VAULT_TLS_SERVER_NAME environment
      variable.

  -tls-skip-verify
      Disable verification of TLS certificates. Using this option is highly
      discouraged as it decreases the security of data transmissions to and
      from the Vault server. The default is false. This can also be specified
      via the VAULT_SKIP_VERIFY environment variable.

  -unlock-key=<string>
      Key to unlock a namespace API lock. The default is (not set).

  -wrap-ttl=<duration>
      Wraps the response in a cubbyhole token with the requested TTL. The
      response is available via the "vault unwrap" command. The TTL is
      specified as a numeric string with suffix like "30s" or "5m". This can
      also be specified via the VAULT_WRAP_TTL environment variable.

Command Options:

  -allowed-response-headers=<string>
      Response header value that plugins will be allowed to set. To specify
      multiple values, specify this flag multiple times.

  -audit-non-hmac-request-keys=<string>
      Key that will not be HMAC'd by audit devices in the request data object.
      To specify multiple values, specify this flag multiple times.

  -audit-non-hmac-response-keys=<string>
      Key that will not be HMAC'd by audit devices in the response data
      object. To specify multiple values, specify this flag multiple times.

  -default-lease-ttl=<duration>
      The default lease TTL for this auth method. If unspecified, this
      defaults to the Vault server's globally configured default lease TTL.

  -description=<string>
      Human-friendly description for the purpose of this auth method.

  -external-entropy-access
      Enable auth method to access Vault's external entropy source. The
      default is false.

  -identity-token-key=<string>
      Select the key used to sign plugin identity tokens. The default is
      default.

  -listing-visibility=<string>
      Determines the visibility of the mount in the UI-specific listing
      endpoint.

  -local
      Mark the auth method as local-only. Local auth methods are not
      replicated nor removed by replication. The default is false.

  -max-lease-ttl=<duration>
      The maximum lease TTL for this auth method. If unspecified, this
      defaults to the Vault server's globally configured maximum lease TTL.

  -options=<key=value>
      Key-value pair provided as key=value for the mount options. This can be
      specified multiple times.

  -passthrough-request-headers=<string>
      Request header value that will be sent to the plugin. To specify
      multiple values, specify this flag multiple times.

  -path=<string>
      Place where the auth method will be accessible. This must be unique
      across all auth methods. This defaults to the "type" of the auth method.
      The auth method will be accessible at "/auth/<path>".

  -plugin-name=<string>
      Name of the auth method plugin. This plugin name must already exist in
      the Vault server's plugin catalog.

  -plugin-version=<string>
      Select the semantic version of the plugin to enable.

  -seal-wrap
      Enable seal wrapping of critical values in the secrets engine. The
      default is false.

  -token-type=<string>
      Sets a forced token type for the mount.

  -version=<int>
      Select the version of the auth method to run. Not supported by all auth
      methods.
```
