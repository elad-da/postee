package main

import (
	"encoding/json"
	"log"

	"github.com/aquasecurity/postee/v2/data"
	"github.com/aquasecurity/postee/v2/router"
)

var body = []byte(`{
	"description": "This is a test message!",
    "image": "test/test-message:1.0",
    "registry": "Test Hub",
    "digest": "sha256:a563877b915679675b2b8ff3cd8b13863b40e9a43c218bed75e6aa14fcae23dc",
    "metadata": {
        "id": "sha256:a563877b915679675b2b8ff3cd8b13863b40e9a43c218bed75e6aa14fcae23dc",
        "repo_digests": [
            "test/sensitivedata-test@sha256:0ace851f52a9579f5087cf39fc15c3b2c4b433bf6be28e6b1def05287ae15c5c"
        ],
        "created": "2021-09-09T07:49:21.291753899Z"
    },
    "image_assurance_results": {
        "disallowed": true,
        "audit_required": true,
        "policy_failures": [
            {
                "policy_id": 7,
                "policy_name": "Sensitive-Data-Default-Policy",
                "blocking": true,
                "controls": [
                    "sensitive_data"
                ]
            }
        ],
        "checks_performed": [
            {
                "policy_id": 6,
                "policy_name": "Malware-Default-Policy",
                "control": "malware"
            },
            {
                "failed": true,
                "policy_id": 7,
                "policy_name": "Sensitive-Data-Default-Policy",
                "control": "sensitive_data",
                "sensitive_data_found": 1
            }
        ],
        "block_required": true
    },
    "vulnerability_summary": {
        "total": 12,
        "critical": 1,
        "high": 9,
        "medium": 2,
        "sensitive": 1,
        "malware": 3,
        "score_average": 7.058333
    },
    "scan_options": {
        "scan_executables": true,
        "scan_sensitive_data": true,
        "scan_malware": true,
        "scan_files": true,
        "scan_timeout": 3600000000000,
        "manual_pull_fallback": true,
        "use_cvss3": true,
        "dockerless": true,
        "telemetry_enabled": true,
        "enable_fast_scanning": true,
        "memoryThrottling": true,
        "suggest_os_upgrade": true
    },
    "previous_digest": "sha256:a563877b915679675b2b8ff3cd8b13863b40e9a43c218bed75e6aa14fcae23dc",
    "vulnerability_diff": {
        "critical": 1
    },
    "initiating_user": "administrator",
    "function_metadata": {},
    "scan_id": 12,
    "image_id": 6
}`)

func main() {

	settings := &data.OutputSettings{
		Name: "Test",
		Type: "slack",
		Url:  "https://hooks.slack.com/services/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",
	}

	// settings := &data.OutputSettings{
	// 	Name: "Test2",
	// 	Type: "teams",
	// 	Url:  "https://aquasecurity.webhook.office.com/webhookb2/7b8eb7a0-5589-4e14-9f89-2b9608bee40b@bc034cf3-566b-41ca-9f24-5dc49474b05e/IncomingWebhook/0354753665fe4606af3b5e5a4c852a70/846afe98-492a-4c55-8f4f-5d587608a784",
	// }

	input := make(map[string]interface{})
	err := json.Unmarshal(body, &input)
	if err != nil {
		log.Fatal(err)
	}

	err = router.TestOutput(input, settings)
	if err != nil {
		log.Fatal(err)
	}
}
