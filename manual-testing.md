# Manual Testing

The Manual Testing I do that needs to be automated -

1. Run a Vault Developer Server. Preferably the latest version, but to test - whatever version needs to be tested - that the tool supports and is compatible with

```bash
vault server -dev -dev-root-token-id root -dev-listen-address 127.0.0.1:8200
```

2. Enable some Kubernetes Auth Methods at different Mount Paths

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

3. Create some Kubernetes Auth Method Config and Roles at these different Mount Paths

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

4. Run `vault-k8s-auth-backup` tool

```bash
export VAULT_ADDR="http://127.0.0.1:8200"
export VAULT_TOKEN="root"

./vault-k8s-auth-backup
```

5. Check the `vault_k8s_auth_backup.json` file and verify if the backup has been done correctly and successfully. Correctness is key!


```bash
cat vault_k8s_auth_backup.json

cat vault_k8s_auth_backup.json | jq
```
