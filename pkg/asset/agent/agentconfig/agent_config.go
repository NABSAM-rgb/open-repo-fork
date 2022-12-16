package agentconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/yaml"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/types/agent"
	"github.com/openshift/installer/pkg/types/agent/conversion"
	"github.com/openshift/installer/pkg/validate"
)

var (
	agentConfigFilename = "agent-config.yaml"
)

// AgentConfig reads the agent-config.yaml file.
type AgentConfig struct {
	File     *asset.File
	Config   *agent.Config
	Template string
}

var _ asset.WritableAsset = (*AgentConfig)(nil)

// Name returns a human friendly name for the asset.
func (*AgentConfig) Name() string {
	return "Agent Config"
}

// Dependencies returns all of the dependencies directly needed to generate
// the asset.
func (*AgentConfig) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates the Agent Config manifest.
func (a *AgentConfig) Generate(dependencies asset.Parents) error {

	// TODO: We are temporarily generating a template of the agent-config.yaml
	// Change this when its interactive survey is implemented.
	agentConfigTemplate := `#
# Note: This is a sample AgentConfig file showing
# which fields are available to aid you in creating your
# own agent-config.yaml file.
#
apiVersion: v1alpha1
kind: AgentConfig
metadata:
  name: example-agent-config
  namespace: cluster0
# All fields are optional
rendezvousIP: your-node0-ip
hosts:
# If a host is listed, then at least one interface
# needs to be specified.
- hostname: change-to-hostname
  role: master
  # For more information about rootDeviceHints:
  # https://docs.openshift.com/container-platform/4.10/installing/installing_bare_metal_ipi/ipi-install-installation-workflow.html#root-device-hints_ipi-install-installation-workflow
  rootDeviceHints:
    deviceName: /dev/sda
  # interfaces are used to identify the host to apply this configuration to
  interfaces:
    - macAddress: 00:00:00:00:00:00
      name: host-network-interface-name
  # networkConfig contains the network configuration for the host in NMState format.
  # See https://nmstate.io/examples.html for examples.
  networkConfig:
    interfaces:
      - name: eth0
        type: ethernet
        state: up
        mac-address: 00:00:00:00:00:00
        ipv4:
          enabled: true
          address:
            - ip: 192.168.122.2
              prefix-length: 23
          dhcp: false
`

	a.Template = agentConfigTemplate

	// TODO: template is not validated
	return nil
}

// PersistToFile writes the agent-config.yaml file to the assets folder
func (a *AgentConfig) PersistToFile(directory string) error {
	templatePath := filepath.Join(directory, agentConfigFilename)
	templateByte := []byte(a.Template)

	err := os.WriteFile(templatePath, templateByte, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Files returns the files generated by the asset.
func (a *AgentConfig) Files() []*asset.File {
	if a.File != nil {
		return []*asset.File{a.File}
	}
	return []*asset.File{}
}

// Load returns agent config asset from the disk.
func (a *AgentConfig) Load(f asset.FileFetcher) (bool, error) {

	file, err := f.FetchByName(agentConfigFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.Wrap(err, fmt.Sprintf("failed to load %s file", agentConfigFilename))
	}

	config := &agent.Config{}
	if err := yaml.UnmarshalStrict(file.Data, config); err != nil {
		return false, errors.Wrapf(err, "failed to unmarshal %s", agentConfigFilename)
	}

	// Upconvert any deprecated fields
	if err := conversion.ConvertAgentConfig(config); err != nil {
		return false, err
	}

	a.File, a.Config = file, config
	if err = a.finish(); err != nil {
		return false, err
	}

	return true, nil
}

func (a *AgentConfig) finish() error {
	if err := a.validateAgent().ToAggregate(); err != nil {
		return errors.Wrapf(err, "invalid Agent Config configuration")
	}

	return nil
}

func (a *AgentConfig) validateAgent() field.ErrorList {
	var allErrs field.ErrorList

	if err := a.validateRendezvousIP(); err != nil {
		allErrs = append(allErrs, err...)
	}

	if err := a.validateHosts(); err != nil {
		allErrs = append(allErrs, err...)
	}

	if err := a.validateAdditionalNTPSources(field.NewPath("AdditionalNTPSources"), a.Config.AdditionalNTPSources); err != nil {
		allErrs = append(allErrs, err...)
	}

	if err := a.validateRendevousIPNotWorker(a.Config.RendezvousIP, a.Config.Hosts); err != nil {
		allErrs = append(allErrs, err...)
	}

	return allErrs
}

func (a *AgentConfig) validateRendezvousIP() field.ErrorList {
	var allErrs field.ErrorList

	rendezvousIPPath := field.NewPath("rendezvousIP")

	//empty rendezvous ip is fine
	if a.Config.RendezvousIP == "" {
		return nil
	}

	if err := validate.IP(a.Config.RendezvousIP); err != nil {
		allErrs = append(allErrs, field.Invalid(rendezvousIPPath, a.Config.RendezvousIP, err.Error()))
	}

	return allErrs
}

func (a *AgentConfig) validateHosts() field.ErrorList {
	var allErrs field.ErrorList

	macs := make(map[string]bool)
	for i, host := range a.Config.Hosts {

		hostPath := field.NewPath("Hosts").Index(i)

		if err := a.validateHostInterfaces(hostPath, host, macs); err != nil {
			allErrs = append(allErrs, err...)
		}

		if err := a.validateHostRootDeviceHints(hostPath, host); err != nil {
			allErrs = append(allErrs, err...)
		}

		if err := a.validateRoles(hostPath, host); err != nil {
			allErrs = append(allErrs, err...)
		}
	}

	return allErrs
}

func (a *AgentConfig) validateHostInterfaces(hostPath *field.Path, host agent.Host, macs map[string]bool) field.ErrorList {
	var allErrs field.ErrorList

	interfacePath := hostPath.Child("Interfaces")
	if len(host.Interfaces) == 0 {
		allErrs = append(allErrs, field.Required(interfacePath, "at least one interface must be defined for each node"))
	}

	for j := range host.Interfaces {
		mac := host.Interfaces[j].MacAddress
		macAddressPath := interfacePath.Index(j).Child("macAddress")

		if mac == "" {
			allErrs = append(allErrs, field.Required(macAddressPath, "each interface must have a MAC address defined"))
			continue
		}

		if err := validate.MAC(mac); err != nil {
			allErrs = append(allErrs, field.Invalid(macAddressPath, mac, err.Error()))
		}

		if _, ok := macs[mac]; ok {
			allErrs = append(allErrs, field.Invalid(macAddressPath, mac, "duplicate MAC address found"))
		}
		macs[mac] = true
	}

	return allErrs
}

func (a *AgentConfig) validateHostRootDeviceHints(hostPath *field.Path, host agent.Host) field.ErrorList {
	var allErrs field.ErrorList

	if host.RootDeviceHints.WWNWithExtension != "" {
		allErrs = append(allErrs, field.Forbidden(
			hostPath.Child("RootDeviceHints", "WWNWithExtension"), "WWN extensions are not supported in root device hints"))
	}

	if host.RootDeviceHints.WWNVendorExtension != "" {
		allErrs = append(allErrs, field.Forbidden(hostPath.Child("RootDeviceHints", "WWNVendorExtension"), "WWN vendor extensions are not supported in root device hints"))
	}

	return allErrs
}

func (a *AgentConfig) validateRoles(hostPath *field.Path, host agent.Host) field.ErrorList {
	var allErrs field.ErrorList

	if len(host.Role) > 0 && host.Role != "master" && host.Role != "worker" {
		allErrs = append(allErrs, field.Forbidden(hostPath.Child("Host"), "host role has incorrect value. Role must either be 'master' or 'worker'"))
	}

	return allErrs
}

func (a *AgentConfig) validateAdditionalNTPSources(additionalNTPSourcesPath *field.Path, sources []string) field.ErrorList {
	var allErrs field.ErrorList

	for i, source := range sources {
		domainNameErr := validate.DomainName(source, true)
		if domainNameErr != nil {
			ipErr := validate.IP(source)
			if ipErr != nil {
				allErrs = append(allErrs, field.Invalid(additionalNTPSourcesPath.Index(i), source, "NTP source is not a valid domain name nor a valid IP"))
			}
		}
	}

	return allErrs
}

func (a *AgentConfig) validateRendevousIPNotWorker(rendezvousIP string, hosts []agent.Host) field.ErrorList {
	var allErrs field.ErrorList

	if rendezvousIP != "" {
		for i, host := range hosts {
			hostPath := field.NewPath("Hosts").Index(i)
			if strings.Contains(string(host.NetworkConfig.Raw), rendezvousIP) && host.Role != "master" {
				if len(host.Role) > 0 {
					errMsg := "Host " + host.Hostname + " is not of role 'master' and has the rendevousIP assigned to it. The rendevousIP must be assigned to a host of role 'master'"
					allErrs = append(allErrs, field.Forbidden(hostPath.Child("Host"), errMsg))
				}
			}
		}
	}

	return allErrs
}

// HostConfigFileMap is a map from a filepath ("<host>/<file>") to file content
// for hostconfig files.
type HostConfigFileMap map[string][]byte

// HostConfigFiles returns a map from filename to contents of the files used for
// host-specific configuration by the agent installer client
func (a *AgentConfig) HostConfigFiles() (HostConfigFileMap, error) {
	if a == nil || a.Config == nil {
		return nil, nil
	}

	files := HostConfigFileMap{}
	for i, host := range a.Config.Hosts {
		name := fmt.Sprintf("host-%d", i)
		if host.Hostname != "" {
			name = host.Hostname
		}

		macs := []string{}
		for _, iface := range host.Interfaces {
			macs = append(macs, strings.ToLower(iface.MacAddress)+"\n")
		}

		if len(macs) > 0 {
			files[filepath.Join(name, "mac_addresses")] = []byte(strings.Join(macs, ""))
		}

		rdh, err := yaml.Marshal(host.RootDeviceHints)
		if err != nil {
			return nil, err
		}
		if len(rdh) > 0 && string(rdh) != "{}\n" {
			files[filepath.Join(name, "root-device-hints.yaml")] = rdh
		}

		if len(host.Role) > 0 {
			files[filepath.Join(name, "role")] = []byte(host.Role)
		}
	}
	return files, nil
}

func unmarshalJSON(b []byte) []byte {
	output, _ := yaml.JSONToYAML(b)
	return output
}
