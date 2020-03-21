package main_test

import (
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
)

const threatString = `main(): The map file wasn't found, symbols wont be available
	main(): Scanning /malware/EICAR...
	EngineScanCallback(): Scanning input
	EngineScanCallback(): Threat Virus:DOS/EICAR_Test_File identified.
  `

const cleanString = `main(): The map file wasn't found, symbols wont be available
	main(): Scanning /malware/EICAR...
	EngineScanCallback(): Scanning input
`

func parseWindowsDefenderOutput(windefout string) (string, error) {

	lines := strings.Split(windefout, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Scanning input") {
			continue
		}
		if strings.Contains(line, "EngineScanCallback") {
			threat := strings.TrimPrefix(strings.TrimSpace(line), "EngineScanCallback():")
			if len(threat) > 0 {
				threat = strings.TrimSpace(threat)
				threat = strings.TrimPrefix(threat, "Threat")
				threat = strings.TrimSuffix(threat, "identified.")
				return strings.TrimSpace(threat), nil
			} else {
				log.Errorf("Umm... len(threat)=%d, threat=%v", len(threat), threat)
			}
		}
	}
	return "CLEAN", nil
}

// TestParseResult tests the ParseFSecureOutput function.
func TestParseThreat(t *testing.T) {

	results, err := parseWindowsDefenderOutput(threatString)

	if err != nil {
		t.Log(err)
	}

	if !strings.EqualFold(results, "Virus:DOS/EICAR_Test_File") {
		t.Error("Threat incorrectly extracted")
		t.Log("results: ", results)
	}

}

// TestParseResult tests the ParseFSecureOutput function.
func TestParseClean(t *testing.T) {

	results, err := parseWindowsDefenderOutput(cleanString)

	if err != nil {
		t.Log(err)
	}

	if !strings.EqualFold(results, "CLEAN") {
		t.Error("Threat incorrectly extracted")
		t.Log("results: ", results)
	}

}
