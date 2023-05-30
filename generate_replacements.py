import re
import urllib.request
import subprocess

url = "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/kustomize/go.mod"

try:
    with urllib.request.urlopen(url) as response:
        content = response.read().decode('utf-8')

    matches = re.findall(r'\s([a-z0-9_.\/-]+)\sv([a-z0-9_.\/+-]+)', content, re.I)

    for match in matches:
        print(f"{match[0]} v{match[1]}")
        command = "go mod edit -replace \"%s=%s@v%s\"" % (match[0], match[0], match[1])
        print(command)
        output = subprocess.check_output(command, shell=True, text=True)
        print(output)

except urllib.error.HTTPError as e:
    print("Failed to download file. HTTP Error:", e.code)
except urllib.error.URLError as e:
    print("Failed to download file. URL Error:", e.reason)
