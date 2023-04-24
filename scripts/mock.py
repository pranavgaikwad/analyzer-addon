#!/usr/bin/python

from requests import request
from json import loads

base = "http://localhost:3000"
headers = { "Content-Type": "application/json" }

def create(url: str, payload: dict):
    response = request("POST", url, headers=headers, json=payload)
    print(response.text)

def get(url: str) -> dict:
    response = request("GET", url)
    return loads(response.text)

def jobfunctions(jfs: list[str]):
    for jf in jfs:
        create(f"{base}/jobfunctions", {
            "name": jf,
        })

def stakeholdergroups(shgroups: list[tuple[str,str]]):
    for name, description in shgroups:
        create(f"{base}/stakeholdergroups", {
                "name": name,
                "username": "default",
                "description": description,
            })

def businessservices(bs: list[tuple[str,str]]):
    for name, description in bs:
        create(f"{base}/businessservices", {
                "name": name,
                "description": description,
            })

def stakeholders(sh: list[dict]):
    jfs = get(f"{base}/jobfunctions")
    bss = get(f"{base}/businessservices")
    shg = get(f"{base}/stakeholdergroups")
    for payload in sh:
        stakeholder = {
            "name": payload["name"],
            "email": payload["email"],
            "businessServices": [],
            "stakeholderGroups": [],
        }
        for jf in jfs:
            if jf["name"] == payload["jobFunction"]["name"]:
                stakeholder["jobFunction"] = {
                    "id": jf["id"]
                }
        for bs in bss:
            for pbs in payload["businessServices"]:
                if bs["name"] == pbs["name"]:
                    stakeholder["businessServices"].append({
                        "id": bs["id"]
                    })
        for sh in shg:
            for psh in payload["stakeholderGroups"]:
                if sh["name"] == psh["name"]:
                    stakeholder["stakeholderGroups"].append({
                        "id": sh["id"]
                    })
        create(f"{base}/stakeholders", stakeholder)

jobfunctions([
    "Software Engineer",
    "IT Engineer",
    "SRE",
    "DevOps Engineer",
    "Consultant",
    "Engineering Manager",
])

stakeholdergroups([
    ("Migrators", "Stakeholders directly involved with migration"),
    ("Consultants", "Consultants helping customers modernize"),
    ("Field Enablement", "Field enablement teams"),
    ("Engineering", "Engineers working with the applications in the portfolio"),
    ("Technical Management", "Technical management folks closely involved with engineering"),
    ("Management", "Management not involved with engineering"),
    ("Operations", "Technical operations teams"),
])

businessservices([
    ("Retail", "Retail software"),
    ("Insurance", "Insurance software"),
    ("Compliance", "Compliance software"),
    ("R&D", "Ongoing research software"),
    ("Accounting", "Accounting software"),
])

stakeholders([
    {
        "name": "John Doe",
        "email": "john.doe@redhat.com",
        "jobFunction": {
            "name": "Software Engineer",
        },
        "stakeholderGroups": [
            {"name": "Migrators"},
            {"name": "Engineering"},
            {"name": "Operations"},
        ],
        "businessServices": [
            {"name": "Retail"},
            {"name": "Insurance"},
        ]
    },
    {
        "name": "Piff Jenkins",
        "email": "piff.jenkins@redhat.com",
        "jobFunction": {
            "name": "Software Engineer",
        },
        "stakeholderGroups": [
            {"name": "Migrators"},
            {"name": "Engineering"},
        ],
        "businessServices": [
            {"name": "Compliance"},
        ]
    },
    {
        "name": "Desmond Eagle",
        "email": "deagle@redhat.com",
        "jobFunction": {
            "name": "IT Engineer",
        },
        "stakeholderGroups": [
            {"name": "Migrators"},
            {"name": "Engineering"},
            {"name": "Operations"},
        ],
        "businessServices": [
            {"name": "Retail"},
            {"name": "Compliance"},
        ]
    },
    {
        "name": "Bodrum Salvador",
        "email": "brumsalv@redhat.com",
        "jobFunction": {
            "name": "Engineering Manager",
        },
        "stakeholderGroups": [
            {"name": "Technical Management"},
        ],
        "businessServices": [
            {"name": "Retail"},
        ]
    },
    {
        "name": "Russell Sprout",
        "email": "rsprout@redhat.com",
        "jobFunction": {
            "name": "Consultant",
        },
        "stakeholderGroups": [
            {"name": "Consultants"},
            {"name": "Field Enablement"},
        ],
        "businessServices": [
            {"name": "Compliance"},
            {"name": "Retail"},
        ]
    },
    {
        "name": "Indigo Violet",
        "email": "dviolet@redhat.com",
        "jobFunction": {
            "name": "Consultant",
        },
        "stakeholderGroups": [
            {"name": "Consultants"},
            {"name": "Field Enablement"},
        ],
        "businessServices": [
            {"name": "Insurance"},
            {"name": "Retail"},
            {"name": "Accounting"},
        ]
    },
])
