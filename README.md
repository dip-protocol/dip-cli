\# DIP CLI



Reference command-line implementation for the \*\*Decision Integrity Protocol (DIP)\*\*.



`dip-cli` generates signed \*\*decision artifacts\*\* that can be verified independently.



---



\# Features



\* Create DIP artifacts

\* Deterministic canonicalization

\* Ed25519 signing

\* Artifact hash generation



---



\# Installation



Clone the repository and build:



```

go build

```



---



\# Usage



Create a decision file:



```

decision.json

```



Example:



```json

{

&nbsp; "decision\_id": "decision-001",

&nbsp; "timestamp": "2026-03-08T10:00:00Z",

&nbsp; "inputs": {

&nbsp;   "amount": 100

&nbsp; },

&nbsp; "outputs": {

&nbsp;   "approved": true

&nbsp; }

}

```



Generate an artifact:



```

dip sign decision.json

```



This produces:



```

artifact.json

```



---



\# Artifact Structure



```

artifact\_version

artifact\_id

decision

signature

```



---



\# Verification



Artifacts produced by `dip-cli` can be verified using:



```

dip-go-verifier

dip-js-verifier

```



---



\# License



Apache License 2.0



