# Changelog

## [v0.0.8](https://github.com/kbrew-dev/kbrew/tree/v0.0.8) (2021-08-02)

[Full Changelog](https://github.com/kbrew-dev/kbrew/compare/v0.0.7...v0.0.8)

**Closed issues:**

- Raw args with int values are not working [\#94](https://github.com/kbrew-dev/kbrew/issues/94)

**Merged pull requests:**

- Correct binary install link in README [\#108](https://github.com/kbrew-dev/kbrew/pull/108) ([PrasadG193](https://github.com/PrasadG193))
- Add license to code files [\#106](https://github.com/kbrew-dev/kbrew/pull/106) ([gauravgahlot](https://github.com/gauravgahlot))
- Add Mergify integration to automatically merge PRs on approval [\#104](https://github.com/kbrew-dev/kbrew/pull/104) ([PrasadG193](https://github.com/PrasadG193))
- Switch release repo to kbrew [\#103](https://github.com/kbrew-dev/kbrew/pull/103) ([PrasadG193](https://github.com/PrasadG193))
- Bump helm.sh/helm/v3 from 3.5.4 to 3.6.1 [\#101](https://github.com/kbrew-dev/kbrew/pull/101) ([dependabot[bot]](https://github.com/apps/dependabot))

## [v0.0.7](https://github.com/kbrew-dev/kbrew/tree/v0.0.7) (2021-07-26)

[Full Changelog](https://github.com/kbrew-dev/kbrew/compare/v0.0.6...v0.0.7)

**Closed issues:**

- kbrew remove doesn't delete resources from pre-install steps [\#74](https://github.com/kbrew-dev/kbrew/issues/74)
- kbrew info command for information on packages [\#22](https://github.com/kbrew-dev/kbrew/issues/22)

**Merged pull requests:**

- render templates before unmarshalling [\#100](https://github.com/kbrew-dev/kbrew/pull/100) ([sahil-lakhwani](https://github.com/sahil-lakhwani))
- Use golangci-lint run for linting [\#95](https://github.com/kbrew-dev/kbrew/pull/95) ([sanketsudake](https://github.com/sanketsudake))
- Logging framework [\#93](https://github.com/kbrew-dev/kbrew/pull/93) ([PrasadG193](https://github.com/PrasadG193))

## [v0.0.6](https://github.com/kbrew-dev/kbrew/tree/v0.0.6) (2021-07-01)

[Full Changelog](https://github.com/kbrew-dev/kbrew/compare/v0.0.5...v0.0.6)

**Closed issues:**

- Create namespace if does not exist for raw apps [\#80](https://github.com/kbrew-dev/kbrew/issues/80)
- Installation should continue if any dependency app already exists [\#79](https://github.com/kbrew-dev/kbrew/issues/79)

**Merged pull requests:**

- Silence usage and error print [\#92](https://github.com/kbrew-dev/kbrew/pull/92) ([sahil-lakhwani](https://github.com/sahil-lakhwani))
- Skip namespace creation if name is nil [\#91](https://github.com/kbrew-dev/kbrew/pull/91) ([PrasadG193](https://github.com/PrasadG193))
- Use correct reference while creating namespace [\#90](https://github.com/kbrew-dev/kbrew/pull/90) ([PrasadG193](https://github.com/PrasadG193))
- Enable code scan with CodeQL for Go source [\#89](https://github.com/kbrew-dev/kbrew/pull/89) ([sanketsudake](https://github.com/sanketsudake))
- add info and args command [\#88](https://github.com/kbrew-dev/kbrew/pull/88) ([sahil-lakhwani](https://github.com/sahil-lakhwani))
- App cleanup improvements [\#87](https://github.com/kbrew-dev/kbrew/pull/87) ([PrasadG193](https://github.com/PrasadG193))
- Skip helm app installation if already exists [\#85](https://github.com/kbrew-dev/kbrew/pull/85) ([PrasadG193](https://github.com/PrasadG193))
- Fix incorrect exit codes [\#84](https://github.com/kbrew-dev/kbrew/pull/84) ([PrasadG193](https://github.com/PrasadG193))
- create namespace if it doesn't exist  [\#83](https://github.com/kbrew-dev/kbrew/pull/83) ([mahendrabagul](https://github.com/mahendrabagul))

## [v0.0.5](https://github.com/kbrew-dev/kbrew/tree/v0.0.5) (2021-06-02)

[Full Changelog](https://github.com/kbrew-dev/kbrew/compare/v0.0.4...v0.0.5)

**Closed issues:**

- Add a flag to set custom wait for ready timeout [\#75](https://github.com/kbrew-dev/kbrew/issues/75)
- Collect events for analytics [\#65](https://github.com/kbrew-dev/kbrew/issues/65)
- Support for helm functions in recipe yaml [\#64](https://github.com/kbrew-dev/kbrew/issues/64)
- Add support for passing args to the kbrew app [\#63](https://github.com/kbrew-dev/kbrew/issues/63)
- Add support for kbrew help command [\#57](https://github.com/kbrew-dev/kbrew/issues/57)

**Merged pull requests:**

- Support helm functions on raw app args [\#81](https://github.com/kbrew-dev/kbrew/pull/81) ([PrasadG193](https://github.com/PrasadG193))
- Add autocomplete support for kbrew CLI [\#78](https://github.com/kbrew-dev/kbrew/pull/78) ([sanketsudake](https://github.com/sanketsudake))
- Collect events for analytics [\#77](https://github.com/kbrew-dev/kbrew/pull/77) ([PrasadG193](https://github.com/PrasadG193))
- Support custom wait-for-ready timeout [\#76](https://github.com/kbrew-dev/kbrew/pull/76) ([PrasadG193](https://github.com/PrasadG193))
- \[CI\] Fix go vet failures [\#72](https://github.com/kbrew-dev/kbrew/pull/72) ([PrasadG193](https://github.com/PrasadG193))
- Support for raw app args [\#71](https://github.com/kbrew-dev/kbrew/pull/71) ([sahil-lakhwani](https://github.com/sahil-lakhwani))
- change args value to empty interface [\#70](https://github.com/kbrew-dev/kbrew/pull/70) ([sahil-lakhwani](https://github.com/sahil-lakhwani))
- yaml evaluator [\#69](https://github.com/kbrew-dev/kbrew/pull/69) ([sahil-lakhwani](https://github.com/sahil-lakhwani))
- Enhanced version command to check for latest version and warn user [\#68](https://github.com/kbrew-dev/kbrew/pull/68) ([vishal-biyani](https://github.com/vishal-biyani))

## [v0.0.4](https://github.com/kbrew-dev/kbrew/tree/v0.0.4) (2021-05-13)

[Full Changelog](https://github.com/kbrew-dev/kbrew/compare/v0.0.3...v0.0.4)

**Closed issues:**

- Add Go lint checks in CI [\#37](https://github.com/kbrew-dev/kbrew/issues/37)

**Merged pull requests:**

- Args templating engine [\#67](https://github.com/kbrew-dev/kbrew/pull/67) ([sahil-lakhwani](https://github.com/sahil-lakhwani))
- arguments to helm app [\#66](https://github.com/kbrew-dev/kbrew/pull/66) ([sahil-lakhwani](https://github.com/sahil-lakhwani))
- Remove stale recipes dir [\#61](https://github.com/kbrew-dev/kbrew/pull/61) ([PrasadG193](https://github.com/PrasadG193))
- Added link to registry [\#53](https://github.com/kbrew-dev/kbrew/pull/53) ([vishal-biyani](https://github.com/vishal-biyani))
- Add go lint checks [\#52](https://github.com/kbrew-dev/kbrew/pull/52) ([mahendrabagul](https://github.com/mahendrabagul))
- fix exact match for registry [\#51](https://github.com/kbrew-dev/kbrew/pull/51) ([sahil-lakhwani](https://github.com/sahil-lakhwani))
- Add support of Mac build [\#50](https://github.com/kbrew-dev/kbrew/pull/50) ([PrasadG193](https://github.com/PrasadG193))

## [v0.0.3](https://github.com/kbrew-dev/kbrew/tree/v0.0.3) (2021-04-05)

[Full Changelog](https://github.com/kbrew-dev/kbrew/compare/v0.0.2...v0.0.3)

**Closed issues:**

- Install script does not support Mac [\#49](https://github.com/kbrew-dev/kbrew/issues/49)

## [v0.0.2](https://github.com/kbrew-dev/kbrew/tree/v0.0.2) (2021-04-05)

[Full Changelog](https://github.com/kbrew-dev/kbrew/compare/v0.0.1...v0.0.2)

**Merged pull requests:**

- \[README\] Update installation steps and command structs [\#48](https://github.com/kbrew-dev/kbrew/pull/48) ([PrasadG193](https://github.com/PrasadG193))
- Add support to update kbrew registries [\#47](https://github.com/kbrew-dev/kbrew/pull/47) ([PrasadG193](https://github.com/PrasadG193))
- Install kbrew in existing exec bin dir during update [\#46](https://github.com/kbrew-dev/kbrew/pull/46) ([PrasadG193](https://github.com/PrasadG193))

## [v0.0.1](https://github.com/kbrew-dev/kbrew/tree/v0.0.1) (2021-04-04)

[Full Changelog](https://github.com/kbrew-dev/kbrew/compare/55e699e0ef038af0f8e29bbca6b3697752e1fe77...v0.0.1)

**Closed issues:**

- Add release process/script [\#39](https://github.com/kbrew-dev/kbrew/issues/39)
- Host and fetch recipes at runtime [\#28](https://github.com/kbrew-dev/kbrew/issues/28)
- Recipe: Rook + Ceph [\#21](https://github.com/kbrew-dev/kbrew/issues/21)
- Trial/Recipe: Longhorn [\#20](https://github.com/kbrew-dev/kbrew/issues/20)
- Trial/Recipe: OpenEBS [\#19](https://github.com/kbrew-dev/kbrew/issues/19)
- Add GitHub actions CI to run build and test [\#11](https://github.com/kbrew-dev/kbrew/issues/11)
- Add preinstall/postinstall hooks in installation workflow [\#6](https://github.com/kbrew-dev/kbrew/issues/6)
- Wait resources to be ready after install/upgrade [\#5](https://github.com/kbrew-dev/kbrew/issues/5)
- Add support for repository in config [\#4](https://github.com/kbrew-dev/kbrew/issues/4)

**Merged pull requests:**

- Add update command to upgrade kbrew [\#45](https://github.com/kbrew-dev/kbrew/pull/45) ([PrasadG193](https://github.com/PrasadG193))
- Add kbrew version command [\#44](https://github.com/kbrew-dev/kbrew/pull/44) ([PrasadG193](https://github.com/PrasadG193))
- Publish releases to kbrew-release repo [\#43](https://github.com/kbrew-dev/kbrew/pull/43) ([PrasadG193](https://github.com/PrasadG193))
- Add release process [\#42](https://github.com/kbrew-dev/kbrew/pull/42) ([PrasadG193](https://github.com/PrasadG193))
- Add support for kbrew registries [\#35](https://github.com/kbrew-dev/kbrew/pull/35) ([PrasadG193](https://github.com/PrasadG193))
- Cleanup - remove redundant files [\#34](https://github.com/kbrew-dev/kbrew/pull/34) ([PrasadG193](https://github.com/PrasadG193))
- Change repo owner to kbrew-dev [\#33](https://github.com/kbrew-dev/kbrew/pull/33) ([PrasadG193](https://github.com/PrasadG193))
- Revert "Add Mergify integration to automatically merge PRs on approval" [\#32](https://github.com/kbrew-dev/kbrew/pull/32) ([PrasadG193](https://github.com/PrasadG193))
- Add Mergify integration to automatically merge PRs on approval [\#31](https://github.com/kbrew-dev/kbrew/pull/31) ([PrasadG193](https://github.com/PrasadG193))
- Fix resource name in WaitForReady methods [\#30](https://github.com/kbrew-dev/kbrew/pull/30) ([PrasadG193](https://github.com/PrasadG193))
- Added details to readme such as instructions and usage commands [\#26](https://github.com/kbrew-dev/kbrew/pull/26) ([vishal-biyani](https://github.com/vishal-biyani))
-  Set app/chart versions in the recipes [\#18](https://github.com/kbrew-dev/kbrew/pull/18) ([PrasadG193](https://github.com/PrasadG193))
- Add recipe for kafka operator installation [\#17](https://github.com/kbrew-dev/kbrew/pull/17) ([PrasadG193](https://github.com/PrasadG193))
- Add support for shell command in install hooks [\#15](https://github.com/kbrew-dev/kbrew/pull/15) ([PrasadG193](https://github.com/PrasadG193))
- Add support for pre/post install hooks [\#14](https://github.com/kbrew-dev/kbrew/pull/14) ([PrasadG193](https://github.com/PrasadG193))
- Remove redundant file [\#13](https://github.com/kbrew-dev/kbrew/pull/13) ([PrasadG193](https://github.com/PrasadG193))
- Added a simple workflow to build and test code [\#12](https://github.com/kbrew-dev/kbrew/pull/12) ([vishal-biyani](https://github.com/vishal-biyani))
- Update imports to point to infracloudio repo [\#10](https://github.com/kbrew-dev/kbrew/pull/10) ([PrasadG193](https://github.com/PrasadG193))
- Wait for workloads to be ready after install [\#9](https://github.com/kbrew-dev/kbrew/pull/9) ([PrasadG193](https://github.com/PrasadG193))
- Add support for app repository [\#7](https://github.com/kbrew-dev/kbrew/pull/7) ([PrasadG193](https://github.com/PrasadG193))
- Add support to pass and parse config [\#3](https://github.com/kbrew-dev/kbrew/pull/3) ([PrasadG193](https://github.com/PrasadG193))
- Add support to install, remove and search apps [\#2](https://github.com/kbrew-dev/kbrew/pull/2) ([PrasadG193](https://github.com/PrasadG193))



