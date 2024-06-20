
```bash
export VAULT_ADDR="http://127.0.0.1:8200"
export VAULT_TOKEN="root"

curl -X POST -H "X-Vault-Request: true" -H "X-Vault-Token: $VAULT_TOKEN" -d '{"type":"kubernetes","description":"","config":{"options":null,"default_lease_ttl":"0s","max_lease_ttl":"0s","force_no_cache":false},"local":false,"seal_wrap":false,"external_entropy_access":false,"options":null}' $VAULT_ADDR/v1/sys/auth/perf

curl -X POST -H "X-Vault-Request: true" -H "X-Vault-Token: $VAULT_TOKEN" -d '{"type":"kubernetes","description":"","config":{"options":null,"default_lease_ttl":"0s","max_lease_ttl":"0s","force_no_cache":false},"local":false,"seal_wrap":false,"external_entropy_access":false,"options":null}' $VAULT_ADDR/v1/sys/auth/staging

curl -X POST -H "X-Vault-Request: true" -H "X-Vault-Token: $VAULT_TOKEN" -d '{"type":"kubernetes","description":"","config":{"options":null,"default_lease_ttl":"0s","max_lease_ttl":"0s","force_no_cache":false},"local":false,"seal_wrap":false,"external_entropy_access":false,"options":null}' $VAULT_ADDR/v1/sys/auth/production

curl -X POST -H "X-Vault-Request: true" -H "X-Vault-Token: $VAULT_TOKEN" -d '{"type":"kubernetes","description":"","config":{"options":null,"default_lease_ttl":"0s","max_lease_ttl":"0s","force_no_cache":false},"local":false,"seal_wrap":false,"external_entropy_access":false,"options":null}' $VAULT_ADDR/v1/sys/auth/dummy-test

curl -X PUT -H "X-Vault-Token: $VAULT_TOKEN" -H "X-Vault-Request: true" -d '{"disable_local_ca_jwt":"true","kubernetes_ca_cert":"-----BEGIN CERTIFICATE-----\nMIIDBTCCAe2gAwIBAgIIGkf4cKDZKbgwDQYJKoZIhvcNAQELBQAwFTETMBEGA1UE\nAxMKa3ViZXJuZXRlczAeFw0yNDA0MDMwNzQ5MjhaFw0zNDA0MDEwNzU0MjhaMBUx\nEzARBgNVBAMTCmt1YmVybmV0ZXMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK\nAoIBAQDf8V5Mxa5OAtZckVZo4RslyD1zgtpnkLCJYIPqD23U1cmmF75VV04httOE\nsPHv610WhkOje2jNEZO4SY0wJi6A9QVOJyyCfXAzehY4IYZCWWbfFL99dg28WH7N\ntEpU648GUB9M8Sd/sngof3/CRfi0OELKejmn3xmEYV74Vj3hB57KC8dvNpU0Zgs1\n62oF/ZXXMWLOugM8WonekIwpjy71b3VfRatgBCcqr5yQvyR3r9MJjxgQG/eJAZHI\nUCiIF4GKFsRdCl7hSl+MRf4beg5N2Qc/FeomxxFD8Mc7guYaA5errLlapdHYd0Kx\nVHRwLQ/+hnmnP5FINV+kj7k782c1AgMBAAGjWTBXMA4GA1UdDwEB/wQEAwICpDAP\nBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBSCGsrJ79HjbShkkJQrjmOtl9BX4DAV\nBgNVHREEDjAMggprdWJlcm5ldGVzMA0GCSqGSIb3DQEBCwUAA4IBAQAbC0IJLDVZ\nsYK5sOZi4Z0WwdLKHn5jTF0cS+6E4gWO/3qZMH1lQELlvUa6B3rCOvJ12/a4MoWL\n9JGzRNqi1G9Nox93OW9MrJsfRN+a6HB7cq0qMPhCRv+h7KBurN+MRZu0AZuSWJ5G\nBJ3eIrIFQbBpHtho2Ueu4JYlifJIEmn5yWNvIHYCumEevXPB5dEASXE7djywteE+\nw3Pi64gYnj3Tb3T8ZIFyWsBqdWZzeFPDUasVi/IuFY/7plDuIOY27BDhhvX2TirH\nOzwEMN9nZ9PWaSRyeHLSslFTjCndoVO90Y95UbBTjz2YO/nueG+4UN08ApqSffzk\n1kpyPWcIU0/8\n-----END CERTIFICATE-----\n","kubernetes_host":"https://dummy.my-cluster.com","token_reviewer_jwt":"dummy_jwt_token"}' $VAULT_ADDR/v1/auth/perf/config

curl -H "X-Vault-Request: true" -H "X-Vault-Token: $VAULT_TOKEN" $VAULT_ADDR/v1/auth/perf/config

curl -H "X-Vault-Request: true" -H "X-Vault-Token: $VAULT_TOKEN" $VAULT_ADDR/v1/auth/perf/config | jq

curl -X PUT -H "X-Vault-Request: true" -H "X-Vault-Token: $VAULT_TOKEN" -d '{"alias_name_source":"serviceaccount_uid","audience":"dummy","bound_service_account_names":"dummy-service-account","bound_service_account_namespace_selector":"{\"matchLabels\":{\"environment\":\"perf\",\"org\":\"company1\"}}","bound_service_account_namespaces":"dummy-namespace","token_bound_cidrs":"192.168.1.0/24,10.0.1.0/24","token_explicit_max_ttl":"7200s","token_max_ttl":"7200s","token_no_default_policy":"true","token_num_uses":"10","token_period":"3600.1s","token_policies":"dummy-policy","token_ttl":"1h","token_type":"default"}' $VAULT_ADDR/v1/auth/perf/role/dummy

curl -H "X-Vault-Request: true" -H "X-Vault-Token: $VAULT_TOKEN" $VAULT_ADDR/v1/auth/perf/role?list=true

curl -H "X-Vault-Request: true" -H "X-Vault-Token: $VAULT_TOKEN" $VAULT_ADDR/v1/auth/perf/role?list=true

curl -H "X-Vault-Request: true" -H "X-Vault-Token: $VAULT_TOKEN" $VAULT_ADDR/v1/auth/perf/role?list=true | jq

curl -H "X-Vault-Request: true" -H "X-Vault-Token: $VAULT_TOKEN" $VAULT_ADDR/v1/auth/perf/role/dummy

curl -H "X-Vault-Request: true" -H "X-Vault-Token: $VAULT_TOKEN" $VAULT_ADDR/v1/auth/perf/role/dummy | jq

curl -X PUT -H "X-Vault-Request: true" -H "X-Vault-Token: $VAULT_TOKEN" -d '{"disable_local_ca_jwt":"true","kubernetes_ca_cert":"-----BEGIN CERTIFICATE-----\nMIIDBTCCAe2gAwIBAgIIGkf4cKDZKbgwDQYJKoZIhvcNAQELBQAwFTETMBEGA1UE\nAxMKa3ViZXJuZXRlczAeFw0yNDA0MDMwNzQ5MjhaFw0zNDA0MDEwNzU0MjhaMBUx\nEzARBgNVBAMTCmt1YmVybmV0ZXMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK\nAoIBAQDf8V5Mxa5OAtZckVZo4RslyD1zgtpnkLCJYIPqD23U1cmmF75VV04httOE\nsPHv610WhkOje2jNEZO4SY0wJi6A9QVOJyyCfXAzehY4IYZCWWbfFL99dg28WH7N\ntEpU648GUB9M8Sd/sngof3/CRfi0OELKejmn3xmEYV74Vj3hB57KC8dvNpU0Zgs1\n62oF/ZXXMWLOugM8WonekIwpjy71b3VfRatgBCcqr5yQvyR3r9MJjxgQG/eJAZHI\nUCiIF4GKFsRdCl7hSl+MRf4beg5N2Qc/FeomxxFD8Mc7guYaA5errLlapdHYd0Kx\nVHRwLQ/+hnmnP5FINV+kj7k782c1AgMBAAGjWTBXMA4GA1UdDwEB/wQEAwICpDAP\nBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBSCGsrJ79HjbShkkJQrjmOtl9BX4DAV\nBgNVHREEDjAMggprdWJlcm5ldGVzMA0GCSqGSIb3DQEBCwUAA4IBAQAbC0IJLDVZ\nsYK5sOZi4Z0WwdLKHn5jTF0cS+6E4gWO/3qZMH1lQELlvUa6B3rCOvJ12/a4MoWL\n9JGzRNqi1G9Nox93OW9MrJsfRN+a6HB7cq0qMPhCRv+h7KBurN+MRZu0AZuSWJ5G\nBJ3eIrIFQbBpHtho2Ueu4JYlifJIEmn5yWNvIHYCumEevXPB5dEASXE7djywteE+\nw3Pi64gYnj3Tb3T8ZIFyWsBqdWZzeFPDUasVi/IuFY/7plDuIOY27BDhhvX2TirH\nOzwEMN9nZ9PWaSRyeHLSslFTjCndoVO90Y95UbBTjz2YO/nueG+4UN08ApqSffzk\n1kpyPWcIU0/8\n-----END CERTIFICATE-----\n","kubernetes_host":"https://dummy.my-staging-cluster.com","token_reviewer_jwt":"dummy_jwt_token_2"}' $VAULT_ADDR/v1/auth/staging/config

curl -H "X-Vault-Request: true" -H "X-Vault-Token: $VAULT_TOKEN" $VAULT_ADDR/v1/auth/staging/config

curl -H "X-Vault-Request: true" -H "X-Vault-Token: $VAULT_TOKEN" $VAULT_ADDR/v1/auth/staging/config | jq

curl -X PUT -H "X-Vault-Request: true" -H "X-Vault-Token: $VAULT_TOKEN" -d '{"disable_local_ca_jwt":"true","kubernetes_ca_cert":"-----BEGIN CERTIFICATE-----\nMIIDBTCCAe2gAwIBAgIIGkf4cKDZKbgwDQYJKoZIhvcNAQELBQAwFTETMBEGA1UE\nAxMKa3ViZXJuZXRlczAeFw0yNDA0MDMwNzQ5MjhaFw0zNDA0MDEwNzU0MjhaMBUx\nEzARBgNVBAMTCmt1YmVybmV0ZXMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK\nAoIBAQDf8V5Mxa5OAtZckVZo4RslyD1zgtpnkLCJYIPqD23U1cmmF75VV04httOE\nsPHv610WhkOje2jNEZO4SY0wJi6A9QVOJyyCfXAzehY4IYZCWWbfFL99dg28WH7N\ntEpU648GUB9M8Sd/sngof3/CRfi0OELKejmn3xmEYV74Vj3hB57KC8dvNpU0Zgs1\n62oF/ZXXMWLOugM8WonekIwpjy71b3VfRatgBCcqr5yQvyR3r9MJjxgQG/eJAZHI\nUCiIF4GKFsRdCl7hSl+MRf4beg5N2Qc/FeomxxFD8Mc7guYaA5errLlapdHYd0Kx\nVHRwLQ/+hnmnP5FINV+kj7k782c1AgMBAAGjWTBXMA4GA1UdDwEB/wQEAwICpDAP\nBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBSCGsrJ79HjbShkkJQrjmOtl9BX4DAV\nBgNVHREEDjAMggprdWJlcm5ldGVzMA0GCSqGSIb3DQEBCwUAA4IBAQAbC0IJLDVZ\nsYK5sOZi4Z0WwdLKHn5jTF0cS+6E4gWO/3qZMH1lQELlvUa6B3rCOvJ12/a4MoWL\n9JGzRNqi1G9Nox93OW9MrJsfRN+a6HB7cq0qMPhCRv+h7KBurN+MRZu0AZuSWJ5G\nBJ3eIrIFQbBpHtho2Ueu4JYlifJIEmn5yWNvIHYCumEevXPB5dEASXE7djywteE+\nw3Pi64gYnj3Tb3T8ZIFyWsBqdWZzeFPDUasVi/IuFY/7plDuIOY27BDhhvX2TirH\nOzwEMN9nZ9PWaSRyeHLSslFTjCndoVO90Y95UbBTjz2YO/nueG+4UN08ApqSffzk\n1kpyPWcIU0/8\n-----END CERTIFICATE-----\n","kubernetes_host":"https://dummy.my-prod-cluster.com","token_reviewer_jwt":"dummy_jwt_token_3"}' $VAULT_ADDR/v1/auth/production/config

curl -H "X-Vault-Request: true" -H "X-Vault-Token: $VAULT_TOKEN" $VAULT_ADDR/v1/auth/production/config

curl -H "X-Vault-Request: true" -H "X-Vault-Token: $VAULT_TOKEN" $VAULT_ADDR/v1/auth/production/config | jq

curl -X PUT -H "X-Vault-Request: true" -H "X-Vault-Token: $VAULT_TOKEN" -d '{"alias_name_source":"serviceaccount_uid","audience":"dummy","bound_service_account_names":"dummy-service-account","bound_service_account_namespace_selector":"{\"matchLabels\":{\"environment\":\"prod\",\"org\":\"company-1\"}}","bound_service_account_namespaces":"dummy-namespace","token_bound_cidrs":["192.168.1.0/24","10.0.1.0/24","10.0.2.0/24"],"token_explicit_max_ttl":"1200s","token_max_ttl":"1200s","token_no_default_policy":"true","token_num_uses":"10","token_period":"600.9s","token_policies":"dummy-policy","token_ttl":"10m","token_type":"default"}' $VAULT_ADDR/v1/auth/production/role/dummy

curl -H "X-Vault-Request: true" -H "X-Vault-Token: $VAULT_TOKEN" $VAULT_ADDR/v1/auth/production/role?list=true

curl -H "X-Vault-Request: true" -H "X-Vault-Token: $VAULT_TOKEN" $VAULT_ADDR/v1/auth/production/role?list=true

curl -H "X-Vault-Request: true" -H "X-Vault-Token: $VAULT_TOKEN" $VAULT_ADDR/v1/auth/production/role?list=true | jq

curl -H "X-Vault-Request: true" -H "X-Vault-Token: $VAULT_TOKEN" $VAULT_ADDR/v1/auth/production/role/dummy

curl -H "X-Vault-Request: true" -H "X-Vault-Token: $VAULT_TOKEN" $VAULT_ADDR/v1/auth/production/role/dummy | jq

curl -X PUT -H "X-Vault-Token: $VAULT_TOKEN" -H "X-Vault-Request: true" -d '{"alias_name_source":"serviceaccount_uid","audience":"dummy-2","bound_service_account_names":"dummy-2-service-account","bound_service_account_namespace_selector":"{\"matchLabels\":{\"environment\":\"prod\",\"org\":\"company-1\"}}","bound_service_account_namespaces":"dummy-2-namespace","token_bound_cidrs":["192.168.1.0/24","10.0.1.0/24","10.0.2.0/24"],"token_explicit_max_ttl":"1200s","token_max_ttl":"1200s","token_no_default_policy":"true","token_num_uses":"10","token_period":"600.9s","token_policies":"dummy-2-policy","token_ttl":"10m","token_type":"default"}' $VAULT_ADDR/v1/auth/production/role/dummy-2

curl -H "X-Vault-Request: true" -H "X-Vault-Token: $VAULT_TOKEN" $VAULT_ADDR/v1/auth/production/role?list=true

curl -H "X-Vault-Token: $VAULT_TOKEN" -H "X-Vault-Request: true" $VAULT_ADDR/v1/auth/production/role?list=true

curl -H "X-Vault-Request: true" -H "X-Vault-Token: $VAULT_TOKEN" $VAULT_ADDR/v1/auth/production/role?list=true | jq

curl -H "X-Vault-Request: true" -H "X-Vault-Token: $VAULT_TOKEN" $VAULT_ADDR/v1/auth/production/role/dummy-2

curl -H "X-Vault-Request: true" -H "X-Vault-Token: $VAULT_TOKEN" $VAULT_ADDR/v1/auth/production/role/dummy-2 | jq
```