# Automatic Ingress Deployment - poleca.to

## Overview

This application listens for files in a specific directory. If the file contains hostname then it 
creates kubernetes ingress yaml and executes it.

This application is being used to create new k8s ingress for MediaPlanet company.

## Usage

1st argument is where you store new files.
2nd is where application should put yaml files.

kube-watch /tmp/new-files /tmp/yaml-files

## Error codes:

- 1: Invalid number of arguments.
- 2: Unable to create file watcher.