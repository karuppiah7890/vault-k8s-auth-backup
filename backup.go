package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/hashicorp/vault/api"
)

type VaultK8sAuthBackup struct {
	K8sAuthMethods VaultK8sAuthMethods `json:"k8sAuthMethods"`
}

type VaultK8sAuthMethods []VaultK8sAuthMethod

type VaultK8sAuthMethod struct {
	MountPath string                    `json:"mountPath"`
	Config    *VaultK8sAuthMethodConfig `json:"config,omitempty"`
	Roles     VaultK8sAuthMethodRoles   `json:"roles,omitempty"`
}

type VaultK8sAuthMethodRoles []VaultK8sAuthMethodRole

type VaultK8sAuthMethodRole struct {
	// This stores the `alias_name_source` field / parameter
	AliasNameSource *string `json:"aliasNameSource,omitempty"`

	// This stores the `audience` field / parameter
	Audience *string `json:"audience,omitempty"`

	// This stores the `bound_service_account_names` field / parameter
	BoundServiceAccountNames []string `json:"boundServiceAccountNames,omitempty"`

	// This stores the `bound_service_account_namespace_selector` field / parameter
	BoundServiceAccountNamespaceSelector *string `json:"boundServiceAccountNamespaceSelector,omitempty"`

	// This stores the `bound_service_account_namespaces` field / parameter
	BoundServiceAccountNamespaces []string `json:"boundServiceAccountNamespaces,omitempty"`

	// This stores the `name` field / parameter
	Name *string `json:"name,omitempty"`

	// This stores the `token_bound_cidrs` field / parameter
	TokenBoundCidrs []string `json:"tokenBoundCidrs,omitempty"`

	// This stores the `token_explicit_max_ttl` field / parameter
	TokenExplicitMaxTtl *json.Number `json:"tokenExplicitMaxTtl,omitempty"`

	// This stores the `token_max_ttl` field / parameter
	TokenMaxTtl *json.Number `json:"tokenMaxTtl,omitempty"`

	// This stores the `token_no_default_policy` field / parameter
	TokenNoDefaultPolicy *bool `json:"tokenNoDefaultPolicy,omitempty"`

	// This stores the `token_num_uses` field / parameter
	TokenNumUses *int64 `json:"tokenNumUses,omitempty"`

	// This stores the `token_period` field / parameter
	TokenPeriod *json.Number `json:"tokenPeriod,omitempty"`

	// This stores the `token_policies` field / parameter
	TokenPolicies []string `json:"tokenPolicies,omitempty"`

	// This stores the `token_ttl` field / parameter
	TokenTtl *json.Number `json:"tokenTtl,omitempty"`

	// This stores the `token_type` field / parameter
	TokenType *string `json:"tokenType,omitempty"`
}

// Note: Not considering deprecated parameters in the config, like:
// - disable_iss_validation
// - issuer
type VaultK8sAuthMethodConfig struct {
	// This stores the `kubernetes_host` field / parameter
	KubernetesHost *string `json:"kubernetesHost,omitempty"`

	// This stores the `kubernetes_ca_cert` field / parameter
	KubernetesCaCert *string `json:"kubernetesCaCert,omitempty"`

	// This stores the `disable_local_ca_jwt` field / parameter
	DisableLocalCaJwt *bool `json:"disableLocalCaJwt,omitempty"`

	// TODO: We CANNOT backup the Token Reviewer JWT token as the Vault Server
	// never provides it as data over the HTTP API once it's set. We can only
	// check if it's set or not, that's all. We cannot see / get the token
	// itself.
	// TokenReviewerJwt              *string  `json:"tokenReviewerJwt,omitempty"`

	// This stores the `token_reviewer_jwt_set` field / parameter
	TokenReviewerJwtSet *bool `json:"tokenReviewerJwtSet,omitempty"`

	// This stores the `pem_keys` field / parameter
	PemKeys []string `json:"pemKeys,omitempty"`

	// This stores the `use_annotations_as_alias_metadata` field / parameter
	UseAnnotationsAsAliasMetadata *bool `json:"useAnnotationsAsAliasMetadata,omitempty"`
}

func convertVaultK8sAuthBackupToJSON(vaultK8sAuthBackup VaultK8sAuthBackup) ([]byte, error) {
	vaultK8sAuthBackupJSON, err := toJSON(vaultK8sAuthBackup)
	if err != nil {
		return nil, err
	}
	return vaultK8sAuthBackupJSON, nil
}

func backupK8sAuthMethods(client *api.Client, k8sAuthMethodMountPaths []string, quietProgress bool) VaultK8sAuthBackup {
	vaultK8sAuthMethods := make(VaultK8sAuthMethods, 0)

	fmt.Fprintf(os.Stdout, "\nbacking up the vault k8s auth methods at the following mount paths: %+v\n", k8sAuthMethodMountPaths)

	// Backup all the given Vault K8s Auth Methods
	for _, k8sAuthMethodMountPath := range k8sAuthMethodMountPaths {
		vaultK8sAuthMethod := getVaultK8sAuthMethod(client, k8sAuthMethodMountPath)

		if quietProgress {
			fmt.Fprintf(os.Stdout, ".")
		} else {
			fmt.Fprintf(os.Stdout, "\nbacking up the vault k8s auth method (mount path = `%s`)\n", k8sAuthMethodMountPath)
			fmt.Fprintf(os.Stdout, "\ndata of vault k8s auth method (mount path = `%s`): %+v\n", k8sAuthMethodMountPath, vaultK8sAuthMethod)
		}

		vaultK8sAuthMethods = append(vaultK8sAuthMethods, vaultK8sAuthMethod)
	}

	fmt.Fprintf(os.Stdout, "\n")

	return VaultK8sAuthBackup{K8sAuthMethods: vaultK8sAuthMethods}
}

func getVaultK8sAuthMethod(client *api.Client, k8sAuthMethodMountPath string) VaultK8sAuthMethod {
	config, err := getVaultK8sAuthMethodConfig(client, k8sAuthMethodMountPath)
	// TODO: Think on this.
	// Note: We are abruptly stopping here, if we cannot read the config,
	// instead of letting the caller handle the error.
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading vault k8s auth method (mount path = `%s`) config from vault: %s\n", k8sAuthMethodMountPath, err)
		os.Exit(1)
	}

	roles, err := getVaultK8sAuthMethodRoles(client, k8sAuthMethodMountPath)
	// TODO: Think on this.
	// Note: We are abruptly stopping here, if we cannot read the roles,
	// instead of letting the caller handle the error.
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading vault k8s auth method (mount path = `%s`) roles from vault: %s\n", k8sAuthMethodMountPath, err)
		os.Exit(1)
	}

	return VaultK8sAuthMethod{
		MountPath: k8sAuthMethodMountPath,
		Config:    config,
		Roles:     roles,
	}
}

const KUBERNETES_HOST_CONFIG_KEY = "kubernetes_host"
const KUBERNETES_CA_CERT_CONFIG_KEY = "kubernetes_ca_cert"
const DISABLE_LOCAL_CA_JWT_CONFIG_KEY = "disable_local_ca_jwt"
const TOKEN_REVIEWER_JWT_SET_CONFIG_KEY = "token_reviewer_jwt_set"
const PEM_KEYS_CONFIG_KEY = "pem_keys"
const USE_ANNOTATIONS_AS_ALIAS_METADATA_CONFIG_KEY = "use_annotations_as_alias_metadata"

func getVaultK8sAuthMethodConfig(client *api.Client, k8sAuthMethodMountPath string) (*VaultK8sAuthMethodConfig, error) {
	k8sAuthMethodConfigPath := path.Join("auth", k8sAuthMethodMountPath, "config")
	rawConfig, err := client.Logical().Read(k8sAuthMethodConfigPath)

	if err != nil {
		return nil, fmt.Errorf("error reading k8s auth method (mount path = `%s`) config at path `%s`: %v", k8sAuthMethodMountPath, k8sAuthMethodConfigPath, err)
	}

	if rawConfig == nil {
		return nil, nil
	}

	config := VaultK8sAuthMethodConfig{}

	if rawKubernetesHost, ok := rawConfig.Data[KUBERNETES_HOST_CONFIG_KEY]; ok {
		if kubernetesHost, stringOk := rawKubernetesHost.(string); stringOk {
			config.KubernetesHost = &kubernetesHost
		}
	}

	if rawKubernetesCaCert, ok := rawConfig.Data[KUBERNETES_CA_CERT_CONFIG_KEY]; ok {
		if kubernetesCaCert, stringOk := rawKubernetesCaCert.(string); stringOk {
			config.KubernetesCaCert = &kubernetesCaCert
		}
	}

	if rawDisableLocalCaJwt, ok := rawConfig.Data[DISABLE_LOCAL_CA_JWT_CONFIG_KEY]; ok {
		if disableLocalCaJwt, booleanOk := rawDisableLocalCaJwt.(bool); booleanOk {
			config.DisableLocalCaJwt = &disableLocalCaJwt
		}
	}

	if rawTokenReviewerJwtSet, ok := rawConfig.Data[TOKEN_REVIEWER_JWT_SET_CONFIG_KEY]; ok {
		if tokenReviewerJwtSet, booleanOk := rawTokenReviewerJwtSet.(bool); booleanOk {
			config.TokenReviewerJwtSet = &tokenReviewerJwtSet
		}
	}

	if rawPemKeys, ok := rawConfig.Data[PEM_KEYS_CONFIG_KEY]; ok {
		if pemKeys, arrayOfStringOk := convertInterfaceToStringArray(rawPemKeys); arrayOfStringOk {
			config.PemKeys = pemKeys
		}
	}

	if rawUseAnnotationsAsAliasMetadata, ok := rawConfig.Data[USE_ANNOTATIONS_AS_ALIAS_METADATA_CONFIG_KEY]; ok {
		if useAnnotationsAsAliasMetadata, booleanOk := rawUseAnnotationsAsAliasMetadata.(bool); booleanOk {
			config.UseAnnotationsAsAliasMetadata = &useAnnotationsAsAliasMetadata
		}
	}

	return &config, nil
}

func getVaultK8sAuthMethodRoles(client *api.Client, k8sAuthMethodMountPath string) (VaultK8sAuthMethodRoles, error) {
	roleNames, err := getVaultK8sAuthMethodRoleNames(client, k8sAuthMethodMountPath)
	if err != nil {
		return nil, fmt.Errorf("error getting k8s auth method (mount path = `%s`) role names: %v", k8sAuthMethodMountPath, err)
	}

	if roleNames == nil {
		return nil, nil
	}

	var roles VaultK8sAuthMethodRoles

	for _, roleName := range roleNames {
		role, err := getVaultK8sAuthMethodRole(client, k8sAuthMethodMountPath, roleName)

		if err != nil {
			return nil, fmt.Errorf("error getting k8s auth method (mount path = `%s`) role info for `%s` role: %v", k8sAuthMethodMountPath, roleName, err)
		}

		if role != nil {
			roles = append(roles, *role)
		}

	}

	return roles, nil
}

func getVaultK8sAuthMethodRoleNames(client *api.Client, k8sAuthMethodMountPath string) ([]string, error) {
	k8sAuthMethodListRolesPath := path.Join("auth", k8sAuthMethodMountPath, "role")
	rawRolesList, err := client.Logical().List(k8sAuthMethodListRolesPath)

	if err != nil {
		return nil, fmt.Errorf("error reading k8s auth method (mount path = `%s`) role names at path `%s`: %v", k8sAuthMethodMountPath, k8sAuthMethodListRolesPath, err)
	}

	if rawRolesList == nil {
		return nil, nil
	}

	var allRoleNames []string

	if rawRoleNames, ok := rawRolesList.Data["keys"]; ok {
		if roleNames, arrayOfStringOk := convertInterfaceToStringArray(rawRoleNames); arrayOfStringOk {
			allRoleNames = roleNames
		}
	}

	return allRoleNames, nil
}

const ALIAS_NAME_SOURCE_ROLE_KEY = "alias_name_source"
const AUDIENCE_ROLE_KEY = "audience"
const BOUND_SERVICE_ACCOUNT_NAMES_ROLE_KEY = "bound_service_account_names"
const BOUND_SERVICE_ACCOUNT_NAMESPACE_SELECTOR_ROLE_KEY = "bound_service_account_namespace_selector"
const BOUND_SERVICE_ACCOUNT_NAMESPACES_ROLE_KEY = "bound_service_account_namespaces"
const TOKEN_BOUND_CIDRS_ROLE_KEY = "token_bound_cidrs"
const TOKEN_EXPLICIT_MAX_TTL_ROLE_KEY = "token_explicit_max_ttl"
const TOKEN_MAX_TTL_ROLE_KEY = "token_max_ttl"
const TOKEN_NO_DEFAULT_POLICY_ROLE_KEY = "token_no_default_policy"
const TOKEN_NUM_USES_ROLE_KEY = "token_num_uses"
const TOKEN_PERIOD_ROLE_KEY = "token_period"
const TOKEN_POLICIES_ROLE_KEY = "token_policies"
const TOKEN_TTL_ROLE_KEY = "token_ttl"
const TOKEN_TYPE_ROLE_KEY = "token_type"

func getVaultK8sAuthMethodRole(client *api.Client, k8sAuthMethodMountPath string, roleName string) (*VaultK8sAuthMethodRole, error) {
	k8sAuthMethodReadRolePath := path.Join("auth", k8sAuthMethodMountPath, "role", roleName)
	rawRole, err := client.Logical().Read(k8sAuthMethodReadRolePath)

	if err != nil {
		return nil, fmt.Errorf("error reading k8s auth method (mount path = `%s`) role `%s` at path `%s`: %v", k8sAuthMethodMountPath, roleName, k8sAuthMethodReadRolePath, err)
	}

	role := VaultK8sAuthMethodRole{
		Name: &roleName,
	}

	// Ideally something should be present in the role as one cannot create a role
	// with just a name and without any other field. But we still catch this case
	// here, just in case, because you never know when things break :) or when
	// something is or is NOT nil / null. We don't want nil pointer exceptions
	// or errors or panics or segmentation fault or segmentation violation
	if rawRole == nil {
		return &role, nil
	}

	if rawAliasNameSource, ok := rawRole.Data[ALIAS_NAME_SOURCE_ROLE_KEY]; ok {
		if aliasNameSource, stringOk := rawAliasNameSource.(string); stringOk {
			role.AliasNameSource = &aliasNameSource
		}
	}

	if rawAudience, ok := rawRole.Data[AUDIENCE_ROLE_KEY]; ok {
		if audience, stringOk := rawAudience.(string); stringOk {
			role.Audience = &audience
		}
	}

	if rawBoundServiceAccountNames, ok := rawRole.Data[BOUND_SERVICE_ACCOUNT_NAMES_ROLE_KEY]; ok {
		if boundServiceAccountNames, arrayOfStringOk := convertInterfaceToStringArray(rawBoundServiceAccountNames); arrayOfStringOk {
			role.BoundServiceAccountNames = boundServiceAccountNames
		}
	}

	if rawBoundServiceAccountNamespaceSelector, ok := rawRole.Data[BOUND_SERVICE_ACCOUNT_NAMESPACE_SELECTOR_ROLE_KEY]; ok {
		if boundServiceAccountNamespaceSelector, stringOk := rawBoundServiceAccountNamespaceSelector.(string); stringOk {
			role.BoundServiceAccountNamespaceSelector = &boundServiceAccountNamespaceSelector
		}
	}

	if rawBoundServiceAccountNamespaces, ok := rawRole.Data[BOUND_SERVICE_ACCOUNT_NAMESPACES_ROLE_KEY]; ok {
		if boundServiceAccountNamespaces, arrayOfStringOk := convertInterfaceToStringArray(rawBoundServiceAccountNamespaces); arrayOfStringOk {
			role.BoundServiceAccountNamespaces = boundServiceAccountNamespaces
		}
	}

	if rawTokenBoundCidrs, ok := rawRole.Data[TOKEN_BOUND_CIDRS_ROLE_KEY]; ok {
		if tokenBoundCidrs, arrayOfStringOk := convertInterfaceToStringArray(rawTokenBoundCidrs); arrayOfStringOk {
			role.TokenBoundCidrs = tokenBoundCidrs
		}
	}

	if rawTokenExplicitMaxTtl, ok := rawRole.Data[TOKEN_EXPLICIT_MAX_TTL_ROLE_KEY]; ok {
		if tokenExplicitMaxTtl, jsonNumberOk := rawTokenExplicitMaxTtl.(json.Number); jsonNumberOk {
			role.TokenExplicitMaxTtl = &tokenExplicitMaxTtl
		}
	}

	if rawTokenMaxTtl, ok := rawRole.Data[TOKEN_MAX_TTL_ROLE_KEY]; ok {
		if tokenMaxTtl, jsonNumberOk := rawTokenMaxTtl.(json.Number); jsonNumberOk {
			role.TokenMaxTtl = &tokenMaxTtl
		}
	}

	if rawTokenNoDefaultPolicy, ok := rawRole.Data[TOKEN_NO_DEFAULT_POLICY_ROLE_KEY]; ok {
		if tokenNoDefaultPolicy, booleanOk := rawTokenNoDefaultPolicy.(bool); booleanOk {
			role.TokenNoDefaultPolicy = &tokenNoDefaultPolicy
		}
	}

	if rawTokenNumUses, ok := rawRole.Data[TOKEN_NUM_USES_ROLE_KEY]; ok {
		if tokenNumUses, jsonNumberOk := rawTokenNumUses.(json.Number); jsonNumberOk {
			// TODO: Ignoring the error here as server will always return an integer only.
			// But if something goes wrong, gotta see how to handle the error here
			value, _ := tokenNumUses.Int64()
			role.TokenNumUses = &value
		}
	}

	if rawTokenPeriod, ok := rawRole.Data[TOKEN_PERIOD_ROLE_KEY]; ok {
		if tokenPeriod, stringOk := rawTokenPeriod.(json.Number); stringOk {
			role.TokenPeriod = &tokenPeriod
		}
	}

	if rawTokenPolicies, ok := rawRole.Data[TOKEN_POLICIES_ROLE_KEY]; ok {
		if tokenPolicies, arrayOfStringOk := convertInterfaceToStringArray(rawTokenPolicies); arrayOfStringOk {
			role.TokenPolicies = tokenPolicies
		}
	}

	if rawTokenTtl, ok := rawRole.Data[TOKEN_TTL_ROLE_KEY]; ok {
		if tokenTtl, jsonNumberOk := rawTokenTtl.(json.Number); jsonNumberOk {
			role.TokenTtl = &tokenTtl
		}
	}

	if rawTokenType, ok := rawRole.Data[TOKEN_TYPE_ROLE_KEY]; ok {
		if tokenType, stringOk := rawTokenType.(string); stringOk {
			role.TokenType = &tokenType
		}
	}

	return &role, nil
}
