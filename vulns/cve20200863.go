package vulns

import (
	"strings"
)

// Check if vulnerable
func (cve *CVE20200863) Check(build string, version string, kbs []string) bool {
	major := strings.Split(build, ".")[2]
	supersedence := []string{}
	switch major {
	case "18362", "18363":
		supersedence = append(supersedence, "KB4540673", "KB4551762", "KB4549951", "KB4550945", "KB4554364", "KB4541335")
	default:
		return false
	}
	for _, i := range kbs {
		for _, j := range supersedence {
			if i == j {
				return false
			}
		}
	}
	return true
}

// Name returns the CVE Number
func (cve *CVE20200863) Name() string {
	return "CVE-2020-0863"
}

// Description returns the Description
func (cve *CVE20200863) Description() string {
	return "Diagnostic Tracking Service: Arbitrary File Read as SYSTEM"
}

// CVE20200863 is the cve type
type CVE20200863 struct {
	name string
}
