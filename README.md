## Dagu for gokrazy

This package contains the static build of https://github.com/dagu-org/dagu

### Usage

Dagu is a lightweight workflow engine with a modern Web UI. Workflows are defined in a simple, declarative YAML format and can be executed on schedule. It supports shell commands, remote execution via SSH, and container images.

1. Install Dagu onto the remote machine:

```
gok add github.com/ascension-association/dagu
gok update
```

2. Load <device IP address>:8080 in your browser

3. Click on `DAG Definitions` in the left-hand column

4. Click on the `+ New` button in the top-right area

5. Call the DAG Name `test` and click `Create`

6. Scroll down to the `Definition` section and replace the contents with:

```
env:
  - PATH: /bin:/usr/local/bin

steps:
  - echo "It worked!"
```

7. Click the `Save Changes` button

8. Scroll to the top and click the play triangle button to run the test (click `Start` on the popup)

9. The `cmd_1` box outline should turn green and the _Run Status_ section should say `succeeded` and the _Steps_ section should show the result
