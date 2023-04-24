#!/bin/bash

host="${HOST:-localhost:3000}"

curl -X POST ${host}/tasks -d \
'{
    "name":"Windup",
    "state": "Ready",
    "locator": "windup",
    "addon": "windup",
    "application": {"id": 1},
    "data": {
        "rulesets": [
            {
                "id": 12,
                "name": "OpenShift Ruleset"
            }
        ]
    }
}' | jq -M .