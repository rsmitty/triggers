Usage:

```bash
docker run -ti -e PROJECT="testbed-xxxx" -e SUBSCRIPTION="test-sub" -e CHANNEL="mychan" -e NAMESPACE="rsmitty" -e GOOGLE_APPLICATION_CREDENTIALS="/tmp/creds.json" -v /path/to/creds.json:/tmp/creds.json rsmitty/gcssource:0.1.0
```

TODO:
- Verify env vars before issuing
- Create a kube manifest for deployment