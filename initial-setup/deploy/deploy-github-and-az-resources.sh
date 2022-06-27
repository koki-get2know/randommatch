#!/bin/bash

# Includes
source _helpers.sh
source _common.sh
source _gh.sh

declare AZURE_LOCATION=""
declare GITHUB_ORG_NAME="expressoa2"
declare TEAM_NAME=""
declare AZURE_DEPLOYMENT=true

declare -r GITHUB_TEMPLATE_OWNER="Microsoft-OpenHack"
declare -r GITHUB_TEMPLATE_REPO="devops-artifacts"
declare -r NAME_PREFIX="randommatch"
declare -r USAGE_HELP="Usage: ./deploy-gh.sh -l <AZURE_LOCATION> [-o <GITHUB_ORG_NAME> -t <TEAM_NAME> -a <AZURE_DEPLOYMENT> true/false]"
declare -r AZURE_SP_JSON="azuresp.json"
declare -r DETAILS_FILE="details-gh.json"


# Verify the type of input and number of values
# Display an error message if the input is not correct
# Exit the shell script with a status of 1 using exit 1 command.
if [ $# -eq 0 ]; then
    _error "${USAGE_HELP}"
    exit 1
fi

# Initialize parameters specified from command line
while getopts ":l:o:t:a:" arg; do
    case "${arg}" in
    l) # Process -l (Location)
        AZURE_LOCATION="${OPTARG}"
        ;;
    o) # Process -o (GitHub Organization)
        GITHUB_ORG_NAME="${OPTARG}"
        ;;
    t) # Process -t (Team Name)
        TEAM_NAME=$(echo "${OPTARG}" | LC_CTYPE=C tr '[:upper:]' '[:lower:]' | LC_CTYPE=C tr -d '[:space:]')
        ;;
    a) # Process -a (Azure Deployment)
        AZURE_DEPLOYMENT="${OPTARG}"
        ;;
    \?)
        _error "Invalid options found: -${OPTARG}."
        _error "${USAGE_HELP}"
        exit 1
        ;;
    esac
done
shift $((OPTIND - 1))

if [ ${#AZURE_LOCATION} -eq 0 ]; then
    _error "Required AZURE_LOCATION parameter is not set!"
    _error "${USAGE_HELP}"
    exit 1
fi

declare -a unsupported_azure_regions=("koreasouth" "westindia" "australiacentral" "australiacentral2" "brazilsoutheast" "francesouth" "germanynorth" "swedencentral" "swedensouth" "uaecentral" "centraluseuap" "eastus2euap" "norwaywest" "westcentralus" "southafricawest" "southindia")
if [[ "${unsupported_azure_regions[*]}" =~ "${AZURE_LOCATION}" ]]; then
    _error "Provided region (${AZURE_LOCATION}) is not supported."
    _error "Unsupported regions:"
    printf '%s\n' "${unsupported_azure_regions[@]}"
    exit 1
fi

if [ ${#GITHUB_ORG_NAME} -eq 0 ]; then
    _error "Required GITHUB_ORG_NAME parameter is not set!"
    _error "${USAGE_HELP}"
    exit 1
fi

# Check for GITHUB_TOKEN
if [ -z ${GITHUB_TOKEN+x} ]; then
    _error "GITHUB_TOKEN does not set!"
    exit 1
fi

# Check for programs
declare -a commands=("az" "jq" "gh" "curl")
check_commands "${commands[@]}"
check_tool_semver "azure-cli" $(az version --output tsv --query \"azure-cli\") "2.34.1"


# CREATE REPOSITORY SECRET
_gh_create_repository_secret() {
    local _repository_secret_name="$1"
    local _repository_full_name="$2"
    local _value="$3"

    gh secret set "${_repository_secret_name}" -b "${_value}" --repo "${_repository_full_name}"
}

gh_create_repository_secrets() {
    local _organization_repository_fullname="$1"

    _azure_parse_json

    _gh_create_repository_secret "ACTIONS_RUNNER_DEBUG" "${_organization_repository_fullname}" "false"
    # _gh_create_repository_secret "RESOURCES_PREFIX" "${_organization_repository_fullname}" "${UNIQUE_NAME}"
    _gh_create_repository_secret "LOCATION" "${_organization_repository_fullname}" "${AZURE_LOCATION}"
    _gh_create_repository_secret "TFSTATE_RESOURCES_GROUP_NAME" "${_organization_repository_fullname}" "${UNIQUE_NAME}staterg"
    _gh_create_repository_secret "TFSTATE_STORAGE_ACCOUNT_NAME" "${_organization_repository_fullname}" "${UNIQUE_NAME}statest"
    _gh_create_repository_secret "TFSTATE_STORAGE_CONTAINER_NAME" "${_organization_repository_fullname}" "tfstate"
    _gh_create_repository_secret "TFSTATE_KEY" "${_organization_repository_fullname}" "terraform.tfstate"
    _gh_create_repository_secret "AZURE_CREDENTIALS" "${_organization_repository_fullname}" "$(cat azuresp.json)"
    _gh_create_repository_secret "ARM_CLIENT_ID" "${_organization_repository_fullname}" "${ARM_CLIENT_ID}"
    _gh_create_repository_secret "ARM_CLIENT_SECRET" "${_organization_repository_fullname}" "${ARM_CLIENT_SECRET}"
    _gh_create_repository_secret "ARM_SUBSCRIPTION_ID" "${_organization_repository_fullname}" "${ARM_SUBSCRIPTION_ID}"
    _gh_create_repository_secret "ARM_TENANT_ID" "${_organization_repository_fullname}" "${ARM_TENANT_ID}"
}

gh_logout() {
    export GITHUB_TOKEN=0
}

save_details() {
    local _board_url="$1"
    local _team_url="$2"
    local _repo_url="$3"

    jq -n \
        --arg teamName "${UNIQUE_NAME}" \
        --arg orgName "${GITHUB_ORG_NAME}" \
        --arg boardUrl "${_board_url}" \
        --arg teamUrl "${_team_url}" \
        --arg repoUrl "${_repo_url}" \
        --arg azRgTfState "https://portal.azure.com/#resource/subscriptions/${ARM_SUBSCRIPTION_ID}/resourceGroups/${UNIQUE_NAME}staterg/overview" \
        --arg TFSTATE_RESOURCES_GROUP_NAME "${UNIQUE_NAME}staterg" \
        --arg TFSTATE_STORAGE_ACCOUNT_NAME "${UNIQUE_NAME}statest" \
        --arg TFSTATE_STORAGE_CONTAINER_NAME "tfstate" \
        --arg TFSTATE_KEY "terraform.tfstate" \
        '{teamName: $teamName, orgName: $orgName, boardUrl: $boardUrl, teamUrl: $teamUrl, repoUrl: $repoUrl, azRgTfState: $azRgTfState, TFSTATE_RESOURCES_GROUP_NAME: $TFSTATE_RESOURCES_GROUP_NAME, TFSTATE_STORAGE_ACCOUNT_NAME: $TFSTATE_STORAGE_ACCOUNT_NAME, TFSTATE_STORAGE_CONTAINER_NAME: $TFSTATE_STORAGE_CONTAINER_NAME, TFSTATE_KEY: $TFSTATE_KEY}' >"${DETAILS_FILE}"
}

# EXECUTE
_information "Getting unique name..."
get_unique_name

_information "Checking for ${AZURE_SP_JSON} file..."
check_azuresp_json

_information "Creating Azure resources..."
create_azure_resources

_organization_repository_fullname="https://github.com/expressoa2/randommatch"

_information "Creating repository secrets..."
gh_create_repository_secrets "${_organization_repository_fullname}"

_information "GitHub logout..."
gh_logout

_team_url="https://github.com/orgs/expressoa2/teams/coremembers"
_repository_url="https://github.com/expressoa2/randommatch"
_project_url="https://github.com/expressoa2"

_information "Saving details to ${DETAILS_FILE} file..."
save_details "${_project_url}" "${_team_url}" "${_repository_url}"

# OUTPUT
echo -e "\n"
_information "Team Name: ${UNIQUE_NAME}"
_information "Project URL: ${_project_url}"
_information "Team URL: ${_team_url}"
_information "Repo URL: ${_repository_url}"
_information "Azure RG for TF State: https://portal.azure.com/#resource/subscriptions/${ARM_SUBSCRIPTION_ID}/resourceGroups/${UNIQUE_NAME}staterg/overview"
echo -e "\n"
_success "Done!"
