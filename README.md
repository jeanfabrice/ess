[<img src="https://github.com/jeanfabrice/ess/actions/workflows/build-and-release.yml/badge.svg?style=flat">](https://github.com/jeanfabrice/ess/actions/workflows/build-and-release.yml) [<img src="https://img.shields.io/github/v/release/jeanfabrice/ess?style=flat">](https://github.com/jeanfabrice/ess/releases/latest)

# ESS

This is a multi purpose command-line utility written in Go to interact with an [Elasticsearch Service](https://www.elastic.co/guide/en/cloud/current/index.html) deployment.

It can:

- fetch details of [traffic filter](https://www.elastic.co/guide/en/cloud/current/ec-traffic-filtering-deployment-configuration.html) rulesets associated with the deployment
- collect [Elasticsearch diagnostics](https://www.elastic.co/guide/en/cloud-enterprise/current/capture-deployment-resource-diagnostics.html) for the deployment
- send simple GET API calls to the Elasticsearch cluster and display the result

## Installation

- Download a [compiled release binary](https://github.com/jeanfabrice/ess/releases/latest) according to your architecture

- Copy it somewhere in your PATH and rename it to `ess` (or `ess.exe` on Windows)

- Define an environment variable `ELASTIC_ESS_KEY` with the value of your [Elastic Cloud API key](https://www.elastic.co/guide/en/cloud/current/ec-api-keys.html)

## Usage

Once installed, you can use the utility as follows:

```bash
ess [-v] [-d|-t] <deployment_id> [Elasticsearch GET command]
```

- `-v`: Verbose mode (optional).
- `-d`: Diagnostics mode, collect an Elasticsearch support diagnostics for the deployment.
- `-t`: Traffic Filters mode, collect Traffic filters ruleset associated with the deployment.
- `<deployment_id>`: ID of the Elasticsearch Service deployment.
- [Elasticsearch GET command]: Command mode, run an Elasticsearch API GET query, like `_cat/indices`

Command mode is mutually exclusive with Diagnostics or Traffic Filters mode.

### Examples

```bash
ess -t abcdef1234567890abcdef1234567890
```

This command fetches and display ruleset details for the deployment with ID `abcdef1234567890abcdef1234567890`.

```bash
ess abcdef1234567890abcdef1234567890 _cat/indices
ess abcdef1234567890abcdef1234567890 /_cat/indices
```

This command run the `_cat/indices` API command on the Elasticsearch cluster of the deployment and display the results

```bash
ess abcdef1234567890abcdef1234567890 | jq '.version.number'
ess abcdef1234567890abcdef1234567890 / | jq '.version.number'
```

This command extracts the version number from the Elasticsearch API endpoint `/`

## License

This utility is released under the MIT License. See [LICENSE](LICENSE) for details.
