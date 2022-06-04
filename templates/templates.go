package templates

const ChartYaml string = `apiVersion: {{.apiVersion}}
name: {{.name}}
description: {{.description}}

type: {{.type}}
version: {{.version}}
appVersion: {{.appVersion}}
`
