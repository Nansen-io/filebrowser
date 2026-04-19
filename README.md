<div align="center">
  <h3>AcornDrive</h3>
  AcornDrive based on Quantum Fielbrowser is a part of the Acorn.Tools toolset.It is integrated witht he Chain-FS API and represents a fork of the filebrowser code.
  <br/><br/>
</div>

## Pinned


## About

AcornDrive provides an easy way to access and manage your files from the web. It has a modern responsive interface that has many advanced features to manage users, access, sharing, data protection and file preview and editing.

# Deploy to DEV
git push origin dev

# Deploy to UAT
git push origin uat

# Deploy to PROD
git push origin main

RUN: cd /mnt/c/Users/Andrew/Development/filebrowser/frontend && npm run build && cd ../backend && go build -o filebrowser && ./filebrowser --chainfs-bypass -c config.dev.yaml

