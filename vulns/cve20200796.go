package vulns

import (
	"strings"
)

// Check if vulnerable
func (cve *CVE20200796) Check(build string, version string, kbs []string) bool {
	major := strings.Split(build, ".")[2]
	supersedence := []string{}
	switch major {
	// earlier versions are not vulnerable
	case "18362", "18363":
		supersedence = append(supersedence, "KB4551762", "KB4549951", "KB4550945", "KB4554364", "KB4541335")
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
func (cve *CVE20200796) Name() string {
	return "CVE-2020-0796"
}

// Description returns the Description
func (cve *CVE20200796) Description() string {
	return "SMBv3: Local Privilege Escalation \"SMBGhost\""
}

// CVE20200796 is the cve type
type CVE20200796 struct {
}
