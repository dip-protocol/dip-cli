package main

import (
	"archive/zip"
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/dip-protocol/dip-cli/internal/proof"
	"github.com/dip-protocol/dip-cli/internal/verifyproof"
)

type Artifact struct {
	ArtifactVersion string    `json:"artifact_version"`
	ArtifactID      string    `json:"artifact_id"`
	Decision        any       `json:"decision"`
	Signature       Signature `json:"signature"`
}

type Signature struct {
	Algorithm string `json:"algorithm"`
	PublicKey []byte `json:"public_key"`
	Value     []byte `json:"value"`
}

func canonicalizeJSON(data []byte) ([]byte, error) {

	var obj any

	err := json.Unmarshal(data, &obj)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}

	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

	err = enc.Encode(obj)
	if err != nil {
		return nil, err
	}

	return bytes.TrimSpace(buf.Bytes()), nil
}

func computeArtifactID(canonical []byte) string {

	hash := sha256.Sum256(canonical)

	return hex.EncodeToString(hash[:])
}

func appendToRegistry() {

	cmd := exec.Command("..\\dip-registry\\registry.exe", "append", "artifact.json")

	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("Registry append failed:", err)
		return
	}

	fmt.Println(string(output))
}

func signDecision(inputFile string) {

	raw, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Error reading decision file:", err)
		return
	}

	canonical, err := canonicalizeJSON(raw)
	if err != nil {
		fmt.Println("Canonicalization error:", err)
		return
	}

	var decision any
	json.Unmarshal(canonical, &decision)

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Println("Key generation error:", err)
		return
	}

	artifactID := computeArtifactID(canonical)

	sig := ed25519.Sign(priv, canonical)

	artifact := Artifact{
		ArtifactVersion: "1",
		ArtifactID:      artifactID,
		Decision:        decision,
		Signature: Signature{
			Algorithm: "ed25519",
			PublicKey: pub,
			Value:     sig,
		},
	}

	out, err := json.MarshalIndent(artifact, "", "  ")
	if err != nil {
		fmt.Println("Artifact encoding error:", err)
		return
	}

	err = os.WriteFile("artifact.json", out, 0644)
	if err != nil {
		fmt.Println("Artifact write error:", err)
		return
	}

	fmt.Println("DIP artifact created: artifact.json")
	fmt.Println("Artifact ID:", artifactID)

	appendToRegistry()
}

func verifyArtifact(file string) {

	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("Verification error:", err)
		return
	}

	var artifact Artifact

	err = json.Unmarshal(data, &artifact)
	if err != nil {
		fmt.Println("Invalid artifact format")
		return
	}

	decisionBytes, err := json.Marshal(artifact.Decision)
	if err != nil {
		fmt.Println("Decision encoding error:", err)
		return
	}

	canonical, err := canonicalizeJSON(decisionBytes)
	if err != nil {
		fmt.Println("Canonicalization error:", err)
		return
	}

	expectedID := computeArtifactID(canonical)

	if expectedID != artifact.ArtifactID {
		fmt.Println("Artifact ID mismatch")
		return
	}

	valid := ed25519.Verify(
		artifact.Signature.PublicKey,
		canonical,
		artifact.Signature.Value,
	)

	if valid {
		fmt.Println("Artifact verification: VALID")
	} else {
		fmt.Println("Artifact verification: INVALID")
	}
}

func generateProof(file string) {

	err := proof.GenerateProof(file, "..\\dip-registry")

	if err != nil {
		fmt.Println("Proof generation failed:", err)
		return
	}

	fmt.Println("Proof generation complete")
}

func verifyProof(artifact string, proofFile string) {

	err := verifyproof.VerifyProof(artifact, proofFile)

	if err != nil {
		fmt.Println("Verification failed:", err)
	}
}

func bundle(artifact string, proofFile string) {

	bundleFile, err := os.Create("decision.dip")
	if err != nil {
		fmt.Println("Bundle creation failed:", err)
		return
	}

	defer bundleFile.Close()

	zipWriter := zip.NewWriter(bundleFile)
	defer zipWriter.Close()

	files := []string{artifact, proofFile}

	for _, file := range files {

		f, err := os.Open(file)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}

		defer f.Close()

		w, err := zipWriter.Create(file)
		if err != nil {
			fmt.Println("Error adding file to bundle:", err)
			return
		}

		_, err = io.Copy(w, f)
		if err != nil {
			fmt.Println("Error writing file to bundle:", err)
			return
		}
	}

	fmt.Println("DIP bundle created: decision.dip")
}

func verifyBundle(bundle string) {

	reader, err := zip.OpenReader(bundle)
	if err != nil {
		fmt.Println("Bundle open failed:", err)
		return
	}

	defer reader.Close()

	os.Mkdir("dip_tmp", 0755)

	for _, file := range reader.File {

		path := "dip_tmp/" + file.Name

		rc, err := file.Open()
		if err != nil {
			fmt.Println("Bundle read error:", err)
			return
		}

		out, err := os.Create(path)
		if err != nil {
			fmt.Println("File create error:", err)
			return
		}

		io.Copy(out, rc)

		out.Close()
		rc.Close()
	}

	artifact := "dip_tmp/artifact.json"
	proof := "dip_tmp/proof.json"

	err = verifyproof.VerifyProof(artifact, proof)

	if err != nil {
		fmt.Println("Bundle verification failed:", err)
		return
	}

	fmt.Println("DIP bundle verification: VALID")

	os.RemoveAll("dip_tmp")
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  dip sign <decision.json>")
		fmt.Println("  dip verify <artifact.json | decision.dip>")
		fmt.Println("  dip proof <artifact.json>")
		fmt.Println("  dip verify-proof <artifact.json> <proof.json>")
		fmt.Println("  dip bundle <artifact.json> <proof.json>")
		return
	}

	cmd := os.Args[1]

	switch cmd {

	case "sign":

		if len(os.Args) < 3 {
			fmt.Println("Usage: dip sign <decision.json>")
			return
		}

		signDecision(os.Args[2])

	case "verify":

		if len(os.Args) < 3 {
			fmt.Println("Usage: dip verify <artifact.json | decision.dip>")
			return
		}

		file := os.Args[2]

		if strings.HasSuffix(file, ".dip") {
			verifyBundle(file)
		} else {
			verifyArtifact(file)
		}

	case "proof":

		if len(os.Args) < 3 {
			fmt.Println("Usage: dip proof <artifact.json>")
			return
		}

		generateProof(os.Args[2])

	case "verify-proof":

		if len(os.Args) < 4 {
			fmt.Println("Usage: dip verify-proof <artifact.json> <proof.json>")
			return
		}

		verifyProof(os.Args[2], os.Args[3])

	case "bundle":

		if len(os.Args) < 4 {
			fmt.Println("Usage: dip bundle <artifact.json> <proof.json>")
			return
		}

		bundle(os.Args[2], os.Args[3])

	default:

		fmt.Println("Unknown command")
	}
}