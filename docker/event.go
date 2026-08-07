package docker

type Event struct {
	Type   EventType
	Action EventAction
	Actor  struct {
		Id         string `json:"ID"`
		Attributes EventAttributes
	}
	Scope    EventScope `json:"scope"`
	Time     int64      `json:"time"`
	TimeNano int64      `json:"timeNano"`
}

///////////////////////////////////////////////////////////////////////////////////////////////////

type EventAttributes struct {
	attributes map[string]string
}

func (this EventAttributes) Get(key string) (string, bool) {
	if this.attributes == nil {
		return "", false
	}
	val, ok := this.attributes[key]
	return val, ok
}

func (this EventAttributes) Image() (string, bool) {
	return this.Get("image")
}

func (this EventAttributes) Name() (string, bool) {
	return this.Get("name")
}

///////////////////////////////////////////////////////////////////////////////////////////////////

type EventType string

var EventTypes = eventTypes{
	Container: "container",
	Network:   "network",
}

type eventTypes struct {
	Container EventType
	Network   EventType
}

///////////////////////////////////////////////////////////////////////////////////////////////////

type EventAction string

var EventActions = eventActions{
	Start: "start",
	Stop:  "stop",
}

type eventActions struct {
	Start EventAction
	Stop  EventAction
}

///////////////////////////////////////////////////////////////////////////////////////////////////

type EventScope string

var EventScopes = eventScopes{
	Local: "local",
}

type eventScopes struct {
	Local EventScope
}

///////////////////////////////////////////////////////////////////////////////////////////////////

type ContainerListItem struct {
	Command      string
	CreatedAt    string
	HealthStatus string
	ID           string
	Image        string
	Labels       string
	LocalVolumes string
	Mounts       string
	Names        string
	Networks     string
	Platform     struct {
		Architecture string `json:"architecture"`
		OS           string `json:"os"`
		Variant      string `json:"variant"`
	}
	Ports      string
	RunningFor string
	Size       string
	State      string
	Status     string
}

///////////////////////////////////////////////////////////////////////////////////////////////////

type InspectItem struct {
	Id      string
	Created string
	Path    string
	Args    []string
	State   struct {
		Status     string
		Running    bool
		Paused     bool
		Restarting bool
		OOMKilled  bool
		Dead       bool
		Pid        int
		ExitCode   int
		Error      string
		StartedAt  string
		FinishedAt string
	}
	Image           string
	Name            string
	ResolvConfPath  string
	HostnamePath    string
	HostsPath       string
	LogPath         string
	Driver          string
	Platform        string
	MountLabel      string
	ProcessLabel    string
	AppArmorProfile string
	ExecIDs         []string
	HostConfig      struct {
		Binds           []string
		ContainerIDFile string
		LogConfig       struct {
			Type   string
			Config map[string]string
		}
		NetworkMode  string
		PortBindings map[string][]struct {
			HostIp   string `json:"HostIp"`
			HostPort string `json:"HostPort"`
		} `json:"PortBindings"`
		RestartPolicy struct {
			Name              string
			MaximumRetryCount int
		}
		AutoRemove   bool
		VolumeDriver string
		VolumesFrom  []string
		ConsoleSize  []int
		CapAdd       []string
		CapDrop      []string
		CgroupnsMode string
		Dns          []string
		DnsOptions   []string
		DnsSearch    []string
		ExtraHosts   []string
	}
}

// [
//   {
//     "HostConfig": {
//       "GroupAdd": null,
//       "IpcMode": "private",
//       "Cgroup": "",
//       "Links": null,
//       "OomScoreAdj": 0,
//       "PidMode": "",
//       "Privileged": false,
//       "PublishAllPorts": false,
//       "ReadonlyRootfs": false,
//       "SecurityOpt": null,
//       "UTSMode": "",
//       "UsernsMode": "",
//       "ShmSize": 67108864,
//       "Runtime": "runc",
//       "Isolation": "",
//       "CpuShares": 0,
//       "Memory": 0,
//       "NanoCpus": 0,
//       "CgroupParent": "",
//       "BlkioWeight": 0,
//       "BlkioWeightDevice": [],
//       "BlkioDeviceReadBps": [],
//       "BlkioDeviceWriteBps": [],
//       "BlkioDeviceReadIOps": [],
//       "BlkioDeviceWriteIOps": [],
//       "CpuPeriod": 0,
//       "CpuQuota": 0,
//       "CpuRealtimePeriod": 0,
//       "CpuRealtimeRuntime": 0,
//       "CpusetCpus": "",
//       "CpusetMems": "",
//       "Devices": [],
//       "DeviceCgroupRules": null,
//       "DeviceRequests": null,
//       "MemoryReservation": 0,
//       "MemorySwap": 0,
//       "MemorySwappiness": null,
//       "OomKillDisable": null,
//       "PidsLimit": null,
//       "Ulimits": [],
//       "CpuCount": 0,
//       "CpuPercent": 0,
//       "IOMaximumIOps": 0,
//       "IOMaximumBandwidth": 0,
//       "MaskedPaths": [
//         "/proc/acpi",
//         "/proc/asound",
//         "/proc/interrupts",
//         "/proc/kcore",
//         "/proc/keys",
//         "/proc/latency_stats",
//         "/proc/sched_debug",
//         "/proc/scsi",
//         "/proc/timer_list",
//         "/proc/timer_stats",
//         "/sys/devices/virtual/powercap",
//         "/sys/firmware"
//       ],
//       "ReadonlyPaths": [
//         "/proc/bus",
//         "/proc/fs",
//         "/proc/irq",
//         "/proc/sys",
//         "/proc/sysrq-trigger"
//       ]
//     },
//     "Storage": {
//       "RootFS": {
//         "Snapshot": {
//           "Name": "overlayfs"
//         }
//       }
//     },
//     "Mounts": [],
//     "Config": {
//       "Hostname": "0cfca5669eb6",
//       "Domainname": "",
//       "User": "",
//       "AttachStdin": false,
//       "AttachStdout": true,
//       "AttachStderr": true,
//       "Tty": false,
//       "OpenStdin": false,
//       "StdinOnce": false,
//       "Env": [
//         "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
//       ],
//       "Cmd": [
//         "/hello"
//       ],
//       "Image": "hello-world",
//       "Volumes": null,
//       "WorkingDir": "/",
//       "Entrypoint": null,
//       "Labels": {}
//     },
//     "NetworkSettings": {
//       "SandboxID": "",
//       "SandboxKey": "",
//       "Ports": {},
//       "Networks": {
//         "bridge": {
//           "IPAMConfig": null,
//           "Links": null,
//           "Aliases": null,
//           "DriverOpts": null,
//           "GwPriority": 0,
//           "NetworkID": "5f4de74bdaaf98dfae97ccdad3b16660bba5cc70e4893ed3ff06418f7137df86",
//           "EndpointID": "",
//           "Gateway": "",
//           "IPAddress": "",
//           "MacAddress": "",
//           "IPPrefixLen": 0,
//           "IPv6Gateway": "",
//           "GlobalIPv6Address": "",
//           "GlobalIPv6PrefixLen": 0,
//           "DNSNames": null
//         }
//       }
//     },
//     "ImageManifestDescriptor": {
//       "mediaType": "application/vnd.oci.image.manifest.v1+json",
//       "digest": "sha256:5099b89d7666cc2186cad769ddc262ddc7c335b33f5fe79f9ffe50a01282b23e",
//       "size": 1027,
//       "annotations": {
//         "com.docker.official-images.bashbrew.arch": "arm64v8",
//         "org.opencontainers.image.base.name": "scratch",
//         "org.opencontainers.image.created": "2026-03-23T21:33:57Z",
//         "org.opencontainers.image.revision": "0b0efba82b82ace81ab2fb42d25116f9488e6cb4",
//         "org.opencontainers.image.source": "https://github.com/docker-library/hello-world.git#0b0efba82b82ace81ab2fb42d25116f9488e6cb4:arm64v8",
//         "org.opencontainers.image.url": "https://hub.docker.com/_/hello-world",
//         "org.opencontainers.image.version": "linux"
//       },
//       "platform": {
//         "architecture": "arm64",
//         "os": "linux",
//         "variant": "v8"
//       }
//     }
//   }
// ]
