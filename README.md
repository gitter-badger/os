# RancherOS

[![Join the chat at https://gitter.im/imikushin/os](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/imikushin/os?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

The smallest, easiest way to run Docker in production at scale.  Everything in RancherOS is a container managed by Docker.  This includes system services such as udev and rsyslog.  RancherOS includes only the bare minimum amount of software needed to run Docker.  This keeps the binary download of RancherOS to about 25MB.  Everything else can be pulled in dynamically through Docker.

## How this works

Everything in RancherOS is a Docker container.  We accomplish this by launching two instances of
Docker.  One is what we call the system Docker which runs as PID 1.  System Docker then launches
a container that runs the user Docker.  The user Docker is then the instance that gets primarily
used to create containers.  We created this separation because it seemed logical and also
it would really be bad if somebody did `docker rm -f $(docker ps -qa)` and deleted the entire OS.

![How it works](docs/rancheros.png "How it works")


## Latest Release

**v0.4.1 - Docker 1.9.1-rc1 - Linux 4.2**

### ISO

https://releases.rancher.com/os/latest/rancheros.iso  
https://releases.rancher.com/os/v0.4.1/rancheros.iso  

**Note**: you must login using `rancher` for username and password.

### Additional Downloads

https://releases.rancher.com/os/latest/initrd
https://releases.rancher.com/os/latest/iso-checksums.txt
https://releases.rancher.com/os/latest/rancheros-v0.4.1.tar.gz
https://releases.rancher.com/os/latest/rancheros.iso
https://releases.rancher.com/os/latest/vmlinuz

https://releases.rancher.com/os/v0.4.1/initrd
https://releases.rancher.com/os/v0.4.1/iso-checksums.txt
https://releases.rancher.com/os/v0.4.1/rancheros-v0.4.1.tar.gz
https://releases.rancher.com/os/v0.4.1/rancheros.iso
https://releases.rancher.com/os/v0.4.1/vmlinuz

**Note**: you can use `http` instead of `https` in the above URLs, e.g. for iPXE.  

### Amazon

We have 2 different [virtualization types of AMIs](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/virtualization_types.html). SSH keys are added to the **`rancher`** user, so you must log in using the **rancher** user.

**HVM**

Region | Type | AMI |
-------|------|------
ap-northeast-1 | HVM |  [ami-8fd6f6e1](https://console.aws.amazon.com/ec2/home?region=ap-northeast-1#launchInstanceWizard:ami=ami-8fd6f6e1)
ap-southeast-1 | HVM |  [ami-db5395b8](https://console.aws.amazon.com/ec2/home?region=ap-southeast-1#launchInstanceWizard:ami=ami-db5395b8)
ap-southeast-2 | HVM |  [ami-010a5362](https://console.aws.amazon.com/ec2/home?region=ap-southeast-2#launchInstanceWizard:ami=ami-010a5362)
eu-central-1 | HVM |  [ami-0d637161](https://console.aws.amazon.com/ec2/home?region=eu-central-1#launchInstanceWizard:ami=ami-0d637161)
eu-west-1 | HVM |  [ami-9d6fb4ee](https://console.aws.amazon.com/ec2/home?region=eu-west-1#launchInstanceWizard:ami=ami-9d6fb4ee)
sa-east-1 | HVM |  [ami-227dc74e](https://console.aws.amazon.com/ec2/home?region=sa-east-1#launchInstanceWizard:ami=ami-227dc74e)
us-east-1 | HVM |  [ami-11740e7b](https://console.aws.amazon.com/ec2/home?region=us-east-1#launchInstanceWizard:ami=ami-11740e7b)
us-west-1 | HVM |  [ami-25422c45](https://console.aws.amazon.com/ec2/home?region=us-west-1#launchInstanceWizard:ami=ami-25422c45)
us-west-2 | HVM |  [ami-49cedf28](https://console.aws.amazon.com/ec2/home?region=us-west-2#launchInstanceWizard:ami=ami-49cedf28)

**Paravirtual**

Region | Type | AMI
---|--- | ---
ap-northeast-1 | PV |  [ami-2bd5f545](https://console.aws.amazon.com/ec2/home?region=ap-northeast-1#launchInstanceWizard:ami=ami-2bd5f545)
ap-southeast-1 | PV |  [ami-76529415](https://console.aws.amazon.com/ec2/home?region=ap-southeast-1#launchInstanceWizard:ami=ami-76529415)
ap-southeast-2 | PV |  [ami-0b085168](https://console.aws.amazon.com/ec2/home?region=ap-southeast-2#launchInstanceWizard:ami=ami-0b085168)
eu-central-1 | PV |  [ami-e8617384](https://console.aws.amazon.com/ec2/home?region=eu-central-1#launchInstanceWizard:ami=ami-e8617384)
eu-west-1 | PV |  [ami-686fb41b](https://console.aws.amazon.com/ec2/home?region=eu-west-1#launchInstanceWizard:ami=ami-686fb41b)
sa-east-1 | PV |  [ami-847ac0e8](https://console.aws.amazon.com/ec2/home?region=sa-east-1#launchInstanceWizard:ami=ami-847ac0e8)
us-east-1 | PV |  [ami-956a10ff](https://console.aws.amazon.com/ec2/home?region=us-east-1#launchInstanceWizard:ami=ami-956a10ff)
us-west-1 | PV |  [ami-5b412f3b](https://console.aws.amazon.com/ec2/home?region=us-west-1#launchInstanceWizard:ami=ami-5b412f3b)
us-west-2 | PV |  [ami-1ecedf7f](https://console.aws.amazon.com/ec2/home?region=us-west-2#launchInstanceWizard:ami=ami-1ecedf7f)

### Google Compute Engine (Experimental)

We are providing a disk image that users can download and import for use in Google Compute Engine. The image can be obtained from the release artifacts for RancherOS v0.3.0 or later.

[Download Image](https://github.com/rancher/os/releases/download/v0.4.1/rancheros-v0.4.1.tar.gz)

Please follow the directions at our [docs to launch in GCE](http://docs.rancher.com/os/running-rancheros/cloud/gce/). 

#### Known issues/ToDos
 * Add GCE daemon support. (Manages users)

## Documentation for RancherOS

Please refer to our [RancherOS Documentation](http://docs.rancher.com/os/) website to read all about RancherOS. It has detailed information on how RancherOS works, getting-started and other details.

## Support, Discussion, and Community
If you need any help with RancherOS or Rancher, please join us at either our [Rancher forums](http://forums.rancher.com) or [#rancher IRC channel](http://webchat.freenode.net/?channels=rancher) where most of our team hangs out at.

Please submit any **RancherOS** bugs, issues, and feature requests to [rancher/os](//github.com/rancher/os/issues).

Please submit any **Rancher** bugs, issues, and feature requests to [rancher/rancher](//github.com/rancher/rancher/issues).

#License
Copyright (c) 2014-2015 [Rancher Labs, Inc.](http://rancher.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

